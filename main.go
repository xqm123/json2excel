package main

import (
	"fmt"
	"json2excel/common"
	"json2excel/logic"
	"os"
	"path/filepath"
	"runtime/debug"
)

func main() {
	l := len(os.Args)
	if l < 2 {
		fmt.Println("params lost... eg: ./main data.json")
		return
	}

	workdir := filepath.Dir(os.Args[0])
	reqFile := os.Args[1]

	var err error
	defer func() {
		//如果发生了错误，则在程序结束前打印堆栈信息
		if err != nil {
			debug.PrintStack()
		}
	}()

	// 获取json文件数据
	filePath := workdir + "/" + reqFile
	jsBytes, err := common.ReadFile(filePath)
	if err != nil {
		fmt.Println("ReadJsonFile err :", err)
		return
	}

	json2Excel := new(logic.Json2Excel)

	savePath, err := json2Excel.Json2Excel(jsBytes, workdir)
	if err != nil {
		fmt.Println("Json2Excel err:", err)
		return
	}

	fmt.Println("Json2Excel success excelPath:", savePath)
}
