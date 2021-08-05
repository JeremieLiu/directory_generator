package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	SourcePath     = `./test_dir`
	PathTitle      = "|--"
	GidFolder      = ".git"
	OutPutFileName = "./project_folder_%d.md"
	res            = ""
)

func main() {
	result := GetAllInnerFileOrDir(SourcePath)
	outFileName := fmt.Sprintf(OutPutFileName, time.Now().Unix())
	if err := CheckOrCreateFile(outFileName); err != nil {
		fmt.Printf("[ERROR] CheckOrCreateFile error:%s \r\n", err)
	}
	if err := AppendContentToFile(outFileName, result); err != nil {
		fmt.Printf("[ERROR] AppendContentToFile error:%s \r\n", err)
	}
}

// directory utils------------------------------------------------------------------------------------------------------
// 获取文件目录下所有的文件名称
func GetAllInnerFileOrDir(rootPath string) string {
	fileInfoList, err := ioutil.ReadDir(rootPath)
	if err != nil {
		fmt.Printf("ReadDir err:%s \r\n", err)
	}

	for _, v := range fileInfoList {
		currentName := v.Name()

		// 过滤.git目录
		if currentName == GidFolder {
			continue
		}

		// 当前文件名的全路径
		currentPath := filepath.Join(rootPath, currentName)
		res = fmt.Sprintf("%s%s\n", res, fmt.Sprintf(FormatPathPattern(currentPath)))
		if IsDir(currentPath) {
			GetAllInnerFileOrDir(currentPath)
		}
	}
	return res
}

// 判断是否是文件夹
func IsDir(path string) bool {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil || len(fileInfoList) == 0 {
		return false
	}
	return true
}

//格式化文件开头
func FormatPathPattern(fullPath string) string {
	tabCount := GetCountOfFileFolder(fullPath)
	tabStr := ""
	for i := tabCount; i > 0; i-- {
		tabStr = fmt.Sprintf("%s%s", tabStr, "\t")
	}
	fileName := GetLastNameOfPath(fullPath)
	return fmt.Sprintf("%s%s %s", tabStr, PathTitle, fileName)
}

//获取文件层数
func GetCountOfFileFolder(fullPathName string) int {
	rest := strings.Replace(fullPathName, SourcePath, "", 1)
	return strings.Count(rest, string(os.PathSeparator)) - 1
}

// 全路径获取最后文件路劲
func GetLastNameOfPath(fullPathName string) string {
	slice := strings.Split(fullPathName, string(os.PathSeparator))
	length := len(slice)
	return slice[length-1]
}

// file utils-----------------------------------------------------------------------------------------------------------
func CheckOrCreateFile(fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if _, err := os.Create(fileName); err != nil {
			return err
		}
	}
	return nil
}

func AppendContentToFile(fileName, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0x666)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.WriteString(f, content)
	return err
}
