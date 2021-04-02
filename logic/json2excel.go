package logic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
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
		return "", err
	}
	if len(jsonMap) < 1 {
		return "", fmt.Errorf("json data count lt 1")
	}

	header := make([]string, 0)
	temp := map[string]struct{}{}

	for _, v := range jsonMap {
		for km, _ := range v {
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
	fmt.Println("eIndexX:", eIndexX)
	indexY := 2
	sheetName := "Sheet1"

	//设置列宽
	if err = f.SetColWidth(sheetName, sIndexX, eIndexX, 20); err != nil {
		return "", err
	}

	//设置header行高
	if err = f.SetRowHeight("Sheet1", indexY, 20); err != nil {
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

	err = f.SetSheetRow(sheetName, fmt.Sprintf("%s%d", sIndexX, indexY), &header)
	if err != nil {
		return "", err
	}
	indexY++

	//设置值
	values := make([]interface{}, len(header))
	for _, row := range jsonMap {
		for i, vm := range header {
			if val, ok := row[vm]; ok {
				valRef := reflect.ValueOf(val)
				if !valRef.IsValid() {
					val = nil
				} else if valRef.Kind() == reflect.Slice || valRef.Kind() == reflect.Array || valRef.Kind() == reflect.Map {
					valjs, _ := json.Marshal(val)
					val = string(valjs)
				} else if valRef.Kind() == reflect.Bool {
					val = strconv.FormatBool(val.(bool))
				}
				values[i] = val
			} else {
				values[i] = nil
			}
		}

		err = f.SetSheetRow(sheetName, fmt.Sprintf("%s%d", sIndexX, indexY), &values)
		if err != nil {
			return "", err
		}
		indexY++
	}

	savePath = fmt.Sprintf("%s/%s.xlsx", saveDir, time.Now().Format("20060102150405"))
	if err := f.SaveAs(savePath); err != nil {
		return "", err
	}

	return savePath, nil
}
