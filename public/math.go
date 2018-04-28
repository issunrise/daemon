package public

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func CreateFile(fileName, data string) {
	if checkFileIsExist(fileName) {
		if f, err := os.OpenFile(fileName, os.O_APPEND, 0666); err != nil {
			log.Println("open file error:", err)
		} else if _, err := io.WriteString(f, data); err != nil {
			log.Println("write file error:", err)
		}
	} else if err := ioutil.WriteFile(fileName, []byte(data), 0666); err != nil {
		log.Println("create file error:", err)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func NewDir(dirName string) error {
	var path string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
	dirNames := strings.Split(dirName, path)
	// log.Println("name", dirNames)
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, name := range dirNames {
		dir = dir + path + name
		// log.Println("start", dir)
		if checkFileIsExist(dir) {
			// log.Println("checkFileIsExist", dir)
			continue
		}
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			log.Println("Mkdir", dir, err)
			return err
		}
	}
	return nil
}

func StringToFloat(str string) float64 {
	var value float64
	fmt.Sscanf(str, "%f", &value)
	return value
}
