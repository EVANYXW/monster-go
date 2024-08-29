// Package common @Author evan_yxw
// @Date 2024/8/27 21:17:00
// @Desc
package common

import (
	"encoding/json"
	"os"
)

var RunDirPrefix string

func init() {
	var err error
	RunDirPrefix, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func xml_loader(name string) ([]byte, error) {
	path := RunDirPrefix + "/scp/" + name + ".xml"
	data, err := os.ReadFile(path)
	if err != nil {
		//xsf_log.Error("CSVReader read file error, " + err.Error())
		return nil, err
	}

	return data, nil
}

func json_loader(name string) ([]map[string]interface{}, error) {
	path := RunDirPrefix + "/scp/" + name + ".json"
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		jsonData := make([]map[string]interface{}, 0)
		if err = json.Unmarshal(bytes, &jsonData); err != nil {
			return nil, err
		}
		return jsonData, nil
	}
}
