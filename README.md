# Nacos 运维检查工具

方便运维查看nacos注册服务，快速查找服务，同时生成prometheus自动发现所需要的json文件。   
golang 运维萌新，学习项目... 😊

### 使用

```shell
  -cluster
        查看集群状态
  -find string
        查找服务
  -ipfile string
        ip解析文件 (default "salt_ip.json")
  -json
        输出json
  -noconsole
        不输出console
  -second duration
        监控服务间隔刷新时间 (default 2s)
  -url string
        nacos地址 (default "http://dev-k8s-nacos:8848")
  -v2upgrade
        查看2.0升级状态,和-cluster一起使用
  -version
        查看版本
  -watch
        监控服务
  -write string
        prometheus 自动发现文件路径
```

#### 显示所有实例注册信息
![image](https://ddn-md.oss-cn-beijing.aliyuncs.com/images/md/2022/06/10/20220610104702.png)

#### 集群和升级状态

![image](https://ddn-md.oss-cn-beijing.aliyuncs.com/images/md/2022/06/10/20220610104930.png)

### 安装
```shell
curl -L https://github.com/typ431127/Nacos_check/releases/download/0.4.1/nacos_check-linux-amd64 -o nacos_check
chmod +x nacos_check
./nacos_check --url https://nacos地址
```

### 基本使用
##### 运维命令
```shell
./nacos_check --url https://nacos地址
```

#####  Prometheus自动发现

##### 写入自动发现json文件

```shell

nacos_check -write discover.json
```

##### 控制台输出json
```shell
nacos_check -json
```
##### 指定nacos url
```shell
nacos_check -url http://192.168.100.190:8848 -cluster
```
##### 查看nacos 集群和升级状态
```shell
nacos_check -url http://192.168.100.190:8848 -cluster -v2upgrade
```
#####  prometheus 可以结合blackbox_exporter使用

```yml
file_sd_configs:
  - files:
      - '/data/work/prometheus/discovery/*.json'
      refresh_interval: 3m
```

#### find 快速查找服务，支持以下👇匹配
- 匹配命名空间
- 匹配服务名
- 匹配IP端口

```shell
# 模糊匹配命名空间
nacos_check -find public
# 模糊匹配服务
nacos_check -find gateway-service
# 匹配端口
nacos_check -find 8080
# 模糊匹配IP
nacos_check -find 172.30.
```
#### 主机名解析
因为默认只获取到主机ip，获取不到主机名,可以指定ipfile解析主机名，有条件可以二次开发对接自己cmdb, 文件格式如下 (可选)

```shell
{
    "test1": "10.x.x.x",
    "test2": "10.x.x.x",
}
```

### 效果
![image](https://user-images.githubusercontent.com/20376675/154187473-96ced8e9-2c04-46aa-85b7-f3e44100e68d.png)

### grafana 展示出图

grafana控制台导入`grafana.json` 此模板默认匹配blackbox_exporter

![image](https://user-images.githubusercontent.com/20376675/154186534-35eed3db-70d8-461a-9aa6-df8cdcd7aa6c.png)
