package crypt

func testSrc() []byte { return []byte("hello, world") }

func testDesIV() []byte   { return []byte("56785678") }
func testDesKey() []byte  { return []byte("12341234") }
func test3DesKey() []byte { return []byte("123456781234567812345678") }

func testAesIV() []byte     { return []byte("5678567856785678") }
func testAes128Key() []byte { return []byte("1234567812345678") }
func testAes192Key() []byte { return []byte("123456781234567812345678") }
func testAes256Key() []byte { return []byte("12345678123456781234567812345678") }
