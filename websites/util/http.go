package util

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/jeek120/seostation/websites/iproxy"
)

func HttpGetHeader(u string, h map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("创建请求失败: %s", err)
		return nil, err
	}

	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := HttpDo(req)
	if err != nil {
		log.Printf("发送请求出错: %s", err)
		return nil, err
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取流失败: %s, url = %s", err, u)
		return nil, err
	}
	return bs, err
}

func HttpDo(req *http.Request) (res *http.Response, err error) {
	var client http.Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	for i := 0; i < 3; i++ {

		p := iproxy.Get()
		if p.Ip != "" {
			hUrl := fmt.Sprintf("http://%s:%d", p.Ip, p.Port)
			u, err1 := url.Parse(hUrl)
			if err1 != nil {
				log.Printf("生成代理ip失败: %s[%s]", err1, hUrl)
			}
			tr.Proxy = http.ProxyURL(u)
		}
		client.Transport = tr

		req.Header.Set("User-Agent", browser.Computer())

		res, err = client.Do(req)
		if err == nil {
			return
		}
		// time.Sleep(time.Duration(i+1) * time.Second)
	}
	return
}
