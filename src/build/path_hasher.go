package build

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"core"
)

// pathHasher implements hashing of arbitrary paths to help with build incrementality.
// It memoizes hashes internally to improve performance, which is *almost* transparent -
// we can overwrite previously hashed files in the case where we rebuild a target.
// Please is often performance-bound by hashing so the memoisation is still worth a little
// exposed complexity.
type pathHasher struct {
	pathHashMemoizer map[string][]byte `guard:"pathHashMutex"`
	pathHashMutex    sync.RWMutex
}

// hasher is the singleton pathHasher.
var hasher = pathHasher{
	pathHashMemoizer: map[string][]byte{},
}

// PathHash calculates the hash of a single path which might be a file or a directory.
// This is the memoized form that only hashes each path once, unless recalc is true in which
// case it will force a recalculation of the hash.
// If sha1 is true then we will calculate a sha1 hash for this target, otherwise we use sha256
// for output hashes.
func (p *pathHasher) PathHash(path string, recalc bool, sha1 bool) ([]byte, error) {
	path = p.ensureRelative(path)
	if !recalc {
		p.pathHashMutex.RLock()
		cached, present := p.pathHashMemoizer[path]
		p.pathHashMutex.RUnlock()
		if present {
			return cached, nil
		}
	}
	result, err := p.pathHash(path, sha1)
	if err == nil {
		p.pathHashMutex.Lock()
		p.pathHashMemoizer[path] = result
		p.pathHashMutex.Unlock()
	}
	return result, err
}

// TargetPathHash is a convenience function that calculates a path hash for a build target,
// working out whether it needs SHA1 or not.
func (p *pathHasher) TargetPathHash(path string, target *core.BuildTarget) ([]byte, error) {
	return p.PathHash(path, false, p.NeedsSHA1Hash(target))
}

// NeedsSHA1Hash returns true for a target that has any hashes specified on it that are sha1.
// Note that the logic here implies that you cannot mix sha1 and sha256 hashes on a target.
// We could of course fix that but it seems unnecessary since the goal is to migrate off
// sha1 hashes in BUILD files completely.
func (p *pathHasher) NeedsSHA1Hash(target *core.BuildTarget) bool {
	for _, h := range target.Hashes {
		// SHA1 hashes are 40 characters long when hex-encoded. SHA256 would be 64 characters.
		if strings.HasPrefix(h, "sha1:") || (len(canonicalHash(h)) == 40 && !strings.HasPrefix(h, "sha256:")) {
			return true
		}
	}
	return false
}

// MustPathHash is as PathHash, but panics if the hash cannot be calculated.
func (p *pathHasher) MustPathHash(path string, sha1 bool) []byte {
	hash, err := p.PathHash(path, false, sha1)
	if err != nil {
		panic(err)
	}
	return hash
}

// MustTargetPathHash is as TargetPathHash, but panics if the hash cannot be calculated.
func (p *pathHasher) MustTargetPathHash(path string, target *core.BuildTarget) []byte {
	return p.MustPathHash(path, p.NeedsSHA1Hash(target))
}

// HashType returns an appropriate hash type (i.e. either sha1 or sha256)
// Eventually we will remove this and use sha256 exclusively.
func (p *pathHasher) HashType(needsSha1 bool) hash.Hash {
	if needsSha1 {
		return sha1.New()
	}
	return sha256.New()
}

// pathHash calculates a hash for one particular path.
func (p *pathHasher) pathHash(path string, sha1 bool) ([]byte, error) {
	h := p.HashType(sha1)
	info, err := os.Lstat(path)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		// Dereference symlink and try again
		deref, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, err
		}
		// Write something indicating this was a link; important so we rebuild correctly
		// when eg. a filegroup is changed from link=False to link=True.
		// Don't want to hash all file mode bits, the others could change depending on
		// whether we retrieved from cache or not so they're probably a bit too fragile.
		// TODO(pebers): is this still correct? filegroups no longer have link, and we attempt
		//               to disallow symlinks in plz-out more generally, but other bits
		//               could change (e.g. executable...)
		h.Write(boolTrueHashValue)
		d, err := p.pathHash(deref, sha1)
		h.Write(d)
		return h.Sum(nil), err
	}
	if err == nil && info.IsDir() {
		err = filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.Mode()&os.ModeSymlink != 0 {
				// Is a symlink, must verify that it's not a link outside the tmp dir.
				deref, err := filepath.EvalSymlinks(name)
				if err != nil {
					return err
				}
				if !strings.HasPrefix(deref, path) {
					return fmt.Errorf("Output %s links outside the build dir (to %s)", name, deref)
				}
				// Deliberately do not attempt to read it. We will read the contents later since
				// it is a link within the temp dir anyway, and if it's a link to a directory
				// it can introduce a cycle.
				// Just write something to the hash indicating that we found something here,
				// otherwise rules might be marked as unchanged if they added additional symlinks.
				h.Write(boolTrueHashValue)
			} else if !info.IsDir() {
				return p.fileHash(&h, name)
			}
			return nil
		})
	} else {
		err = p.fileHash(&h, path) // let this handle any other errors
	}
	return h.Sum(nil), err
}

// MovePathHash is used when we move files from tmp to out and there was one there before; that's
// the only case in which the hash of a filepath could change.
func (p *pathHasher) MovePathHash(oldPath, newPath string, copy bool) {
	oldPath = p.ensureRelative(oldPath)
	newPath = p.ensureRelative(newPath)
	p.pathHashMutex.Lock()
	defer p.pathHashMutex.Unlock()
	p.pathHashMemoizer[newPath] = p.pathHashMemoizer[oldPath]
	// If the path is in plz-out/tmp we aren't ever going to use it again, so free some space.
	if !copy && strings.HasPrefix(oldPath, core.TmpDir) {
		delete(p.pathHashMemoizer, oldPath)
	}
}

// ensureRelative ensures a path is relative to the repo root.
// This is important for getting best performance from memoizing the path hashes.
func (p *pathHasher) ensureRelative(path string) string {
	if strings.HasPrefix(path, core.RepoRoot) {
		return strings.TrimLeft(strings.TrimPrefix(path, core.RepoRoot), "/")
	}
	return path
}

// fileHash calculates the hash of a single file
func (p *pathHasher) fileHash(h *hash.Hash, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(*h, file)
	file.Close()
	return err
}
