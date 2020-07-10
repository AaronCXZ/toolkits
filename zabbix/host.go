package zabbix

import (
	"github.com/tidwall/gjson"
)

type host struct {
	Output []string `json:"output"`
}

func (z *Zabbix) GetHosts() (hosts []string) {
	body := DefaultRequest(z.Token)
	body.Method = "host.get"
	body.Params = host{
		Output: []string{"hostid", "host"},
	}
	data := z.httpPost(body)
	results := gjson.GetBytes(data, "result")
	for _, result := range results.Array() {
		hosts = append(hosts, gjson.Get(result.String(), "host").String())
	}
	return hosts
}
