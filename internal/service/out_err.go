package service

import "errors"

const errMsgUserNotInTheChat = "пользователь не состоит в чате"
const errMsgChatNotFound = "чат не найден"
const errMsgCantCreateChatConnection = "не удалось создать соединение с чатом"

var ErrUserNotInTheChat = errors.New(errMsgUserNotInTheChat)                 // ErrUserNotInTheChat пользователь не состоит в чате
var ErrChatNotFound = errors.New(errMsgChatNotFound)                         // ErrChatNotFound чат не найден
var ErrCantCreateChatConnection = errors.New(errMsgCantCreateChatConnection) // ErrCantCreateChatConnection не удалось создать соединение с чатом
