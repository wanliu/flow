package context

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var WechatTypeEnum = map[string]int{
	"Private": 1,
	"Group":   2,
	"GroupAt": 3,
}

func PrivateChat(ctx Context) bool {
	t := wechatType(ctx)

	return t == 1 || t == 3
}

func GroupChat(ctx Context) bool {
	t := wechatType(ctx)

	return t == 2
}

// {"BotName":"客服",
// "Content":"这个新版的微信, 容易挂掉",
// "MessageId":"859f8ba8-4ff7-44ae-8d13-1ed3deedd6a9",
// "MsgID":"5571466961726285032",
// "MsgType":"1",
// "ReceiptHandle":"xxx",
// "RoomID":"@Room\u003c测试群\u003e",
// "RoomName":"测试群",
// "SendId":"@f048968daa9e01bb4fbf98067a21ec32",
// "SenderNickName":"胡义",
// "TaskId":"528219641",
// "WeixinUin":"528219641"}
func wechatType(ctx Context) int {
	// log.Printf("[WECHAT GROUP] BEGIN")
	data, _ := json.Marshal(ctx.Value("WECHAT_INFO"))
	log.Printf("[WECHAT GROUP] ctx value: %v", string(data))
	// log.Printf("[WECHAT GROUP] END")

	wechatInfo := ctx.Value("WECHAT_INFO")

	if wechatInfo == nil {
		return WechatTypeEnum["Private"]
	} else {
		message := wechatInfo.(map[string]string)

		rooId, ok := message["RoomID"]
		if ok && rooId != "" {
			if bootName, ok := message["BotName"]; ok {
				content, _ := message["Content"]
				if strings.Contains(content, fmt.Sprintf("@%v", bootName)) {
					return WechatTypeEnum["GroupAt"]
				}
			}

			return WechatTypeEnum["Group"]
		} else {
			return WechatTypeEnum["Private"]
		}
	}

	return 0
}
