package tests

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	repoMocks "github.com/mistandok/chat-server/internal/repository/mocks"
	"github.com/mistandok/chat-server/internal/service/chat"
	serviceMocks "github.com/mistandok/chat-server/internal/service/mocks"
	"github.com/mistandok/platform_common/pkg/db/mocks"
	"github.com/mistandok/platform_common/pkg/db/pg"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCreate_SuccessCreateChat(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	var chatID int64 = 1
	var userIDs = []int64{1, 2, 3}

	txFaker := serviceMocks.NewTxFaker(t)
	ctxWithTx := pg.MakeContextTx(ctx, txFaker)

	txFaker.On("Commit", ctxWithTx).Return(nil).Once()

	transactorMock := mocks.NewTransactor(t)
	transactorMock.On("BeginTx", ctx, txOpts).Return(txFaker, nil)

	txManagerMock := pg.NewTransactionManager(transactorMock)

	chatRepoMock := repoMocks.NewChatRepository(t)
	chatRepoMock.On("Create", ctxWithTx).Return(chatID, nil).Once()
	chatRepoMock.On("LinkChatAndUsers", ctxWithTx, chatID, userIDs).Return(nil).Once()

	userRepoMock := repoMocks.NewUserRepository(t)
	userRepoMock.On("CreateMass", ctxWithTx, userIDs).Return(nil).Once()

	messageRepoMock := repoMocks.NewMessageRepository(t)

	service := chat.NewService(&logger, txManagerMock, chatRepoMock, userRepoMock, messageRepoMock)
	resultChatID, err := service.Create(ctx, userIDs)

	require.NoError(t, err)
	require.Equal(t, chatID, resultChatID)
}
