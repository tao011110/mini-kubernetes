package gpu_job_uploader

import (
	"fmt"
	"github.com/jakehl/goid"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/util"
)

func generateImage(job *def.GPUJob) {
	newImageName := fmt.Sprintf("gpuJob-%s-%s", job.Name, goid.NewV4UUID().String())
	slurmContent := slurmGenerator(job.Slurm)
	apiServerIPFileContent := apiServerIPFileGenerator()
	sourceCodeContent := util.ReadFile(job.SourceCodePath)
	makefilePath := job.MakefilePath
	container := def.Container{
		Image: def.GPUJobUploaderImage,
	}
	containerID := docker.CreateContainer(container, newImageName)
	docker.CopyToContainer(containerID, def.GPUSlurmScriptParentDirPath, def.GPUSlurmScriptFileName, slurmContent)
	docker.CopyToContainer(containerID, def.GPUApiServerIpAndPortFileParentDirPath, def.GPUApiServerIpAndPortFileFileName, apiServerIPFileContent)
	docker.CopyToContainer(containerID, def.GPUJOBMakefileParentDirPath, def.GPUJOBMakefileFileName, makefilePath)
	docker.CopyToContainer(containerID, def.GPUJobSourceCodeParentDirPath, def.GPUJobSourceCodeFileName, sourceCodeContent)
	docker.CopyToContainer(containerID, def.GPUJobNameParentDirName, def.GPUJobNameFileName, job.Name)

	docker.CommitContainer(containerID, newImageName)
	docker.PushImage(newImageName)
	docker.StopContainer(containerID)
	_, _ = docker.RemoveContainer(containerID)
	job.ImageName = newImageName
}
