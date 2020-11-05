package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Token struct {
	Errcode      int32  `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Access_token string `json:"access_token"`
	Expires_in   int32  `json:"expires_in"`
}

// 企业微信text类型消息
type WeComTextMsg struct {
	Touser                   string            `json:"touser"`
	Toparty                  string            `json:"toparty"`
	Totag                    string            `json:"totag"`
	Msgtype                  string            `json:"msgtype"`
	Agentid                  int32             `json:"agentid"`
	Text                     *WecomContentText `json:"text"`
	Safe                     int32             `json:"safe"`
	Enable_id_trans          int32             `json:"enable_id_trans"`
	Enable_duplicate_check   int32             `json:"enable_duplicate_check"`
	Duplicate_check_interval int32             `json:"duplicate_check_interval"`
}

// 企业微信text类型消息消息内容
type WecomContentText struct {
	Content string `json:"content"`
}

type WeComTextMsgResult struct {
	Errcode     int32  `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	Invaliduser string `json:"invaliduser"`
}

const (
	GET_TOKEN     string = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	POST_TEXT_MSG string = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
)

// 发送企业微信text消息
func PostText(corpid string, corpsecret string, topartys string, agentId int32, content string, client *http.Client) {
	token := GetToken(corpid, corpsecret, client)
	msg := getTextJsonMsg(content, topartys, agentId)
	jsonByte, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Marshal failed , err:%v \n", err)
		os.Exit(1)
	}
	fmt.Printf("post body:%s\n", string(jsonByte))
	buffer := bytes.NewBuffer(jsonByte)
	if token == "" {
		fmt.Printf("post text msg failed : token is empty")
		os.Exit(1)
	}
	url := POST_TEXT_MSG + "?access_token=" + token
	resp, err := client.Post(url, "application/json", buffer)
	if err != nil {
		fmt.Printf("post text message failed , err:%v \n", err)
		os.Exit(1)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("read body failed , err:%v \n", err)
		os.Exit(1)
	}
	if responseBody == nil || len(responseBody) == 0 {
		fmt.Println("response body is empty! ")
		os.Exit(1)
	}
	var result WeComTextMsgResult
	err = json.Unmarshal(responseBody, &result)
	if responseBody == nil || len(responseBody) == 0 {
		fmt.Printf("body json format failed,jsonStr:%s, err:%v \n", string(responseBody), err)
		os.Exit(1)
	}
	if result.Errmsg != "ok" {
		fmt.Printf("response not ok !,err:%v \n", err)
		os.Exit(1)
	}
}

func GetHttpClient(proxyAddr string) *http.Client {
	if proxyAddr == "" {
		return http.DefaultClient
	}
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil
	}
	netTransport := &http.Transport{
		//Proxy: http.ProxyFromEnvironment,
		Proxy: http.ProxyURL(proxy),
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,                             //每个host最大空闲连接
		ResponseHeaderTimeout: time.Second * time.Duration(5), //数据收发5秒超时
	}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

// 获取text类型消息的msg
func getTextJsonMsg(contentString string, topartys string, agentId int32) WeComTextMsg {
	var msg WeComTextMsg
	var content WecomContentText
	msg.Text = &content
	content.Content = contentString
	//topartyString := ""
	//if len(topartys) > 0 {
	//	for _, toparty := range topartys{
	//		topartyString = toparty + "|"
	//	}
	//}
	//if strings.HasSuffix(topartyString,"|") {
	//	topartyString = topartyString[:len(topartyString)-1]
	//}
	msg.Msgtype = "text"
	msg.Toparty = topartys
	msg.Agentid = agentId
	return msg
}

// 获取企业微信某应用的token
func GetToken(corpid string, corpsecret string, client *http.Client) string {
	url := GET_TOKEN + "?" +
		"corpid=" + corpid +
		"&corpsecret=" + corpsecret

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("get token failed , err:%v \n", err)
		os.Exit(1)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("get token read body failed , err:%v \n", err)
		os.Exit(1)
	}
	var token Token
	s := string(bytes)
	err = json.Unmarshal([]byte(s), &token)
	if err != nil {
		fmt.Printf("get token json format failed , err:%v \n", err)
		os.Exit(1)
	}
	if "ok" == token.Errmsg {
		return token.Access_token
	}
	return ""
}
