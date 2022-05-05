package kubelet

import (
	"fmt"
	"github.com/robfig/cron"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/resource"
	"mini-kubernetes/tools/util"
	"sort"
)

func ContainerCheck() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/5 * * * * *", checkPodRunning)
	if err != nil {
		fmt.Println(err)
	}

	cron2.Start()
	defer cron2.Stop()
	for {
		if node.ShouldStop {
			break
		}
	}
}

func checkPodRunning() {
	infos := resource.GetAllContainersInfo(node.CadvisorClient)
	var runningContainerIDs []string
	for _, info := range infos {
		runningContainerIDs = append(runningContainerIDs, info.Id)
	}
	sort.Strings(runningContainerIDs)
	for _, instance := range node.PodInstances {
		for _, container := range instance.ContainerSpec {
			if sort.SearchStrings(runningContainerIDs, container.ID) == len(runningContainerIDs) {
				instance.Status = def.FAILED
				util.PersistPodInstance(*instance, node.EtcdClient)
				continue
			}
		}
	}
}
