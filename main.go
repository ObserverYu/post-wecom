package main

import (
	"fmt"
	"os"
	"post-wecom/post"
	"strconv"
)

func main() {
	subject := os.Args[1]
	message := os.Args[2]
	corpid := os.Args[3]
	corpsecret := os.Args[4]
	topartysString := os.Args[5]
	agentIdStr := os.Args[6]
	agentId, _ := strconv.ParseInt(agentIdStr, 10, 32)
	agentId32 := int32(agentId)
	content := subject + "\n" + message
	post.PostText(corpid,
		corpsecret,
		topartysString,
		agentId32,
		content)
	fmt.Println("脚本执行完毕")

}
