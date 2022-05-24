package def

//payload = {'jobName': jobName, 'result': result, 'error': error}

type GPUJobResponse struct {
	JobName string `json:"jobName"`
	Result  string `json:"result"`
	Error   string `json:"error"`
}
