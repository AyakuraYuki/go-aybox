package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5ToString(t *testing.T) {
	src := "hello, world"
	dst := MD5ToString(src)
	assert.EqualValues(t, "e4d7f1b4ed2e42d15898f4b27b019da4", dst)
}

func TestBinaryMD5ToString(t *testing.T) {
	bs, err := testAssetFS.ReadFile("test/sample-taxi.jpg")
	assert.NoError(t, err)
	dst := BinaryMD5ToString(bs)
	assert.EqualValues(t, "328c8ac7fde8eb0f309a772d8fb27fab", dst)
}
