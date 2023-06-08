package asocks

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
		var msg = make([]string, 0)
		for i := 0; i < 5; i++ {
			res, err := http.Get("https://api.asocks.com/api/v1/proxy-list/YgCw9JZbJyAlNSyyMeqa5CRf94IuU3vs.json")
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

			if len(msg) > 0 {
				break
			}
		}

		result := make([]iproxy.Proxy, 0)
		for _, m := range msg {
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
