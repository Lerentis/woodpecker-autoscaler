package utils

import (
	"math/rand"

	log "github.com/sirupsen/logrus"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func BoolPointer(b bool) *bool {
	return &b
}

func CheckError(err error, caller string) {
	if err != nil {
		log.WithFields(log.Fields{
			"Caller": caller,
		}).Warnf("Error from hetzner API: %s", err.Error())
	}
}
