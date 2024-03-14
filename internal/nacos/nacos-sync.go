package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"nacos-check/pkg/fmtd"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Sync struct {
	client config_client.IConfigClient
}

func NewSync() *Sync {
	return &Sync{}
}

// runSyncTask 开启子同步任务
func (c *Sync) runSyncTask(config map[string]string, wg sync.WaitGroup) {
	defer wg.Done()
	var serverConfigs []constant.ServerConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         config["namespace"],
		TimeoutMs:           30000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "error",
	}
	for _, _url := range strings.Split(NACOSURL, ",") {
		parse, _ := url.Parse(_url)
		_host := strings.Split(parse.Host, ":")[0]
		parseUint, err := strconv.ParseUint(parse.Port(), 10, 64)
		if err != nil {
			fmtd.Fatalln(err)
		}
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr:      _host,
			ContextPath: FileConfig.ContextPath,
			Port:        parseUint,
			Scheme:      parse.Scheme,
		})
	}
	var err error
	c.client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	data, err := c.client.GetConfig(vo.ConfigParam{
		DataId: config["dataId"],
		Group:  config["group"],
	})
	if err != nil {
		log.Printf("Failed to get configuration %s", config["dest"])
	} else {
		c.syncWriteFile(config["dest"], []byte(data))
	}
	_ = c.client.ListenConfig(vo.ConfigParam{
		DataId: config["dataId"],
		Group:  config["group"],
		OnChange: func(namespace, group, dataId, data string) {
			log.Printf("Configuration file %s has changes\n", config["dataId"])
			c.syncWriteFile(config["dest"], []byte(data))
		},
	})

}

// syncWriteFile 配置写入
func (c *Sync) syncWriteFile(path string, data []byte) {
	err := c.writeFile(path, data)
	if err != nil {
		log.Printf("Failed to write configuration file %s %v\n", path, err)
	} else {
		log.Printf("The configuration file %s was written successfully", path)
	}
}

// parameterVerification 参数校验
func (c *Sync) parameterVerification(config map[string]string) error {
	for _, key := range []string{"namespace", "dataId", "group", "dest"} {
		if _, ok := config[key]; !ok {
			return fmt.Errorf("%s does not exist", key)
		}
	}
	return nil
}

// createDirectory 创建目录
func (c *Sync) createDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

// writeFile 写入文件
func (c *Sync) writeFile(path string, data []byte) error {
	err := c.createDirectory(filepath.Dir(path))
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// RunSync 开始运行同步
func (c *Sync) RunSync() {
	var wg sync.WaitGroup
	if len(FileConfig.Sync) != 0 {
		for _, config := range FileConfig.Sync {
			err := c.parameterVerification(config)
			if err != nil {
				log.Println(config, "Configuration error, skip")
				continue
			}
			wg.Add(1)
			go func(config map[string]string) {
				c.runSyncTask(config, wg)
			}(config)
		}
		wg.Wait()
	} else {
		fmtd.Fatalln("no configuration information exit")
	}

}
