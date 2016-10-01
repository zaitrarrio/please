package java

import (
	"io"
	"os"
	"testing"

	"github.com/peterebden/ar"
	"github.com/stretchr/testify/assert"
)

func TestMissingPath(t *testing.T) {
	assert.Error(t, CombineAr("test_missing_path.a", "doesnt_exist", nil, nil))
}

func TestCombineGccArFiles(t *testing.T) {
	// .a files written by gcc have an index as well, and often a filename index.
	assert.NoError(t, CombineAr("test_combine_gcc.a", "src/build/java/test_data", []string{".a"}, []string{".x.a"}))
	f, err := os.Open("test_combine_gcc.a")
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
	//expected := []string{"libembedded_file_1.o", "libembedded_file_3.o"}
	expected := []string{"//", "/0", "/22"}
	assert.Equal(t, expected, contents)
}
