package main

import (
	"fmt"
	"io/ioutil"
	"json2excel/logic"
	"os"
	"path/filepath"
)

func main() {
	l := len(os.Args)
	if l < 2 {
		fmt.Println("params lost... eg: ./main data.json")
		return
	}

	workdir := filepath.Dir(os.Args[0])
	reqFile := os.Args[1]

	// 获取json文件数据
	filePath := workdir + "/" + reqFile
	jsBytes, err := ReadReq(filePath)
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

func ReadReq(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	//req := new(proto.CheckPolicyRequest)

	js, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}
	//err = json.Unmarshal(js, req)

	return js, err
}
