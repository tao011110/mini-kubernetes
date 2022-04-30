# kubelet

## kubelet启动的程序
 - etcd: 用于和master之间同步状态
 - cadvisor: 用于采集pod的信息和状态
 - kubelet-main

## kubelet-main的routines 
 - etcd_watcher: 监视etcd中key为`$node-name$_pos`的kv, 发现有pod被添加/暂停/删除时对节点中的对应pod做相应操作
 - resource_watcher: 定时向etcd中更新各container和节点物理资源的使用情况(CPU, memory)
 - pod_watcher: 定期通过http探针监控pod状态并在需要时重启
 - main: 启动etcd和cadvisor
