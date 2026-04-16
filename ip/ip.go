package ip

import (
	"os"

	"github.com/AyakuraYuki/go-aybox/files"
)

func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return files.ReadTextFile("/etc/hostname")
	}
	return name
}
