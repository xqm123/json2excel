package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
)

func ReadFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	js, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}

	return js, err
}

func TransCellVal(val interface{}) (v interface{}) {
	valRef := reflect.ValueOf(val)
	if !valRef.IsValid() {
		v = nil
	} else if valRef.Kind() == reflect.Slice || valRef.Kind() == reflect.Array || valRef.Kind() == reflect.Map {
		valjs, _ := json.Marshal(val)
		v = string(valjs)
	} else if valRef.Kind() == reflect.Bool {
		v = strconv.FormatBool(val.(bool))
	} else {
		v = val
	}

	return v
}
