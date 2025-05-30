package file

import (
	"github.com/cavaliergopher/grab/v3"
)

// Download a given link and save it to dst
func Download(link, dst string) (size int64, err error) {
	rsp, err := grab.Get(dst, link)
	if err != nil {
		return 0, err
	}
	return rsp.Size(), nil
}
