# kubernetes-workshop

这是一个Golang工程，包含三个service组件。用在kuberetes workshop中。

### 1，目录介绍

```bash
├── cmd # service的main文件位置
│   ├── serviceA # serviceA
│   │   └── main.go # serviceA 的 main文件，函数入口
│   ├── serviceB
│   │   └── main.go
│   └── serviceC
│       └── main.go
├── Dockerfile
├── go.mod
├── go.sum
├── hack # 脚本相关的存放目录
│   ├── deployment # 测试使用的yaml文件
│   │   ├── ingress.yaml # ingress ,将serviceA通过ingress controller暴露出去
│   │   ├── namespace.yaml
│   │   ├── servicea.yaml # serviceA的deployment和service
│   │   ├── serviceb.yaml
│   │   └── servicec.yaml
│   └── make-rules
│       └── build.sh # 编译二进制的脚本
├── init-kubernetes.md # kubernetes集群setup相关的说明
├── LICENSE
├── Makefile
├── pkg # 服务的核心代码目录
│   └── handlers
│       └── kubernetes_workshop_handlers.go
├── README.md
```

### 2，编译二进制及Docker Image

#### 本地运行

可以通过`make run WHAT=XXX`命令在本地运行某个服务。如下命令，表示运行serviceA服务，服务运行后默认监听`8080`端口

```bash
root@VM-0-7-ubuntu:~/kubernetes-workshop# make run WHAT=cmd/serviceA
```

访问服务

```shell
liuzhenweideMacBook-Pro kubernetes-workshop % curl http://127.0.0.1:8080/info
{"ServiceName":"serviceA","Version":"","Hostname":"liuzhenweideMacBook-Pro.local"}
```

#### 编译二进制

`注：以下示例命令在一台X86_64的Linux服务器上执行`

可以通过`make all WHAT=XXX`编译指定的service，`WHAT`为服务路径。如下命令，表示编译serviceA服务，编译后的二进制文件在`./bin/cmd/`目录中。

```shell
root@VM-0-7-ubuntu:~/kubernetes-workshop# make all WHAT=cmd/serviceA
root@VM-0-7-ubuntu:~/kubernetes-workshop# ls bin/cmd/ #可以看到编译出的二进制文件,serviceA
```

#### Build Docker Image

Build docker Image和上面的编译二进制并没有直接关系，只是两种不同的出包形态。在build docker image的之后会先执行编译二进制的操作，然后再build docker image。

可以通过`make image WHAT=XXX FULL_IMAGE_NAME={FULL_IMAGE_NAME}`build指定服务的docker镜像。`WHAT`为服务路径，`FULL_IMAGE_NAME`为完整的镜像名，可以指定image repo和tag。如下命令，表示build serviceA的docker镜像，且镜像完整名称为`duizhang/k8s-workshop-servicea`，即image tag为`latest`。

```shell
root@VM-0-7-ubuntu:~/kubernetes-workshop# make image WHAT=cmd/serviceA FULL_IMAGE_NAME=duizhang/k8s-workshop-servicea
```



