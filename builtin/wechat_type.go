package builtin

import (
	"fmt"
	"strings"

	"github.com/wanliu/flow/context"
)

var WechatTypeEnum = map[string]int{
	"Private": 1,
	"Group":   2,
	"GroupAt": 3,
}

func PrivateChat(ctx context.Context) bool {
	t := wechatType(ctx)

	return t == 1 || t == 3
}

func GroupChat(ctx context.Context) bool {
	t := wechatType(ctx)

	return t == 2
}

// SendId:"@296b2ba89984ea1f66dd5f6753245376",
// SenderNickName:"胡义",
// RoomName:"测试群",
// RoomID:"@Room<测试群>",
// TaskId:"528219641",
// MsgID:"8833967666994845528",
// MsgType:"1",
// MessageId:"691086ae-b9ba-47b1-8971-e1c2fc4639cd",
// ReceiptHandle:"xxxxx",
// Content:"@客服 在吗",
// WeixinUin:"528219641",
// BotName:"客服"
func wechatType(ctx context.Context) int {
	wechatInfo := ctx.Value("WECHAT_INFO")

	if wechatInfo == nil {
		return WechatTypeEnum["Private"]
	} else {
		message := wechatInfo.(map[string]string)

		rooId, ok := message["RoomId"]
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
