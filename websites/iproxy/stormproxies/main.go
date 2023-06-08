package stormproxies

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jeek120/seostation/websites/iproxy"
)

func init() {
	iproxy.Setup(func() []iproxy.Proxy {
		var msg IpMsg
		for i := 0; i < 5; i++ {
			res, err := http.Get("http://api.qingtingip.com/ip?app_key=9c16792dcb0e736781dd78bb2f906027&num=5&ptc=http&fmt=json&lb=\n&city=2&port=0&mr=1&area_id=0|")
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

			if msg.Code == 200 {
				break
			}
		}

		result := make([]iproxy.Proxy, 0)
		for _, m := range msg.Data.List {
			ss := strings.Split(m, ":")
			port, err := strconv.ParseInt(ss[1], 10, 32)
			if err != nil {
				log.Printf("转换代理ip的端口失败: %s", err)
			}
			result = append(result, iproxy.Proxy{
				Ip:   ss[0],
				Port: int(port),
			})
		}
		return result
	})
}

type IpMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}
type Data struct {
	List []string `json:"list"`
}
