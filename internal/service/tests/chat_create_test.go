package tests

import (
	"context"
	"errors"
	"testing"

	repoMocks "github.com/mistandok/chat-server/internal/repository/mocks"
	"github.com/mistandok/chat-server/internal/service/chat"
	"github.com/mistandok/platform_common/pkg/db/pg"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCreate_SuccessCreateChat(t *testing.T) {
	startCtx := context.Background()
	logger := zerolog.Nop()
	var chatID int64 = 1
	var userIDs = []int64{1, 2, 3}

	txFaker, ctxWithTx := txFakerAndCtxWithSetup(startCtx, t, true)
	transactorMock := transactorWithSetup(startCtx, t, txFaker)
	txManagerMock := pg.NewTransactionManager(transactorMock)

	chatRepoMock := repoMocks.NewChatRepository(t)
	chatRepoMock.On("Create", ctxWithTx).Return(chatID, nil).Once()
	chatRepoMock.On("LinkChatAndUsers", ctxWithTx, chatID, userIDs).Return(nil).Once()

	userRepoMock := repoMocks.NewUserRepository(t)
	userRepoMock.On("CreateMass", ctxWithTx, userIDs).Return(nil).Once()

	messageRepoMock := repoMocks.NewMessageRepository(t)

	service := chat.NewService(&logger, txManagerMock, chatRepoMock, userRepoMock, messageRepoMock)
	resultChatID, err := service.Create(startCtx, userIDs)

	require.NoError(t, err)
	require.Equal(t, chatID, resultChatID)
}

func TestCreate_FailCreateChat(t *testing.T) {
	startCtx := context.Background()
	logger := zerolog.Nop()
	var userIDs = []int64{1, 2, 3}
	chatError := errors.New("some error")

	txFaker, ctxWithTx := txFakerAndCtxWithSetup(startCtx, t, false)
	transactorMock := transactorWithSetup(startCtx, t, txFaker)
	txManagerMock := pg.NewTransactionManager(transactorMock)

	chatRepoMock := repoMocks.NewChatRepository(t)
	chatRepoMock.On("Create", ctxWithTx).Return(0, chatError).Once()

	userRepoMock := repoMocks.NewUserRepository(t)
	messageRepoMock := repoMocks.NewMessageRepository(t)

	service := chat.NewService(&logger, txManagerMock, chatRepoMock, userRepoMock, messageRepoMock)
	_, err := service.Create(startCtx, userIDs)

	require.Error(t, err)
}

func TestCreate_FailCreateUsers(t *testing.T) {
	startCtx := context.Background()
	logger := zerolog.Nop()
	var chatID int64 = 1
	var userIDs = []int64{1, 2, 3}
	usersError := errors.New("some error")

	txFaker, ctxWithTx := txFakerAndCtxWithSetup(startCtx, t, false)
	transactorMock := transactorWithSetup(startCtx, t, txFaker)
	txManagerMock := pg.NewTransactionManager(transactorMock)

	chatRepoMock := repoMocks.NewChatRepository(t)
	chatRepoMock.On("Create", ctxWithTx).Return(chatID, nil).Once()

	userRepoMock := repoMocks.NewUserRepository(t)
	userRepoMock.On("CreateMass", ctxWithTx, userIDs).Return(usersError).Once()

	messageRepoMock := repoMocks.NewMessageRepository(t)

	service := chat.NewService(&logger, txManagerMock, chatRepoMock, userRepoMock, messageRepoMock)
	_, err := service.Create(startCtx, userIDs)

	require.Error(t, err)
}

func TestCreate_FailLinkChatAndUsers(t *testing.T) {
	startCtx := context.Background()
	logger := zerolog.Nop()
	var chatID int64 = 1
	var userIDs = []int64{1, 2, 3}
	linkError := errors.New("some error")

	txFaker, ctxWithTx := txFakerAndCtxWithSetup(startCtx, t, false)
	transactorMock := transactorWithSetup(startCtx, t, txFaker)
	txManagerMock := pg.NewTransactionManager(transactorMock)

	chatRepoMock := repoMocks.NewChatRepository(t)
	chatRepoMock.On("Create", ctxWithTx).Return(chatID, nil).Once()
	chatRepoMock.On("LinkChatAndUsers", ctxWithTx, chatID, userIDs).Return(linkError).Once()

	userRepoMock := repoMocks.NewUserRepository(t)
	userRepoMock.On("CreateMass", ctxWithTx, userIDs).Return(nil).Once()

	messageRepoMock := repoMocks.NewMessageRepository(t)

	service := chat.NewService(&logger, txManagerMock, chatRepoMock, userRepoMock, messageRepoMock)
	_, err := service.Create(startCtx, userIDs)

	require.Error(t, err)
}
