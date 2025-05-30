package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathExists(t *testing.T) {
	pwd, _ := os.Getwd()
	assert.True(t, PathExist(filepath.Join(pwd, "test", "google.svg")))
	assert.False(t, PathExist(filepath.Join(pwd, "test", "non_exist_file")))
}

func TestListFileNoRecursive(t *testing.T) {
	pwd, _ := os.Getwd()
	files, err := ListFileNoRecursive(pwd)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, len(files))
}

func TestListFiles(t *testing.T) {
	pwd, _ := os.Getwd()
	files := ListFile(pwd)
	assert.EqualValues(t, 11, len(files))
}

func TestIsDir(t *testing.T) {
	pwd, _ := os.Getwd()
	assert.True(t, IsDir(filepath.Join(pwd, "test")))
	assert.False(t, IsDir(filepath.Join(pwd, "file.go")))
	assert.False(t, IsDir(filepath.Join(pwd, "non_exist_dir")))
}

func TestIsFile(t *testing.T) {
	pwd, _ := os.Getwd()
	assert.False(t, IsFile(filepath.Join(pwd, "test")))
	assert.True(t, IsFile(filepath.Join(pwd, "file.go")))
	assert.False(t, IsFile(filepath.Join(pwd, "non_exist_file.go")))
}
