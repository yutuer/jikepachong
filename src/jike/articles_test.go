package jike

import (
	"log"
	"os"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

const (
	fileName    = "tmpFile.html"
	newFileName = "10讲拒绝重复，你的代码百宝箱：如何书写codesnippet.html"
)

func TestCreateFile(t *testing.T) {
	f := newFileName

	if _, e := os.Stat(f); os.IsNotExist(e) {
		_, e := os.Create(f)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestExistFile(t *testing.T) {
	if _, e := os.Stat(newFileName); e == nil {
		os.Remove(newFileName)
	}
}

func TestRenameFile(t *testing.T) {
	if _, e := os.Stat(fileName); e == nil {
		e = os.Rename(fileName, newFileName)
		if e != nil {
			t.Fatal(e)
		}
	}
}
