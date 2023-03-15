package utils

import "os"

func CreateFile() (*os.File, error) {
	return os.Create("../Files/File1.pdf")
}
