package file_hex

import (
	"fmt"
	"os"
	"strings"
)

func Print(path string) {
	content, _ := ReadTo(path)
	fmt.Println(content)
}

func Head(path string, n uint) {
	content, _ := ReadTo(path)
	lines := strings.Split(content, "\n")
	if uint(len(lines)) < n+2 {
		fmt.Println(content)
		return
	}
	for _, line := range lines[:n+2] {
		fmt.Println(line)
	}
	fmt.Println("...")
}

func Tail(path string, n uint) {
	content, _ := ReadTo(path)
	lines := strings.Split(content, "\n")
	count := uint(len(lines))
	if count <= n+2 {
		fmt.Println(content)
		return
	}
	fmt.Println(lines[0])
	fmt.Println(lines[1])
	fmt.Println("...")
	for _, line := range lines[count-n-1 : count-1] {
		fmt.Println(line)
	}
}

func ReadTo(path string) (content string, err error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}

	bytesPerLine := 16 // number of bytes to be printed per line

	builder.WriteString("          00 01 02 03  04 05 06 07  08 09 0A 0B  0C 0D 0E 0F   |  Decoded Text  |\n")
	builder.WriteString("---------------------------------------------------------------------------------\n")

	for i := 0; i < len(bs); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(bs) {
			end = len(bs)
		}

		// print address offset
		builder.WriteString(fmt.Sprintf("%08x  ", i))

		// print hex parts
		for j := i; j < end; j++ {
			builder.WriteString(fmt.Sprintf("%02x ", bs[j]))
			if (j-i)%4 == 3 {
				builder.WriteRune(' ') // add whitespace every 4 bytes
			}
		}

		// align last line
		if end < i+bytesPerLine {
			spaces := (bytesPerLine - (end - i)) * 3
			spaces += ((bytesPerLine - (end - i)) + 3) / 4
			builder.WriteString(fmt.Sprintf("%*s", spaces, ""))
		}

		builder.WriteString(" |")

		// print human-readable ascii characters
		for j := i; j < end; j++ {
			if bs[j] >= 32 && bs[j] <= 126 {
				builder.WriteString(fmt.Sprintf("%c", bs[j]))
			} else {
				builder.WriteRune('.')
			}
		}

		builder.WriteString("|\n")
	}

	return builder.String(), nil
}
