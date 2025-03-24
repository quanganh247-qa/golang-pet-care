package chatbot

import (
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/handlers"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup, chatHandler *handlers.ChatHandler) {
	chatbot := routerGroup.RouterDefault.Group("/")

	// authRoute := routerGroup.RouterAuth(chatbot)
	{
		chatbot.POST("chat", chatHandler.HandleChatRequest)
	}

}
