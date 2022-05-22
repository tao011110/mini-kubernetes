# 关于DNS和Gateway的一些注意事项

1. 每次云主机重启时，要执行：

   ```
   #先停用 systemd-resolved 服务，因为云主机默认开启systemd-resolved监听53端口，导致coredns因端口占用而无法启动
   sudo systemctl stop systemd-resolved
   #修改/etc/resolv.conf
   vim /etc/resolv.conf
   #在第一行写入：
   nameserver 127.0.0.1
   ```

   然后才能正常启动coredns

2. 在每台云主机上创建/home/docker_test.sh文件，并赋予其**可读可写可执行**权限