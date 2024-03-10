package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	impl "github.com/mistandok/chat-server/internal/api/chat"
	"github.com/mistandok/chat-server/internal/model"
	"github.com/mistandok/chat-server/internal/repository"
	serviceMocks "github.com/mistandok/chat-server/internal/service/mocks"
	"github.com/mistandok/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSendMessage_SuccessSendMessage(t *testing.T) {
	ctx := context.TODO()
	var userID int64 = 1
	var chatID int64 = 1
	curTime := time.Now().UTC()

	request := sendMessageRequestWithSetup(userID, chatID, curTime)
	message := messageWithSetup(userID, chatID, curTime)

	chatServiceMock := serviceMocks.NewChatService(t)
	chatServiceMock.On("SendMessage", ctx, *message).Return(nil)

	chatImpl := impl.NewImplementation(chatServiceMock)
	_, err := chatImpl.SendMessage(ctx, request)

	require.NoError(t, err)
}

func TestSendMessage_FailSendMessage(t *testing.T) {
	ctx := context.TODO()
	var userID int64 = 1
	var chatID int64 = 1
	curTime := time.Now().UTC()

	request := sendMessageRequestWithSetup(userID, chatID, curTime)
	message := messageWithSetup(userID, chatID, curTime)

	internalBadError := errors.New("internal error")

	tests := []struct {
		name              string
		sendRequest       *chat_v1.SendMessageRequest
		message           model.Message
		internalError     error
		expectedErrorCode codes.Code
	}{
		{
			name:              "fail send message because chat not found",
			sendRequest:       request,
			message:           *message,
			internalError:     repository.ErrChatNotFound,
			expectedErrorCode: codes.NotFound,
		},
		{
			name:              "fail send message because user not found",
			sendRequest:       request,
			message:           *message,
			internalError:     repository.ErrUserNotFound,
			expectedErrorCode: codes.NotFound,
		},
		{
			name:              "fail send message because user not in the chat",
			sendRequest:       request,
			message:           *message,
			internalError:     internalBadError,
			expectedErrorCode: codes.Internal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chatServiceMock := serviceMocks.NewChatService(t)
			chatServiceMock.On("SendMessage", ctx, *message).Return(test.internalError)

			chatImpl := impl.NewImplementation(chatServiceMock)
			_, err := chatImpl.SendMessage(ctx, request)

			require.Error(t, err)
			if e, ok := status.FromError(err); ok {
				require.Equal(t, e.Code(), test.expectedErrorCode)
			}
		})
	}
}
