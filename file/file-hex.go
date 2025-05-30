package file

import (
	"fmt"
	"os"
)

func ShowInHex(path string) {
	bs, err := os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	bytesPerLine := 16 // number of bytes to be printed per line

	fmt.Println("          00 01 02 03  04 05 06 07  08 09 0A 0B  0C 0D 0E 0F   |  Decoded Text  |")
	fmt.Println("---------------------------------------------------------------------------------")

	for i := 0; i < len(bs); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(bs) {
			end = len(bs)
		}

		// print address offset
		fmt.Printf("%08x  ", i)

		// print hex parts
		for j := i; j < end; j++ {
			fmt.Printf("%02x ", bs[j])
			if (j-i)%4 == 3 {
				fmt.Print(" ") // add whitespace every 4 bytes
			}
		}

		// align last line
		if end < i+bytesPerLine {
			spaces := (bytesPerLine - (end - i)) * 3
			spaces += ((bytesPerLine - (end - i)) + 3) / 4
			fmt.Printf("%*s", spaces, "")
		}

		fmt.Print(" |")

		// print human-readable ascii characters
		for j := i; j < end; j++ {
			if bs[j] >= 32 && bs[j] <= 126 {
				fmt.Printf("%c", bs[j])
			} else {
				fmt.Print(".")
			}
		}

		fmt.Println("|")
	}
}
