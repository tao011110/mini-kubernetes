# 关于DNS和Gateway的一些注意事项

#### 以下操作均在master上执行！

1. 每次作为master的云主机重启时，要执行：

   ```
   #先停用 systemd-resolved 服务，因为云主机默认开启systemd-resolved监听53端口，导致coredns因端口占用而无法启动
   sudo systemctl stop systemd-resolved
   #修改/etc/resolv.conf
   vim /etc/resolv.conf
   #在第一行写入：
   nameserver 127.0.0.1
   ```

   然后才能正常启动coredns

2. 在master的云主机上创建/home/docker_test.sh文件，并赋予其**可读可写可执行**权限

3. 在/etc/hosts中第一行的127.0.0.1 localhost 后写上本台云主机的名称。

   eg：我的一台云主机叫做cloud1，所以这一行就改为  127.0.0.1 localhost cloud1

   这是为了避免出现：`sudo: unable to resolve host cloud1: Temporary failure in name resolution`类似的报错