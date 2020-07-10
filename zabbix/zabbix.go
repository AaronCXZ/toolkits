package zabbix

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var ContextType = "application/json-rpc"

type Zabbix struct {
	Url      string
	user     string
	password string
	Token    string
}

type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
	Auth    string      `json:"auth"`
}

func New(url, user, password string) *Zabbix {
	return &Zabbix{
		Url:      url,
		user:     user,
		password: password,
	}
}

func DefaultRequest(auth string) *Request {
	return &Request{
		Jsonrpc: "2.0",
		Id:      1,
		Auth:    auth,
	}
}

func (z *Zabbix) httpPost(body *Request) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	resp, err := http.Post(z.Url, ContextType, bytes.NewReader(b))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return data
}
