package main

import (
	"fmt"
	"json2excel/common"
	"json2excel/log"
	"json2excel/logic"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	l := len(os.Args)
	if l < 2 {
		fmt.Println("params lost... eg: ./main data.json")
		return
	}

	workdir := filepath.Dir(os.Args[0])
	reqFile := os.Args[1]

	log_path := workdir + "/main.log"
	log.InitLogger(log_path, "info")
	log.Infof("deal request start for command: %s", strings.Join(os.Args, " "))

	// 获取json文件数据
	filePath := workdir + "/" + reqFile
	jsBytes, err := common.ReadFile(filePath)
	if err != nil {
		log.Errorf("ReadJsonFile err: %s", err.Error())
		fmt.Println("ReadJsonFile err: ", err)
		return
	}

	json2Excel := new(logic.Json2Excel)

	savePath, err := json2Excel.Json2Excel(jsBytes, workdir)
	if err != nil {
		log.Errorf("Json2Excel err: %s", err.Error())
		fmt.Println("Json2Excel err: ", err)
		return
	}

	log.Infof("Json2Excel success excelPath: %s", savePath)
	fmt.Println("Json2Excel success excelPath: ", savePath)
}
