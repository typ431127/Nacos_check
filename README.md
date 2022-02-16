# Nacos 检查工具

方便运维查看nacos注册服务，同时生成prometheus自动发现所需要的json文件。

### 使用

```shell
  -find string
        查找服务
  -ipfile string
        ip解析文件 (default "salt_ip.json")
  -noconsole
        输出console
  -url string
        nacos地址 (default "http://nacos-0.nacos-headless.mall.svc.cluster.local:8848")
  -write string
        prometheus 自动发现文件路径 (default "/data/work/prometheus/discovery/nacos.json")
```

因为默认只获取到主机ip，获取不到主机名,可以指定ipfile解析主机名，文件格式如下

```shell
{
    "test1": "10.x.x.x",
    "test2": "10.x.x.x",
}
```

定时任务示例

```shell
*/3 * * * * /data/script/nacos_check -url http://yestae-zk-1:8848  -ipfile /data/script/ip.json -noconsole
```

prometheus 可以结合http探针使用

```yml
file_sd_configs:
  - files:
      - '/data/work/prometheus/check_nginx/*.json'
      refresh_interval: 5m
```