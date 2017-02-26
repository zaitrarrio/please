package build

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"core"
)

// N.B. The expected hashes given here and below are determined empirically - their
//      values are not based on some underlying truth, but it also cannot change
//      casually since that may break any users who specify hashes in BUILD files.
const sha1Hash = "RwNQdxGYd93Aa/trY7QKgX4e+/0"

var content = []byte("testing testing 1 2 3")

func TestPathHash(t *testing.T) {
	const path = "test_path_hash.txt"
	err := ioutil.WriteFile(path, content, 0644)
	assert.NoError(t, err)
	b, err := hasher.PathHash(path, false, true)
	assert.NoError(t, err)
	assert.Equal(t, sha1Hash, b64(b))
	// Redo the above with SHA256
	b, err = hasher.PathHash(path, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "ZW1UsOowRzFm4TjpCmV8DY/VM4ql6omnGwh6kgaCOnk", b64(b))
	// If we write a new file and force a recalculation, we should get a new hash.
	err = ioutil.WriteFile(path, []byte("testing testing 1 2 4"), 0644)
	assert.NoError(t, err)
	b, err = hasher.PathHash(path, true, false)
	assert.Equal(t, "x6sKrKbXC5mJT+3lxYNbIheqjTNx9BccY9YkCuI5iVM", b64(b))
}

func TestTargetPathHashSHA256(t *testing.T) {
	const path = "test_target_path_hash_sha256.txt"
	target := core.NewBuildTarget(core.ParseBuildLabel("//src/build", ""))
	err := ioutil.WriteFile(path, content, 0644)
	assert.NoError(t, err)
	b, err := hasher.TargetPathHash(path, target)
	assert.NoError(t, err)
	assert.Equal(t, "NUIPkeGJA7cerXW+L/grrtnDHxR9LtUejMeA6VRU3jE", b64(b))
}

func TestTargetPathHashSHA1(t *testing.T) {
	const path = "test_target_path_hash_sha1.txt"
	target := core.NewBuildTarget(core.ParseBuildLabel("//src/build", ""))
	target.Hashes = append(target.Hashes, "sha1: "+sha1Hash)
	err := ioutil.WriteFile(path, content, 0644)
	assert.NoError(t, err)
	b, err := hasher.TargetPathHash(path, target)
	assert.NoError(t, err)
	assert.Equal(t, sha1Hash, b64(b))
}

func TestMustPathHash(t *testing.T) {
	assert.Panics(t, func() { hasher.MustPathHash("test_must_path_hash.txt", false) })
}

func TestMustTargetPathHash(t *testing.T) {
	target := core.NewBuildTarget(core.ParseBuildLabel("//src/build", ""))
	target.Hashes = append(target.Hashes, "sha1: "+sha1Hash)
	assert.Panics(t, func() { hasher.MustTargetPathHash("test_must_target_path_hash.txt", target) })
}

func TestMovePathHash(t *testing.T) {
	const path = "test_move_path_hash.txt"
	const path2 = "test_move_path_hash_2.txt"
	err := ioutil.WriteFile(path, content, 0644)
	assert.NoError(t, err)
	b := hasher.MustPathHash(path, false)
	assert.Equal(t, "5c8DBLfyeedGorSLH4aJrXtlo2d9q/Ke742hygM2HE4", b64(b))
	hasher.MovePathHash(path, path2, true)
	b = hasher.MustPathHash(path2, false)
	assert.Equal(t, "5c8DBLfyeedGorSLH4aJrXtlo2d9q/Ke742hygM2HE4", b64(b))
}

func TestPathHashDir(t *testing.T) {
	b, err := hasher.PathHash("src/build/test_data/package1", false, false)
	assert.NoError(t, err)
	assert.Equal(t, "Kho4JL0cSPXcBlfzgxH7/UV+1ZOVRgXZmJqm0EsDp6M", b64(b))
}

func TestSymlink(t *testing.T) {
	// We should be able to handle symlinks, although they should produce a different hash
	// to normal files.
	// If it turns out that symlinks are no longer relevant to us, this test might go away.
	const path = "test_symlink.txt"
	const path2 = "test_symlink_2.txt"
	err := ioutil.WriteFile(path, content, 0644)
	assert.NoError(t, err)
	err = os.Symlink(path, path2)
	assert.NoError(t, err)
	b, err := hasher.PathHash(path2, false, false)
	assert.NoError(t, err)
	assert.Equal(t, "ovS9SIbJdV01H+0pQOLVzWObZI/shDZhYLvKmd1tEkE", b64(b))
}