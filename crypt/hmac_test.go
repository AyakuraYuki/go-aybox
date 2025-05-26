package crypt

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testHmacSrc() []byte { return []byte("apple") }
func testHmacKey() []byte { return []byte("secret") }

func TestHmacSHA1(t *testing.T) {
	dst := HmacSHA1(testHmacSrc(), testHmacKey())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "2651783bdc7367acd2dde6f830ca0b7104368911", got)
}

func TestHmacSHA256(t *testing.T) {
	dst := HmacSHA256(testHmacSrc(), testHmacKey())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "37431003b2d14b6bddb9334c7ec2ff0ea0c65f96ec650952384e56cae83c398f", got)
}

func TestHmacSHA256224(t *testing.T) {
	dst := HmacSHA256224(testHmacSrc(), testHmacKey())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "6612b1a347a167b10d664953966cb70c33c83ff4e8abc1413ebdc497", got)
}

func TestHmacSHA512(t *testing.T) {
	dst := HmacSHA512(testHmacSrc(), testHmacKey())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "33c2f1dbd0b93a8a8354ddb888df1ff97b986959d4d710280f66730a913dc9d4535c43a3d51b3c7ff3708355d3d75ab67a105221b8ca803ed4e604f13514b145", got)
}

func TestHmacSHA512224(t *testing.T) {
	dst := HmacSHA512224(testHmacSrc(), testHmacKey())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "e3ec95df70bfb6b8cb830d070266989303ae1de95c8403eb4e24edf6", got)
}

func TestHmacSHA512256(t *testing.T) {
	dst := HmacSHA512256(testHmacSrc(), testHmacKey())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "398ffae8d29360168f32850ccd1c0f2e09c93b1d9d8cbb7e96087f7d9a7d347d", got)
}

func TestHmacSHA384(t *testing.T) {
	dst := HmacSHA384(testHmacSrc(), testHmacKey())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "d2a847a311701d5871cbd5871f53baf23a038763836acfe4314897c9080b9af7cf349e8195bde1da11a517454ee22ddb", got)
}
