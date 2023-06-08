package qingtingip

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
			u := "http://api.qingtingip.com/ip?app_key=9c16792dcb0e736781dd78bb2f906027&num=5&ptc=http&fmt=json&lb=\\n&city=2&port=0&mr=1&area_id=0|"
			log.Printf("开始获取代理ip: %s", u)
			res, err := http.Get(u)
			if err != nil {
				log.Printf("获取代理ip失败: %s", err)
				return nil
			}
			bs, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("读取代理ip信息流失败: %s", err)
			}
			err = json.Unmarshal(bs, &msg)
			if err != nil {
				log.Printf("解析ip失败: %s", err)
			}

			if msg.Code == 200 {
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
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []Data `json:"data"`
}

type Data struct {
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	ExpireTime string `json:"expire_time"`
	City       string `json:"city"`
	Isp        string `json:"isp"`
}
