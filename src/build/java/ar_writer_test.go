package java

import (
	"io"
	"os"
	"testing"

	"github.com/blakesmith/ar"
	"github.com/stretchr/testify/assert"
)

func TestMissingPath(t *testing.T) {
	assert.Error(t, CombineAr("test_missing_path.a", "doesnt_exist", nil, nil))
}

func TestCombineArFiles(t *testing.T) {
	assert.NoError(t, CombineAr("test_combine.a", "src/build/java/test_data", []string{".a"}, []string{".x.a"}))
	// Read the file back and check the contents match. Crucially we should have duplicate file names.
	f, err := os.Open("test_combine.a")
	assert.NoError(t, err)
	defer f.Close()
	contents := []string{}
	r := ar.NewReader(f)
	for {
		hdr, err := r.Next()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		contents = append(contents, hdr.Name)
	}
	expected := []string{"test1.txt", "test2.txt", "test1.txt", "test2.txt"}
	assert.Equal(t, expected, contents)
}
