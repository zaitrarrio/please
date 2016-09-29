package java

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/blakesmith/ar"
)

// Combines a sequence of ar files into one.
func CombineAr(outFile, inDir string, suffix, excludeSuffix []string) error {
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()
	w := ar.NewWriter(f)
	if err := w.WriteGlobalHeader(); err != nil {
		return err
	}
	if err := filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == outFile || info.IsDir() || !MatchesSuffix(path, suffix) || MatchesSuffix(path, excludeSuffix) {
			return err
		}
		log.Notice("Adding %s", path)
		return addArFile(w, path)
	}); err != nil {
		return err
	}
	return nil
}

func addArFile(w *ar.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := ar.NewReader(f)
	var filenames []byte
	for {
		hdr, err := r.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// Handle gcc index entries
		if hdr.Name == "/" {
			continue // This is the index, we will get ar to regenerate it.
		} else if hdr.Name == "//" {
			// This is some sort of index of filenames, we need to keep them for later.
			filenames = make([]byte, hdr.Size)
			if err := io.Copy(bytes.NewWriter(filenames), r); err != nil {
				return err
			}
			continue
		}
		// Handle "filenames" that are an index into the filenames array
		if strings.HasPrefix(hdr.Name, "/") {
			i, err := strconv.Atoi(strings.TrimPrefix(hdr.Name))
			if err != nil {
				return err // Not sure, maybe we should continue?
			}

		}

		// For unknown reasons they always seem to end in / unnecessarily. Strip it off.
		hdr.Name = strings.TrimSuffix(hdr.Name, "/")
		log.Info("Adding %s from %s, mode %d", hdr.Name, path, hdr.Mode)
		if err := w.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			return err
		}
	}
	return nil
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
