package util

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"regexp"
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
		log.Println(errors.Errorf("%s 不能生成文件", path))
		log.Println(err.Error())
		return
	}

	defer f.Close()

	f.Write([]byte(content))
}

var re = regexp.MustCompile("[/:?|]")

func FilterFileName(oldName string) string {
	title := oldName
	s := re.ReplaceAllString(title, "_")
	return s
}
