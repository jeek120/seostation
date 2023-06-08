package stormproxies

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jeek120/seostation/websites/iproxy"
)

func init() {
	iproxy.Setup(func() []iproxy.Proxy {
		var msg IpMsg
		for i := 0; i < 5; i++ {
			u := "http://api.tq.roxlabs.cn/getProxyIp?num=32&return_type=json&lb=4&sb=&flow=1&regions=&protocol=http"
			log.Printf("开始获取代理ip: %s", u)
			res, err := http.Get(u)
			if err != nil {
				log.Printf("获取代理ip失败: %s", err)
			}
			bs, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("读取代理ip信息流失败: %s", err)
			}
			err = json.Unmarshal(bs, &msg)
			if err != nil {
				log.Printf("解析ip失败: %s", err)
			}

			if msg.Code == 0 {
				break
			}
		}

		result := make([]iproxy.Proxy, 0)
		for _, m := range msg.Data {
			result = append(result, iproxy.Proxy{
				Ip:   m.IP,
				Port: m.Port,
			})
		}
		return result
	})
}

type IpMsg struct {
	Code      int    `json:"code"`
	Success   bool   `json:"success"`
	Msg       string `json:"msg"`
	RequestIP string `json:"request_ip"`
	Data      []Data `json:"data"`
}
type Data struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
