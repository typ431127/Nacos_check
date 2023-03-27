package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (d *Nacos) GET(apiurl string) []byte {
	u, err := url.Parse(apiurl)
	if err != nil {
		panic(err)
	}
	if len(USERNAME) != 0 && len(PASSWORD) != 0 {
		if len(u.RawQuery) == 0 {
			apiurl += "?accessToken=" + url.QueryEscape(d.Token)
		} else {
			apiurl += "&accessToken=" + url.QueryEscape(d.Token)
		}
	}
	req, _ := http.NewRequest("GET", apiurl, nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 501 && u.Path == "/nacos/v1/ns/operator/servers" {
			//panic(fmt.Sprintf("单机不支持集群,请求状态码异常:%d", res.StatusCode))
			_url := fmt.Sprintf("%s/nacos/v2/core/cluster/node/list", d.DefaultUlr)
			return d.GET(_url)
		}
		if res.StatusCode == 501 && u.Path == "/nacos/v1/ns/upgrade/ops/metrics" {
			panic(fmt.Sprintf("此版本不支持查看升级状态:%d", res.StatusCode))
		}
		if res.StatusCode == 403 {
			panic(fmt.Sprintf("%s请求状态码异常:%d 请使用--username --password参数进行鉴权", apiurl, res.StatusCode))
		}
		panic(fmt.Sprintf("%s请求状态码异常:%d", apiurl, res.StatusCode))
	}
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	return resp

}

func (d *Nacos) POST(apiurl string, formData map[string]string) []byte {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for key, val := range formData {
		_ = bodyWriter.WriteField(key, val)
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	var req *http.Request
	u, err := url.Parse(apiurl)
	if u.Path == "/nacos/v1/auth/login" {
		req, _ = http.NewRequest("POST", apiurl, bodyBuf)
		req.Header.Set("Content-Type", contentType)
	}
	res, err := d.Client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		if u.Path == "/nacos/v1/auth/login" && res.StatusCode == 403 {
			panic(fmt.Sprintf("%s请求状态码异常,认证失败!:%d", apiurl, res.StatusCode))
		}
		panic(fmt.Sprintf("%s请求状态码异常:%d", apiurl, res.StatusCode))
	}
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	return resp
}
