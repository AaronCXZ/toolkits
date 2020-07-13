package zabbix

import (
	"fmt"

	"github.com/Muskchen/toolkits/errors"

	"github.com/tidwall/gjson"
)

type trigger struct {
	Host            string `json:"host"`
	Output          string `json:"output"`
	SelectFunctions string `json:"selectfunctions"`
}
type update struct {
	Triggerid string `json:"triggerid"`
	Status    string `json:"status"`
}

func (z *Zabbix) GetTriggers(host string) (triggers map[string]string) {
	var t = make(map[string]string, 0)
	body := DefaultRequest(z.Token)
	body.Method = "trigger.get"
	body.Params = trigger{
		Host:            host,
		Output:          "extend",
		SelectFunctions: "extend",
	}
	data := z.httpPost(body)
	results := gjson.GetBytes(data, "result").Array()
	for _, result := range results {
		//id := gjson.Get(result.String(), "description")
		many := gjson.GetMany(result.String(), "triggerid", "description")
		if len(many) != 2 {
			continue
		}
		t[many[0].String()] = many[1].String()
	}
	return t
}

func (z *Zabbix) StatTrigger(tid, status string) error {
	if status == "stop" {
		status = "1"
	} else if status == "start" {
		status = "0"
	} else {
		return errors.New(fmt.Sprintf("status err: %s", status))
	}

	body := DefaultRequest(z.Token)
	body.Method = "trigger.update"
	body.Params = update{
		Triggerid: tid,
		Status:    status,
	}
	data := z.httpPost(body)
	result := gjson.GetBytes(data, "result")
	fmt.Println(result.String())
	return nil
}

func GetTid(triggers map[string]string, name string) string {
	for k, v := range triggers {
		if v == name {
			return k
		}
	}
	return ""
}
