# Nacos 运维检查工具

旨在方便运维查看nacos注册服务，快速查找服务，同时生成prometheus自动发现所需要的json文件。   
golang 萌新，写的不好大佬勿喷... 😊

### 使用

```shell
  -cluster
        查看集群状态
  -find string
        查找服务
  -ipfile string
        ip解析文件 (default "salt_ip.json")
  -noconsole
        不输出console
  -second duration
        监控服务间隔刷新时间 (default 2s)
  -url string
        nacos地址 (default "http://dev-k8s-nacos:8848")
  -version
        查看版本
  -watch
        监控服务
  -write string
        prometheus 自动发现文件路径
```

因为默认只获取到主机ip，获取不到主机名,可以指定ipfile解析主机名，文件格式如下 (可选)

```shell
{
    "test1": "10.x.x.x",
    "test2": "10.x.x.x",
}
```
也可以使用salt批量获取主机名与ip的对应json关系
```shell
salt '*' network.interface_ip  eth0 --out=json --static -t 10  > /tmp/ip.json
```

定时任务示例

```shell
*/3 * * * * /data/script/nacos_check -url http://nacos-1:8848  -ipfile /data/script/ip.json -write nacos.json -noconsole
```

prometheus 可以结合blackbox_exporter使用

```yml
file_sd_configs:
  - files:
      - '/data/work/prometheus/discovery/*.json'
      refresh_interval: 3m
```
### 效果
![image](https://user-images.githubusercontent.com/20376675/154187473-96ced8e9-2c04-46aa-85b7-f3e44100e68d.png)
find 快速查找服务，支持以下👇匹配
- 匹配命名空间
- 匹配服务名
- 匹配IP端口
  ![image](https://user-images.githubusercontent.com/20376675/154187373-e180e679-0885-48cd-8b46-be3ad89fd53a.png)


### grafana 展示出图
![image](https://user-images.githubusercontent.com/20376675/154186534-35eed3db-70d8-461a-9aa6-df8cdcd7aa6c.png)
