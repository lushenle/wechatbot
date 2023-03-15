package bootstrap

import (
	"log"

	"github.com/eatmoreapple/openwechat"
	"github.com/lushenle/wechatbot/handlers"
)

func Run() {
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 注册消息处理函数
	handler, err := handlers.NewHandler()
	if err != nil {
		log.Fatalf("register error: %v", err)
	}
	bot.MessageHandler = handler

	// 注册登陆二维码回调
	bot.UUIDCallback = handlers.QrCodeCallBack

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	// 执行热登录
	err = bot.HotLogin(reloadStorage, true)
	if err != nil {
		log.Fatalf("login error: %v ", err)
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
