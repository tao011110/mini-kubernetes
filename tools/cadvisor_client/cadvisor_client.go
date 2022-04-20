package cadvisor_client

import "github.com/google/cadvisor/client"

func fuck() {
	client, err := client.NewClient("http://localhost:8080/")
	mInfo, err := client.MachineInfo()
}
