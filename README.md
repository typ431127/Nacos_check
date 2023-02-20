# Nacos 运维便携命令行检查工具

方便运维查看nacos注册服务，快速查找服务，同时生成prometheus自动发现所需要的json文件。   
golang 运维萌新，学习项目... 😊

### 安装
```shell
curl -L https://github.com/typ431127/Nacos_check/releases/download/0.6/nacos_check-linux-amd64 -o nacos_check
chmod +x nacos_check
./nacos_check --url https://nacos地址
```

### 使用帮助

```shell
Nacos

Usage:
  nacos_check [flags]
  nacos_check [command]

Available Commands:
  cluster     集群状态
  completion  Generate the autocompletion script for the specified shell
  config      查看本地配置文件路径
  help        Help about any command
  version     查看版本
  web         开启web api Prometheus http_sd_configs

Flags:
  -f, --find string            查找服务
  -h, --help                   help for nacos_check
  -i, --ipfile string          ip解析文件 (default "salt_ip.json")
      --json                   输出json
  -l, --lable stringToString   添加标签 -l env=dev,pro=java (default [])
  -s, --second duration        监控服务间隔刷新时间 (default 5s)
  -u, --url string             Nacos地址 (default "http://dev-k8s-nacos:8848")
  -w, --watch                  监控服务
  -o, --write string           导出json文件, prometheus 自动发现文件路径

Use "nacos_check [command] --help" for more information about a command.
```

#### 显示所有实例注册信息
```shell
./nacos_check-linux-amd64 --url http://nacos-0:8848 
```
![image](images/1.png)
#### 查看Nacos集群状态
```shell
./nacos_check-linux-amd64 --url http://nacos-0:8848 cluster --v2upgrade
```
![image](images/4.png)

#### 查找注册服务
```shell
./nacos_check-linux-amd64 --url http://nacos-0:8848 -f gateway 
./nacos_check-linux-amd64 --url http://nacos-0:8848 -f 8080
./nacos_check-linux-amd64 --url http://nacos-0:8848 -f 172.30
```
- 支持查找服务名，ip，端口,命名空间
#### 查找注册服务,每10秒刷新一次
```shell
./nacos_check-linux-amd64 --url http://nacos-0:8848 -f gateway  -w -s 10s
```


###  Prometheus自动发现支持

##### 写入自动发现json文件
```shell
./nacos_check-linux-amd64 --url http://nacos-0:8848 -o discovery.json
```

##### 控制台输出json
```shell
./nacos_check-linux-amd64 --url http://nacos-0:8848 --json
# 添加自定义label
./nacos_check-linux-amd64 --url http://nacos-0:8848  -l env=dev,pro=test-pro,k8s=true --json
```

#####  prometheus 可以结合blackbox_exporter使用

```yml
file_sd_configs:
  - files:
      - '/data/work/prometheus/discovery/*.json'
      refresh_interval: 3m
```

```shell
prometheus-file-sd 自动发现
./nacos_check-linux-amd64 --url http://nacos-0.xxxxx:8848 -o  discovery.json

http_sd_configs 自动发现
开启webapi        
./nacos_check-linux-amd64 web --url http://nacos-0.xxxx:8848

开启webapi并添加自定义label
./nacos_check-linux-amd64 web --url http://nacos-0.xxxx:8848 -l env=dev,pro=test-pro,k8s=true
```
**基于http_sd_configs的自动发现**
```yml
scrape_configs:
  - job_name: 'nacos'
    scrape_interval: 10s
    metrics_path: /probe
    params:
      module: [tcp_connect]
    http_sd_configs:
     - url: http://localhost:8099
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9115
```

#### find 快速查找服务，支持以下👇匹配
- 匹配命名空间
- 匹配服务名
- 匹配IP端口

```shell
# 模糊匹配命名空间
./nacos_check-linux-amd64 -f registry
# 模糊匹配服务
./nacos_check-linux-amd64 -f gateway
# 匹配端口
./nacos_check-linux-amd64 -f 8080
# 模糊匹配IP
./nacos_check-linux-amd64 -f 172.30
```
![image](images/3.png)

#### 加载本地配置
每次运行工具都需要指定url很麻烦，可以在本地写一个配置文件，这样默认情况下就会加载配置文件里面的url，就不需要每次都指定了。
查看配置文件路径
```shell
 ./nacos_check-linux-amd64 config
本地配置文件路径: /root/.nacos_conf.toml
```
`/root/.nacos_conf.toml` 示例
```toml
# nacos url地址
url = "http://nacos-0:8848"

# 定义容器网段
container_network = ["172.30.0.0/16","172.16.0.0/16","192.168.0.0/16"]

# 设置默认导出json和web服务附加标签
label = [
    {name = "env",value = "dev"},
    {name = "os",value = "linux"}
]
```
#### docker启动web服务 Prometheus httpd_sd_config 使用
```
docker run -itd -e nacos_url=http://nacos-xx.com:8848 -p 8099:8099 typ431127/nacos-check:0.6
访问 http://localhost:8099
```

#### 主机名解析
因为默认只获取到主机ip，获取不到主机名,可以指定ipfile解析主机名，有条件可以二次开发对接自己cmdb, 文件格式如下 (可选)

```shell
{
    "test1": "10.x.x.x",
    "test2": "10.x.x.x",
}
```
```shell
 ./nacos_check-linux-amd64 -i ../ip.json
```

### 效果
![image](images/1.png)

### grafana 展示出图

grafana控制台导入`grafana.json` 此模板默认匹配blackbox_exporter

![image](images/grafana.png)
