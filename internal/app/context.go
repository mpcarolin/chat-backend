package app

import "chat-backend/internal/chat"

type AppContext struct {
	ChatProvider chat.ChatProvider
}

func NewAppContext(chatProvider chat.ChatProvider) *AppContext {
	return &AppContext{
		ChatProvider: chatProvider,
	}
}