package image_factory

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"os"
	"os/exec"
	"strings"
)

func EchoFactory(preStr string, target string) string {
	newStr := strings.Replace(preStr, "\n", "\\n", -1)
	newStr = strings.Replace(newStr, "\\", "\\\\", -1)
	return fmt.Sprintf("echo -e \\\"%s\\\" > %s", newStr, target)
}

func WriteCmdToFile(filePath string, cmd string) {
	// 需要保证file存在且mode为777
	err := os.Truncate(filePath, 0)
	if err != nil {
		fmt.Println(err)
	}
	file, _ := os.OpenFile(filePath, os.O_RDWR, os.ModeAppend)
	_, err = file.Write([]byte(cmd))
	if err != nil {
		fmt.Println(err)
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func ImageFactory(baseImageName string, newImageName string, commandInContainer []string) {
	cmd := exec.Command("docker", "run", "-itd", "--name", newImageName, baseImageName)
	fmt.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	//container := def.Container{
	//	Image: baseImageName,
	//}
	//containerID := docker.CreateContainer(container, newImageName)
	for _, command := range commandInContainer {
		cmd := exec.Command("docker", "exec", newImageName, "/bin/bash", "-c", fmt.Sprintf("'%s'", command)).String()
		WriteCmdToFile(def.TemplateCmdFilePath, cmd)
		command := fmt.Sprintf(`%s .`, def.TemplateCmdFilePath)
		err = exec.Command("/bin/bash", "-c", command).Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	docker.CommitContainer(newImageName, newImageName)
	docker.PushImage(newImageName)
	docker.StopContainer(newImageName)
	_, _ = docker.RemoveContainer(newImageName)
}
