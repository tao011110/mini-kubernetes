package gpu_job_api

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/apiserver/delete_api"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/util"
	"os"
)

func OutputGPUJOBResponse(li *clientv3.Client, response def.GPUJobResponse) {
	gpuJob := apiserver_utils.GetGPUJobByName(li, response.JobName)
	resultFilePath := gpuJob.ResultPath + "/result.out"
	errorFilePath := gpuJob.ResultPath + "/error.err"
	resultFile, _ := os.Create(resultFilePath)
	errorFile, _ := os.Create(errorFilePath)
	defer func(resultFile *os.File) {
		err := resultFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resultFile)
	defer func(errorFile *os.File) {
		err := errorFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(errorFile)
	_, err := resultFile.WriteString(response.Result)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = errorFile.WriteString(response.Error)
	if err != nil {
		fmt.Println(err)
		return
	}
	nameList := apiserver_utils.GetPodReplicaListByPodName(li, gpuJob.PodName)
	if len(nameList) != 1 {
		panic("job's pod replica name list length!=1")
	}
	podInstanceID := nameList[0]
	podInstance := apiserver_utils.GetPodInstanceByID(li, podInstanceID)
	podInstance.Status = def.SUCCEEDED
	util.PersistPodInstance(podInstance, li)
	delete_api.DeletePod(li, podInstanceID)
}
