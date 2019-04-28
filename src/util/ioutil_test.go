package util

import (
	"log"
	"testing"
)

const path = "G:/jike/42_徐飞_技术与商业案例解读/118 Dremio:在Drill和Arrow上的大数据公司?.html"

func TestWriteFile(t *testing.T) {
	name := FilterFileName(path)
	WriteFile("111", name)
}

func TestFilterFileName2(t *testing.T) {
	name := FilterFileName(path)
	log.Println(name)
}
