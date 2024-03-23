package tests

import (
	"context"
	"errors"
	"testing"

	impl "github.com/mistandok/chat-server/internal/api/chat"
	serviceMocks "github.com/mistandok/chat-server/internal/service/mocks"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
)

func TestDelete_SuccessDeleteChat(t *testing.T) {
	ctx := context.Background()
	var chatID int64 = 1
	request := &chat_v1.DeleteRequest{Id: chatID}

	chatServiceMock := serviceMocks.NewChatService(t)
	chatServiceMock.On("Delete", ctx, chatID).Return(nil)

	chatImpl := impl.NewImplementation(chatServiceMock)
	_, err := chatImpl.Delete(ctx, request)

	require.NoError(t, err)
}

func TestDelete_FailDeleteChat(t *testing.T) {
	ctx := context.Background()
	var chatID int64 = 1
	request := &chat_v1.DeleteRequest{Id: chatID}
	deleteErr := errors.New("some err")

	chatServiceMock := serviceMocks.NewChatService(t)
	chatServiceMock.On("Delete", ctx, chatID).Return(deleteErr)

	chatImpl := impl.NewImplementation(chatServiceMock)
	_, err := chatImpl.Delete(ctx, request)

	require.Error(t, err)
}
