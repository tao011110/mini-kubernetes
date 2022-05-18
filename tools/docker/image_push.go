package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"mini-kubernetes/tools/def"
	"os"
)

// PushImage 需要传入函数的名称
func PushImage(funcName string) {
	image := def.RgistryAddr + funcName
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()

	// RegistryAuth is the base64 encoded credentials for the registry
	// So create a types.AuthConfig and translate it into base64
	authConfig := types.AuthConfig{Username: def.RgistryUsername, Password: def.RgistryPassword}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	resp, err := cli.ImagePush(context.Background(), image, types.ImagePushOptions{
		All:           false,
		RegistryAuth:  authStr,
		PrivilegeFunc: nil,
	})
	if err != nil {
		fmt.Printf("Pushed imageID failed %s\n %v\n", image, err)
		panic(err)
	}
	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	} else {
		fmt.Printf("Pushed imageID %s successully\n", image)
	}
}
