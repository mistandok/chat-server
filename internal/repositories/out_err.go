package repositories

import "github.com/pkg/errors"

const (
	errMsgChatNotFound     = "чата не существует"
	errMsgUserNotFound     = "пользователя не существует"
	errMsgUserNotInTheChat = "пользователь не состоит в чате"
)

var (
	ErrUserNotFound     = errors.New(errMsgUserNotFound)     // ErrUserNotFound сигнальная ошибка в случае отсутствия пользователя.
	ErrChatNotFound     = errors.New(errMsgChatNotFound)     // ErrChatNotFound сигнальная ошибка в случае отсутствия чата.
	ErrUserNotInTheChat = errors.New(errMsgUserNotInTheChat) // ErrUserNotInTheChat сигнальная ошибка в случае отсутствия пользователя в чате.
)
