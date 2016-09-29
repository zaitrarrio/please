package java

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/blakesmith/ar"
)

// Combines a sequence of ar files into one.
func CombineAr(outFile, inDir, suffix, excludeSuffix string) error {
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
		if err != nil || info.IsDir() || !strings.HasSuffix(path, suffix) || strings.HasSuffix(path, excludeSuffix) {
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
	for {
		hdr, err := r.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// For unknown reasons they always seem to end in / unnecessarily. Strip it off.
		hdr.Name = strings.TrimSuffix(hdr.Name, "/")
		log.Info("Adding %s from %s", hdr.Name, path)
		if err := w.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			return err
		}
	}
	return nil
}
