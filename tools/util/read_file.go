package util

import (
	"io/ioutil"
	"os"
)

func ReadFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	f, _ := ioutil.ReadAll(file)
	str := string(f)
	return str
}
