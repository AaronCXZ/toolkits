package zabbix

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

type login struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (z *Zabbix) Login() {
	body := struct {
		Jsonrpc string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  login  `json:"params"`
		Id      int    `json:"id"`
	}{Jsonrpc: "2.0",
		Method: "user.login",
		Params: login{
			User:     z.user,
			Password: z.password,
		}}
	b, err := json.Marshal(body)
	if err != nil {
		return
	}
	resp, err := http.Post(z.Url, ContextType, bytes.NewReader(b))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	token := gjson.GetBytes(data, "result")
	z.Token = token.String()
}
