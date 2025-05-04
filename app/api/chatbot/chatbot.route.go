package chatbot

import (
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/handlers"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup, chatHandler *handlers.ChatHandler) {
	chatbot := routerGroup.RouterDefault.Group("/")

	// Public routes
	{
		chatbot.POST("chat", chatHandler.HandleChatRequest)
		chatbot.OPTIONS("chat", chatHandler.HandleOptionsRequest)
	}

	// Authentication required routes
	authRoute := routerGroup.RouterAuth(chatbot)
	{
		authRoute.GET("conversations", chatHandler.ListConversations)
		authRoute.GET("conversations/:conversation_id", chatHandler.GetConversation)
		authRoute.DELETE("conversations/:conversation_id", chatHandler.DeleteConversation)

	}
}
