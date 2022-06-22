参考文章：[(45条消息) （四）搭建容器云管理平台笔记—安装ETCD(不使用证书)_sas???的博客-CSDN博客](https://blog.csdn.net/weixin_33906657/article/details/92587119)

### etcd.service配置

vim /etc/systemd/system/etcd.service

```
[Unit]
Description=etcd.service

[Service]
Type=notify
TimeoutStartSec=0
ReStartSec=5s
Restart=always
WorkingDirectory=/var/lib/etcd
EnvironmentFile=-/opt/etcd/etcd.conf
ExecStart=/opt/etcd/etcd

[Install]
WantedBy=multi-user.target
```



### /opt/etcd/etcd.conf配置（Cloud1与Cloud2通过内网连接）

Cloud1：

vim /opt/etcd/etcd.conf

```
ETCD_NAME="etcd01"
ETCD_DATA_DIR="/var/lib/etcd/"

ETCD_LISTEN_PEER_URLS="http://192.168.1.7:2380"
ETCD_LISTEN_CLIENT_URLS="http://192.168.1.7:2379,http://127.0.0.1:2379"
ETCD_INITIAL_ADVERTISE_PEER_URLS="http://192.168.1.7:2380"
ETCD_ADVERTISE_CLIENT_URLS="http://192.168.1.7:2379"
ETCD_INITIAL_CLUSTER="etcd01=http://192.168.1.7:2380,etcd02=http://192.168.1.8:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
```



Cloud2：

vim /opt/etcd/etcd.conf

```
ETCD_NAME="etcd02"
ETCD_DATA_DIR="/var/lib/etcd/"

ETCD_LISTEN_PEER_URLS="http://192.168.1.8:2380"
ETCD_LISTEN_CLIENT_URLS="http://192.168.1.8:2379, http://127.0.0.1:2379"
ETCD_INITIAL_ADVERTISE_PEER_URLS="http://192.168.1.8:2380"
ETCD_ADVERTISE_CLIENT_URLS="http://192.168.1.8:2379"
ETCD_INITIAL_CLUSTER="etcd01=http://192.168.1.7:2380,etcd02=http://192.168.1.8:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
```



```
#重新加载etcd服务
systemctl daemon-reload
systemctl enable etcd.service
systemctl start etcd.service
```



```
#查看etcd集群member
etcdctl member list
#查看etcd就集群endpoint状态
etcdctl -w table --endpoint status
```

结点重启，若集群状态变化，要删除原先的历史数据。如该节点原先并未加入任何集群，更改etcd.conf文件后想加入某一集群，就要将原先的数据删除:`rm -rf /var/lib/etcd/*`

例如当遇到这种报错时：request sent was ignored (cluster ID mismatch: peer[3fc33c18ca0433f2]=3b1291d31874b132, local=9818358e4a92d2d9)   request cluster ID mismatch (got 9818358e4a92d2d9 want 3b29dfa314757683)

可以用etcdctl member list验证是不是仍然保留了与配置文件不同的历史数据。若确实如此，则可以予以删除



### 改用外网IP：

```
ETCD_NAME="etcd01"
ETCD_DATA_DIR="/var/lib/etcd/"

ETCD_LISTEN_PEER_URLS="http://192.168.1.7:2380"
ETCD_LISTEN_CLIENT_URLS="http://192.168.1.7:2379,http://127.0.0.1:2379"
ETCD_INITIAL_ADVERTISE_PEER_URLS="http://10.119.11.140:2380"
ETCD_ADVERTISE_CLIENT_URLS="http://10.119.11.140:2379, http://192.168.1.7:2379"
ETCD_INITIAL_CLUSTER="etcd01=http://10.119.11.140:2380, etcd02=http://10.119.11.111:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
```
