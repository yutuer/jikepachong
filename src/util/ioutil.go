package util

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"regexp"
)

func createFile(path string) (*os.File, error) {
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
	return f, err
}

func WriteFile(content string, path string) {
	f, err := createFile(path)
	if f == nil {
		log.Println(errors.Errorf("%s 不能生成文件", path))
		log.Println(err.Error())
		return
	}
	defer f.Close()

	f.Write([]byte(content))
}

func WriteFile_B(bs []byte, path string) {
	f, err := createFile(path)
	if f == nil {
		log.Println(errors.Errorf("%s 不能生成文件", path))
		log.Println(err.Error())
		return
	}
	defer f.Close()

	f.Write(bs)
}

var re = regexp.MustCompile("[/:?|]")

func FilterFileName(oldName string) string {
	title := oldName
	s := re.ReplaceAllString(title, "_")
	return s
}
