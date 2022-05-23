package gpu_job_uploader

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/util"
)

func apiServerIPFileGenerator() string {
	//tempFilePath := def.TemplateFileDir + goid.NewV4UUID().String()
	//fp, err := os.Create(tempFilePath)
	//if err != nil {
	//	return ""
	//}
	//defer func(fp *os.File) {
	//	_ = fp.Close()
	//}(fp)
	//_, _ = fp.WriteString(fmt.Sprintf("%s:%d\n", util.GetLocalIP().String(), def.MasterPort))
	//return tempFilePath
	return fmt.Sprintf("%s:%d\n", util.GetLocalIP().String(), def.MasterPort)
}
