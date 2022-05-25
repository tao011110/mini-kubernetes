# to gpu
cd /home/gpu
# read job name
jobName=`cat /home/job_name`
# upload file to server
targetDirName="/dssg/home/acct-stu/stu614/gpu-${jobName}"
echo $targetDirName
scp -r /home/gpu stu614@sylogin.hpc.sjtu.edu.cn:$targetDirName
# submit job
ret=`ssh stu614@sylogin.hpc.sjtu.edu.cn "cd $targetDirName ; sbatch job.slurm "`
id=`echo ${ret:20}`
echo $id
out=result.out
files=`ssh stu614@sylogin.hpc.sjtu.edu.cn "cd $targetDirName ; ls"`
echo $files
while [[ $files != *$out* ]]; do
    sleep 5
    files=`ssh stu614@sylogin.hpc.sjtu.edu.cn "cd $targetDirName ; ls"`
    echo $files
done
echo finish
# download result from server
sleep 5
scp stu614@sylogin.hpc.sjtu.edu.cn:$targetDirName/result.out /home/result
scp stu614@sylogin.hpc.sjtu.edu.cn:$targetDirName/error.err /home/result
# remove dir in server
ssh stu614@sylogin.hpc.sjtu.edu.cn "cd /dssg/home/acct-stu/stu614 ; rm -rf gpu-${jobName}"
# http send result to server
cd /home/result
python3 sendResultToApiserver.py
