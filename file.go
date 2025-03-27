package aybox

import "os"

func PathExist(path string) bool {
	exist, _ := PathExistE(path)
	return exist
}

func PathExistE(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
