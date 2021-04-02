package common

import (
	"io/ioutil"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
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
