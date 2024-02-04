## 一、安装kubeadm

执行以下命令，让Master节点可以调度Pod

```shell
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
kubectl taint nodes --all node-role.kubernetes.io/master-
```

## 二、安装calico

安装官方文档安装即可：`https://docs.tigera.io/calico/latest/getting-started/kubernetes/quickstart`

安装后，修改`vxlanMode`为 `Always` 

```shell
root@VM-0-7-ubuntu:~# kubectl  edit IPPool default-ipv4-ippool
# vxlanMode: Always
```



## 三、配置Docker节点

### 1,Set up Docker's `apt` repository.

```shell
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
```



### 2，安装最新版Docker

```shell
apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
```

### 3，安装cri-dockerd

1)，下载cri-dockerd

```shell
# 到release页面https://github.com/Mirantis/cri-dockerd/releases，找到适合自己系统的软件包并下载。下面以Linux 64位 为例：
 wget https://github.com/Mirantis/cri-dockerd/releases/download/v0.3.7/cri-dockerd-0.3.7.amd64.tgz 
 tar zvxf cri-dockerd-0.3.7.amd64.tgz
 mkdir -p /usr/local/bin
 cd cri-dockerd/ && install -o root -g root -m 0755 cri-dockerd /usr/local/bin/cri-dockerd
```

2)，配置cri-dockerd systemd service

创建cri-docker.service文件

```shell
vim /etc/systemd/system/cri-docker.service
#输入一下内容，然后保存并退出
[Unit]
Description=CRI Interface for Docker Application Container Engine
Documentation=https://docs.mirantis.com
After=network-online.target firewalld.service docker.service
Wants=network-online.target
Requires=cri-docker.socket

[Service]
Type=notify
ExecStart=/usr/bin/cri-dockerd --container-runtime-endpoint fd://
ExecReload=/bin/kill -s HUP $MAINPID
TimeoutSec=0
RestartSec=2
Restart=always

# Note that StartLimit* options were moved from "Service" to "Unit" in systemd 229.
# Both the old, and new location are accepted by systemd 229 and up, so using the old location
# to make them work for either version of systemd.
StartLimitBurst=3

# Note that StartLimitInterval was renamed to StartLimitIntervalSec in systemd 230.
# Both the old, and new name are accepted by systemd 230 and up, so using the old name to make
# this option work for either version of systemd.
StartLimitInterval=60s

# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity

# Comment TasksMax if your systemd version does not support it.
# Only systemd 226 and above support this option.
TasksMax=infinity
Delegate=yes
KillMode=process

[Install]
WantedBy=multi-user.target
```

创建cri-docker.socket文件

```shell
vim /etc/systemd/system/cri-docker.socket
#输入以下内容，保存并退出
[Unit]
Description=CRI Docker Socket for the API
PartOf=cri-docker.service

[Socket]
ListenStream=%t/cri-dockerd.sock
SocketMode=0660
SocketUser=root
SocketGroup=docker

[Install]
WantedBy=sockets.target
```

启动cri-dockerd服务

```shell
systemctl daemon-reload
systemctl enable --now cri-docker.socket
```

### 4，将节点加入已有的集群

1），在主节点上执行如下命令，获取加入集群的Token以及相关命令

```shell
 kubeadm token create  --print-join-command
 # 输出如下
 # kubeadm join 10.203.0.7:6443 --token 5y9hs7.esm9t71slxujl0sv --discovery-token-ca-cert-hash sha256:9107cad9444756cf261122d33e2c6afa117864c321a790fd397a55497c46ce2f
```

2)，在新的节点上执行以下命令，使当前节点加入集群

```shell
kubeadm join 10.203.0.7:6443 --token 5y9hs7.esm9t71slxujl0sv --discovery-token-ca-cert-hash sha256:9107cad9444756cf261122d33e2c6afa117864c321a790fd397a55497c46ce2f --cri-socket unix:///var/run/cri-dockerd.sock
#注意：需要指定--cri-socket参数为cri-dockerd监听的socket文件
```

