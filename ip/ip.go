package ip

import (
	"os"

	"github.com/AyakuraYuki/go-aybox/file"
)

func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return file.ReadTextFile("/etc/hostname")
	}
	return name
}
