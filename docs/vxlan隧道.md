这是对集群中VXLAN隧道建立过程的详细操作说明：

```
       Cloud1                                                 Cloud2
 ----------------                                        -------------
 |     Cloud1    |                                        |   Cloud2   |
 ----------------                                        -------------
    |        |                                            |      |
    |        |----------------routers---------------------|      |
    br0   ens3（10.119.11.140)          (10.119.11.111) ens3    br0
    |                                                             |
接入内部网络miniK8S-bridge网桥(10.24.1.0/24)      接入内部网络miniK8S-bridge网桥 (10.24.2.0/24)                            
我们以在Cloud1上的配置为例
#安装openvswitch
sudo apt install openvswitch-switch

#### 1-添加br0网桥
sudo ovs-vsctl add-br br0

#### 2-让br0成为miniK8S-bridge网桥的一个端口（docker容器均接入到miniK8S-bridge网桥中）；
brctl addif miniK8S-bridge br0

#### 3-并将br0和miniK8S-bridge的虚拟网卡状态设置为up
ip link set dev br0 up
ip link set dev miniK8S-bridge up

#### 4-配置vxlan接口;
##注意：此处remote_ip为Cloud2公网ｉｐ
ovs-vsctl add-port br0 vx0 -- set interface vx0 type=vxlan options:remote_ip=10.119.11.111

#### 5-添加host2上私有网段的直连路由
ip route add 10.24.2.0/24 dev miniK8S-bridge 

############
#所有配置完成#
############

#查看网桥和端口
ovs-vsctl show
```

