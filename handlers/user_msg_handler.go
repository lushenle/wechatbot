package handlers

import (
	"fmt"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/lushenle/wechatbot/config"
	"github.com/lushenle/wechatbot/gpt"
	"github.com/lushenle/wechatbot/pkg/logger"
	"github.com/lushenle/wechatbot/service"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
	// 接收到消息
	msg *openwechat.Message
	// 发送的用户
	sender *openwechat.User
	// 实现的用户业务
	service service.UserServiceInterface
}

func UserMessageContextHandler() func(ctx *openwechat.MessageContext) {
	return func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		handler, err := NewUserMessageHandler(msg)
		if err != nil {
			logger.Warning(fmt.Sprintf("init user message handler error: %s", err))
		}

		// 处理用户消息
		err = handler.handle()
		if err != nil {
			logger.Warning(fmt.Sprintf("handle user message error: %s", err))
		}
	}
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler(message *openwechat.Message) (MessageHandlerInterface, error) {
	sender, err := message.Sender()
	if err != nil {
		return nil, err
	}
	userService := service.NewUserService(c, sender)
	handler := &UserMessageHandler{
		msg:     message,
		sender:  sender,
		service: userService,
	}

	return handler, nil
}

// handle 处理消息
func (h *UserMessageHandler) handle() error {
	if h.msg.IsText() {
		return h.ReplyText()
	}
	return nil
}

// ReplyText 发送文本消息到群
func (h *UserMessageHandler) ReplyText() error {
	logger.Info(fmt.Sprintf("Received User %v Text Msg : %v", h.sender.NickName, h.msg.Content))
	var (
		reply string
		err   error
	)
	// 1.获取上下文，如果字符串为空不处理
	requestText := h.getRequestText()
	if requestText == "" {
		logger.Info("user message is null")
		return nil
	}

	logger.Info(fmt.Sprintf("h.sender.NickName == %+v", h.sender.NickName))
	// 2.向GPT发起请求，如果回复文本等于空,不回复
	reply, err = gpt.Completions(h.getRequestText())
	if err != nil {
		// 2.1 将GPT请求失败信息输出给用户，省得整天来问又不知道日志在哪里。
		errMsg := fmt.Sprintf("gpt request error: %v", err)
		_, err = h.msg.ReplyText(errMsg)
		if err != nil {
			return fmt.Errorf("response user error: %v ", err)
		}
		return err
	}

	// 2.设置上下文，回复用户
	h.service.SetUserSessionContext(requestText, reply)
	_, err = h.msg.ReplyText(buildUserReply(reply))
	if err != nil {
		return fmt.Errorf("response user error: %v ", err)
	}

	// 3.返回错误
	return err
}

// getRequestText 获取请求接口的文本，要做一些清晰
func (h *UserMessageHandler) getRequestText() string {
	var requestText string
	// 私聊没有前缀不处理
	// 处理前缀大小写，将消息全部转换为小写
	// 去除空格和换行
	if !strings.Contains(strings.ToLower(h.msg.Content), config.LoadConfig().PrivateTrigger) {
		return ""
	} else {
		requestText = strings.ToLower(h.msg.Content)
		requestText = strings.TrimLeft(requestText, config.LoadConfig().PrivateTrigger)
		requestText = strings.Trim(strings.TrimSpace(requestText), "\n")
	}

	// 获取上下文，拼接在一起，如果字符长度超出4000，截取为4000。（GPT按字符长度算），达芬奇3最大为4068，也许后续为了适应要动态进行判断。
	sessionText := h.service.GetUserSessionContext()
	if sessionText != "" {
		requestText = sessionText + "\n" + requestText
	}
	if len(requestText) >= 4000 {
		requestText = requestText[:4000]
	}

	// 返回请求文本
	return requestText
}

// buildUserReply 构建用户回复
func buildUserReply(reply string) string {
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return "请求得不到任何有意义的回复，请具体提出问题。"
	}

	// 2.如果用户有配置前缀，加上前缀
	reply = config.LoadConfig().ReplyPrefix + "\n" + reply
	reply = strings.Trim(reply, "\n")

	// 3.返回拼接好的字符串
	return reply
}
