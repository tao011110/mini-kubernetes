package activer_utils

import (
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
)

func GetPartOfJsonResponce(reg string, response string) string {
	//只支持如下两种
	if reg == "$" {
		return response
	}
	//"$.level1.level2...."
	return fmt.Sprintf("%v", gojsonq.New().FromString(response).Find(string(([]byte(response))[2:])))
}
