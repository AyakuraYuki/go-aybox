package crypt

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSHASrc() []byte { return []byte("apple") }

func TestSHA1(t *testing.T) {
	dst := SHA1(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "d0be2dc421be4fcd0172e5afceea3970e2f3d940", got)
}

func TestSHA256(t *testing.T) {
	dst := SHA256(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "3a7bd3e2360a3d29eea436fcfb7e44c735d117c42d1c1835420b6b9942dd4f1b", got)
}

func TestSHA224(t *testing.T) {
	dst := SHA224(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "b7bbfdf1a1012999b3c466fdeb906a629caa5e3e022428d1eb702281", got)
}

func TestSHA512(t *testing.T) {
	dst := SHA512(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "844d8779103b94c18f4aa4cc0c3b4474058580a991fba85d3ca698a0bc9e52c5940feb7a65a3a290e17e6b23ee943ecc4f73e7490327245b4fe5d5efb590feb2", got)
}

func TestSHA512224(t *testing.T) {
	dst := SHA512224(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "62882a192e3b8956e3f42a3accb9e526e0edc3ccb7dde2d997b0e0ae", got)
}

func TestSHA512256(t *testing.T) {
	dst := SHA512256(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "1f0f662e561e1d4abf27f945db3c5a8305006138161aad4b9933c4a02964ee54", got)
}

func TestSHA384(t *testing.T) {
	dst := SHA384(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "3d8786fcb588c93348756c6429717dc6c374a14f7029362281a3b21dc10250ddf0d0578052749822eb08bc0dc1e68b0f", got)
}

func TestSHA3(t *testing.T) {
	dst := SHA3(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "3f932042b1206955946ffc1095c46b75a6d45257b579476ffea0139348fb14070afeeee5de8250000683c0336e00b52e965a1c5263743e8509eb713818864411", got)
}

func TestSHA3224(t *testing.T) {
	dst := SHA3224(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "a8a4edebe16de798cdddb6cf3bb2e24da900335d9111dc62c6c2b7bd", got)
}

func TestSHA3256(t *testing.T) {
	dst := SHA3256(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "42a990655bffe188c9823a2f914641a32dcbb1b28e8586bd29af291db7dcd4e8", got)
}

func TestSHA3384(t *testing.T) {
	dst := SHA3384(testSHASrc())
	got := hex.EncodeToString(dst)
	assert.EqualValues(t, "6063b43ad451d2b7363d955cd3eb41c19dd7ba5146f20ce72de696c134422035286f4db01a99c6c34cc3a3271b7e3042", got)
}
