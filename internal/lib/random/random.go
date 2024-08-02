package random

import (
	"math/rand"
	"time"
)

// NewRandomString generates random string with given length
func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvxyz" +
		"0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