3)，观察新节点是否成功加入集群，现实新节点为ready状态即为成功

```shell
root@VM-0-7-ubuntu:~# kubectl  get node -o wide
NAME            STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE           KERNEL-VERSION      CONTAINER-RUNTIME
vm-0-6-ubuntu   Ready    <none>          12m   v1.28.4   10.203.0.6    <none>        Ubuntu 22.04 LTS   5.15.0-88-generic   docker://24.0.7
vm-0-7-ubuntu   Ready    control-plane   37d   v1.28.3   10.203.0.7    <none>        Ubuntu 22.04 LTS   5.15.0-86-generic   containerd://1.6.24
```

如图显示：

![image-20231203224014415](./hack/img/image-20231203224014415.png)

5，观察两个节点kubelet配置文件的差异，可以看到kubeadm已经将kubelet的`--container-runtime-endpoint`参数设置当前节点支持的container runtime。

```shell
#节点VM-0-7-ubuntu
root@VM-0-7-ubuntu:~# cat /var/lib/kubelet/kubeadm-flags.env
KUBELET_KUBEADM_ARGS="--container-runtime-endpoint=unix:///var/run/containerd/containerd.sock --pod-infra-container-image=registry.k8s.io/pause:3.9"

#节点VM-0-7-ubuntu
root@VM-0-6-ubuntu:~# cat /var/lib/kubelet/kubeadm-flags.env
KUBELET_KUBEADM_ARGS="--container-runtime-endpoint=unix:///var/run/cri-dockerd.sock --pod-infra-container-image=registry.k8s.io/pause:3.9"
```



## 四、安装Ingress

1，在集群的Master节点执行以下命令安装Ingress Controller到集群中。

```shell
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
```

2，（可选）测试Ingress Controller

**查看Pod**

```shell
root@VM-0-7-ubuntu:~# kubectl -n ingress-nginx get pod
NAME                                        READY   STATUS      RESTARTS   AGE
ingress-nginx-admission-create-p9sq7        0/1     Completed   0          12m
ingress-nginx-admission-patch-f4sdn         0/1     Completed   1          12m
ingress-nginx-controller-8558859656-j5sg2   1/1     Running     0          12m
```

**本地测试**

创建一个测试用web server以及相关的service

```shell
root@VM-0-7-ubuntu:~# kubectl create deployment demo --image=httpd --port=80
root@VM-0-7-ubuntu:~# kubectl get deployment
NAME   READY   UP-TO-DATE   AVAILABLE   AGE
demo   1/1     1            1           25s

root@VM-0-7-ubuntu:~# kubectl expose deployment demo
root@VM-0-7-ubuntu:~# kubectl  get svc
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
demo         ClusterIP   10.99.135.18   <none>        80/TCP    13s
```

创建Ingress

```shell
root@VM-0-7-ubuntu:~# kubectl create ingress demo-localhost --class=nginx --rule="demo.localdev.me/*=demo:80"
root@VM-0-7-ubuntu:~# kubectl get ingress
NAME             CLASS   HOSTS              ADDRESS   PORTS   AGE
demo-localhost   nginx   demo.localdev.me             80      4s
```

测试访问

```shell
root@VM-0-7-ubuntu:~# curl --resolve demo.localdev.me:30656:127.0.0.1 http://demo.localdev.me:30656
<html><body><h1>It works!</h1></body></html>
#注意上面的30656 为ingress-controller service映射到主机的端口，如下图PORTS(S)中的展示
root@VM-0-7-ubuntu:~# kubectl -n ingress-nginx get svc
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.100.54.47     <pending>     80:30656/TCP,443:30599/TCP   7h59m
```



附：安装containerd后的配置修改，`SystemdCgroup = true`，否则会造成容器无法启动，如etcd启动失败。

```yaml
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
  BinaryName = ""
  CriuImagePath = ""
  CriuPath = ""
  CriuWorkPath = ""
  IoGid = 0
  IoUid = 0
  NoNewKeyring = false
  NoPivotRoot = false
  Root = ""
  ShimCgroup = ""
  SystemdCgroup = true #......
```

