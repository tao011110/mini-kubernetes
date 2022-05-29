./kubectl create -f /home/go/src/mini-kubernetes/presentation/gpu/mat_add/mat_add.yaml
./kubectl get gpujob
./kubectl describe gpujob mat_add_test
#登陆GPU平台 ssh stu614@sylogin.hpc.sjtu.edu.cn  9hF8#D7&
#然后去 /home/go/src/mini-kubernetes/tools/gpu_job_uploader_image/home/result/ 下查看
