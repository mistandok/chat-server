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

func TestCreate_SuccessCreateChat(t *testing.T) {
	ctx := context.Background()
	var chatID int64 = 1
	var userIDs = []int64{1, 2, 3}
	request := &chat_v1.CreateRequest{UserIDs: userIDs}

	chatServiceMock := serviceMocks.NewChatService(t)
	chatServiceMock.On("Create", ctx, userIDs).Return(chatID, nil)

	chatImpl := impl.NewImplementation(chatServiceMock)
	resultChatID, err := chatImpl.Create(ctx, request)

	require.NoError(t, err)
	require.Equal(t, &chat_v1.CreateResponse{Id: chatID}, resultChatID)
}

func TestCreate_FailCreateChat(t *testing.T) {
	ctx := context.Background()
	var userIDs = []int64{1, 2, 3}
	request := &chat_v1.CreateRequest{UserIDs: userIDs}
	someErr := errors.New("some err")

	chatServiceMock := serviceMocks.NewChatService(t)
	chatServiceMock.On("Create", ctx, userIDs).Return(int64(0), someErr)

	chatImpl := impl.NewImplementation(chatServiceMock)
	_, err := chatImpl.Create(ctx, request)

	require.Error(t, err)
}
