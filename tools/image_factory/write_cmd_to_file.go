package image_factory

import (
	"fmt"
	"os"
)

func WriteCmdToFile(filePath string, cmd string) {
	file, err := os.Create(filePath)

	//err := os.Truncate(filePath, 0)
	if err != nil {
		fmt.Println(err)
	}
	//file, _ := os.OpenFile(filePath, os.O_RDWR, os.ModeAppend)
	_, err = file.Write([]byte(cmd))
	if err != nil {
		fmt.Println(err)
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
}
