package main

import (
	"github.com/dustin/go-coap"
)

func NewMessage(mType coap.COAPType, mCode coap.COAPCode, mID uint16, token, payload []byte) *coap.Message {
	return &coap.Message{
		Type:      mType,
		Code:      mCode,
		MessageID: mID,
		Token:     token,
		Payload:   payload,
	}
}
