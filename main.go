package main

import (
	"fmt"
	"os"
	"post-wecom/post"
	"strconv"
)

func main() {
	var subject,
		message,
		corpid,
		corpsecret,
		topartysString,
		agentIdStr,
		proxyAddr string
	for i, arg := range os.Args {
		switch i {
		case 1:
			subject = arg
		case 2:
			message = arg
		case 3:
			corpid = arg
		case 4:
			corpsecret = arg
		case 5:
			topartysString = arg
		case 6:
			agentIdStr = arg
		case 7:
			proxyAddr = arg
		}
	}
	agentId, _ := strconv.ParseInt(agentIdStr, 10, 32)
	agentId32 := int32(agentId)
	content := subject + "\n" + message
	client := post.GetHttpClient(proxyAddr)
	post.PostText(corpid,
		corpsecret,
		topartysString,
		agentId32,
		content,
		client)
	fmt.Println("脚本执行完毕")

}
