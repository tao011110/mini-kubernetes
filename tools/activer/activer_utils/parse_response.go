package activer_utils

import (
	"encoding/json"
	"github.com/thedevsaddam/gojsonq/v2"
)

func GetPartOfJsonResponce(reg string, response interface{}) string {
	//只支持如下两种
	if reg == "$" {
		bytes, _ := json.Marshal(response)
		return string(bytes)
	}
	//"$.level1.level2...."
	part := string(([]byte(reg))[2:])
	bytes, _ := json.Marshal(gojsonq.New().FromInterface(response).Find(part))
	return string(bytes)
}
