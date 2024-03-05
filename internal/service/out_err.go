package service

import "errors"

const errMsgUserNotInTheChat = "пользователь не состоит в чате"

var ErrMsgUserNotInTheChat = errors.New(errMsgUserNotInTheChat) // ErrMsgUserNotInTheChat пользователь не состоит в чате
