package util

import (
	"github.com/pkg/errors"
	"log"
	"os"
)

func WriteFile(content string, path string) {
	var f *os.File

	_, err := os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			log.Fatalln(err)
		}

		f, err = os.Create(path)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		if os.IsNotExist(err) {
			f, err = os.Create(path)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	if f == nil {
		log.Fatalln(errors.Errorf("title:%s 不能生成文件", path))
	}
	defer f.Close()

	f.Write([]byte(content))
}
