package crypt

import "embed"

//go:embed test
var testAssetFS embed.FS

func testSrc() []byte { return []byte("hello, world") }
