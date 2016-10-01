package java

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/peterebden/ar"
)

// Combines a sequence of ar files into one.
func CombineAr(outFile, inDir string, suffix, excludeSuffix []string) error {
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()
	w := ar.NewWriter(f)
	// We have to write an initial index, so we have to defer writing the actual contents until then.
	// We could also read them twice if we wanted to save memory, but the code this way is simpler.
	headers := []*ar.Header{}
	contents := [][]byte{}
	if err := w.WriteGlobalHeader(); err != nil {
		return err
	}
	if err := filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == outFile || info.IsDir() || !MatchesSuffix(path, suffix) || MatchesSuffix(path, excludeSuffix) {
			return err
		}
		log.Notice("Adding %s", path)
		h, c, err := addArFile(path)
		headers = append(headers, h...)
		contents = append(contents, c...)
		return err
	}); err != nil {
		return err
	}
	// Now write the the filename index
	index, indexContents := deriveIndex(headers)
	log.Notice("Writing file index (%d bytes)", len(indexContents))
	if err := w.WriteHeader(index); err != nil {
		return err
	}
	if _, err := io.Copy(w, bytes.NewReader(indexContents)); err != nil {
		return err
	}
	// And the actual contents
	for i, hdr := range headers {
		log.Info("Writing %s (%d / %d)", hdr.Name, hdr.Size, len(contents[i]))
		if err := w.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := io.Copy(w, bytes.NewReader(contents[i])); err != nil {
			return err
		}
	}
	return nil
}

func addArFile(path string) ([]*ar.Header, [][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	r := ar.NewReader(f)
	headers := []*ar.Header{}
	contents := [][]byte{}
	var filenames string
	for {
		hdr, err := r.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}
		// Handle gcc index entries
		if hdr.Name == "/" {
			continue // This is the index, we will get ar to regenerate it.
		} else if hdr.Name == "//" {
			// This is an index of filenames that are too long
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, r); err != nil {
				return nil, nil, err
			}
			filenames = buf.String()
			continue
		}
		// Handle "filenames" that are an index into the filenames array
		if strings.HasPrefix(hdr.Name, "/") {
			i, err := strconv.Atoi(strings.TrimPrefix(hdr.Name, "/"))
			if err != nil {
				return nil, nil, err // Not sure, maybe we should continue?
			}
			hdr.Name = filenames[i:]
			hdr.Name = hdr.Name[:strings.IndexRune(hdr.Name, '/')]
		} else {
			// For unknown reasons they always seem to end in / unnecessarily. Strip it off.
			hdr.Name = strings.TrimSuffix(hdr.Name, "/")
		}
		// Normalise all mod times
		hdr.ModTime = time.Unix(0, 0)
		log.Info("Adding %s from %s, mode %d, size %d", hdr.Name, path, hdr.Mode, hdr.Size)
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return nil, nil, err
		}
		headers = append(headers, hdr)
		contents = append(contents, buf.Bytes())
	}
	return headers, contents, nil
}

// MatchesSuffix returns true if the given path matches any one of the given suffixes.
func MatchesSuffix(path string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if suffix != "" && strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}

func deriveIndex(headers []*ar.Header) (*ar.Header, []byte) {
	var buf bytes.Buffer
	for _, hdr := range headers {
		newName := "/" + strconv.Itoa(buf.Len())
		buf.WriteString(hdr.Name + "/\n")
		hdr.Name = newName
	}
	return &ar.Header{
		Name: "//",
		Size: int64(buf.Len()),
	}, buf.Bytes()
}
