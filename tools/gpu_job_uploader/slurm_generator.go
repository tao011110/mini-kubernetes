package gpu_job_uploader

import (
	"fmt"
	"mini-kubernetes/tools/def"
)

func slurmGenerator(config def.GPUSlurmConfig) string {
	//tempFilePath := def.TemplateFileDir + goid.NewV4UUID().String() + ".slurm"
	//fp, err := os.Create(tempFilePath)
	//if err != nil {
	//	return ""
	//}
	//defer func(fp *os.File) {
	//	_ = fp.Close()
	//}(fp)
	//_, _ = fp.WriteString("#!/bin/bash\n\n")
	//_, _ = fp.WriteString(fmt.Sprintf("#SBATCH --job-name=%s\n", config.JobName))
	//_, _ = fp.WriteString(fmt.Sprintf("#SBATCH --partition=%s\n", config.Partition))
	//_, _ = fp.WriteString(fmt.Sprintf("#SBATCH -N %d\n", config.Node))
	//_, _ = fp.WriteString(fmt.Sprintf("#SBATCH --cpus-per-task=%d\n", config.CpusPerTask))
	//_, _ = fp.WriteString(fmt.Sprintf("#SBATCH --ntasks-per-node=%d\n", config.NtasksPerNode))
	//_, _ = fp.WriteString(fmt.Sprintf("#SBATCH --gres=gpu:%d\n", config.GPU))
	//_, _ = fp.WriteString("#SBATCH --output=result.out\n")
	//_, _ = fp.WriteString("#SBATCH --error=error.err\n")
	//_, _ = fp.WriteString(fmt.Sprintf("#SBATCH --time=%s\n\n", config.Time))
	//_, _ = fp.WriteString("module load gcc cuda\n\n")
	//_, _ = fp.WriteString("make\n")
	//_, _ = fp.WriteString(fmt.Sprintf("./%s\n", config.TargetExecutableFileName))
	//return tempFilePath
	content := ``
	content += "#!/bin/bash\n\n"
	content += fmt.Sprintf("#SBATCH --job-name=%s\n", config.JobName)
	content += fmt.Sprintf("#SBATCH --partition=%s\n", config.Partition)
	content += fmt.Sprintf("#SBATCH -N %d\n", config.Node)
	content += fmt.Sprintf("#SBATCH --cpus-per-task=%d\n", config.CpusPerTask)
	content += fmt.Sprintf("#SBATCH --ntasks-per-node=%d\n", config.NtasksPerNode)
	content += fmt.Sprintf("#SBATCH --gres=gpu:%d\n", config.GPU)
	content += "#SBATCH --output=result.out\n"
	content += "#SBATCH --error=error.err\n"
	content += fmt.Sprintf("#SBATCH --time=%s\n\n", config.Time)
	content += "module load gcc cuda\n\n"
	content += "make\n"
	content += fmt.Sprintf("./%s\n", config.TargetExecutableFileName)
	return content
}
