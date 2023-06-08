package util

import "strings"

var two_end_domain = []string{".ac.cn", ".ah.cn", ".bj.cn", ".com.cn", ".cq.cn", ".fj.cn", ".gd.cn", ".gov.cn", ".gs.cn", ".org.cn"}

func GetTopDomain(host string) string {
	var ss = strings.SplitAfterN(host, ".", 4)
	for _, h := range two_end_domain {
		if strings.HasSuffix(host, h) {
			return strings.Join(ss[len(ss)-3:], ".")
		}
	}
	return strings.Join(ss[len(ss)-2:], "")
}
