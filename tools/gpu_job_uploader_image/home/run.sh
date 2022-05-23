cd /home/gpu
# # compile
export PATH=/usr/local/cuda-11.7/bin:$PATH
export LD_LIBRARY_PATH=/usr/local/cuda-11.7/lib64:$LD_LIBRARY_PATH
# upload file to server
scp -r /home/gpu stu614@sylogin.hpc.sjtu.edu.cn:/dssg/home/acct-stu/stu614/gpu
# submit job
ret=`ssh stu614@sylogin.hpc.sjtu.edu.cn "cd /dssg/home/acct-stu/stu614/gpu ; sbatch job.slurm "`
id=`echo ${ret:20}`
echo $id
out=result.out
files=`ssh stu614@sylogin.hpc.sjtu.edu.cn "cd /dssg/home/acct-stu/stu614/gpu ; ls"`
echo $files
while [[ $files != *$out* ]]; do
    sleep 5
    files=`ssh stu614@sylogin.hpc.sjtu.edu.cn "cd /dssg/home/acct-stu/stu614/gpu ; ls"`
    echo $files
done
echo finish
# download result from server
scp stu614@sylogin.hpc.sjtu.edu.cn:/dssg/home/acct-stu/stu614/gpu/result.out /home/result
scp stu614@sylogin.hpc.sjtu.edu.cn:/dssg/home/acct-stu/stu614/gpu/error.err /home/result
# remove dir in server
ssh stu614@sylogin.hpc.sjtu.edu.cn "cd /dssg/home/acct-stu/stu614 ; rm -rf gpu"
# http send result to server
cd /home/result
python3 sendResultToApiserver.py
