package logic

import (
	"encoding/json"
	"fmt"
	"json2excel/common"
	"json2excel/log"
	"sort"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// github.com/360EntSecGroup-Skylar/excelize/v2
// document: https://xuri.me/excelize/zh-hans/

type Json2Excel struct {
}

func (j *Json2Excel) Json2Excel(jsonBytes []byte, saveDir string) (savePath string, err error) {
	jsonMap := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		log.Errorf("json.Unmarshal error. err: %s", err.Error())
		return "", err
	}
	if len(jsonMap) < 1 {
		log.Error("json data count lt 1")
		return "", fmt.Errorf("json data count lt 1")
	}

	header := make([]string, 0)
	temp := map[string]struct{}{}

	for _, v := range jsonMap {
		for km := range v {
			if _, ok := temp[km]; !ok {
				temp[km] = struct{}{}
				header = append(header, km)
			}
		}
	}

	sort.Strings(header)

	f := excelize.NewFile()

	sIndexX := "B"
	eIndexX := string(byte(int(sIndexX[0]) + len(header) - 1))
	indexY := 2
	sheetName := "Sheet1"

	//设置列宽
	if err = f.SetColWidth(sheetName, sIndexX, eIndexX, 20); err != nil {
		log.Errorf("f.SetColWidth error. err: %s", err.Error())
		return "", err
	}

	//设置header行高
	if err = f.SetRowHeight(sheetName, indexY, 20); err != nil {
		log.Errorf("f.SetRowHeight error. err: %s", err.Error())
		return "", err
	}

	//设置头
	headerStyle, err := f.NewStyle(`{
		"font":
		{
			"bold": true,
			"size": 14
		},
		"alignment":
		{
			"horizontal": "left",
			"vertical": "center"
		}
	}`)
	if err != nil {
		return "", err
	}

	err = f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", sIndexX, indexY), fmt.Sprintf("%s%d", eIndexX, indexY), headerStyle)
	if err != nil {
		log.Errorf("f.SetCellStyle error. err: %s", err.Error())
		return "", err
	}

	err = f.SetSheetRow(sheetName, fmt.Sprintf("%s%d", sIndexX, indexY), &header)
	if err != nil {
		log.Errorf("f.SetSheetRow error. err: %s", err.Error())
		return "", err
	}
	indexY++

	//设置值
	values := make([]interface{}, len(header))
	for _, row := range jsonMap {
		for i, vm := range header {
			if val, ok := row[vm]; ok {
				values[i] = common.TransCellVal(val)
			} else {
				values[i] = nil
			}
		}

		err = f.SetSheetRow(sheetName, fmt.Sprintf("%s%d", sIndexX, indexY), &values)
		if err != nil {
			log.Errorf("f.SetSheetRow error. err: %s", err.Error())
			return "", err
		}
		indexY++
	}

	savePath = fmt.Sprintf("%s/%s.xlsx", saveDir, time.Now().Format("20060102150405"))
	if err := f.SaveAs(savePath); err != nil {
		log.Errorf("f.SaveAs error. savePath: %s, err: %s", savePath, err.Error())
		return "", err
	}

	return savePath, nil
}
