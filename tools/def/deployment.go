package def

// ParsedDeployment
//下面不是deployment解析的直接结果
//为了复用pod相关的接口, 后请把template中的内容转为pod并为此pod分配全局唯一的name
type ParsedDeployment struct {
	Name        string `json:"name"`
	ReplicasNum int    `json:"replicas_num"`
	PodName     string `json:"pod_name"`
}
