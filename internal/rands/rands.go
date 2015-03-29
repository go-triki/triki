/*
Package rands generates pseudo-random []byte slices.
*/
package rands // import "gopkg.in/triki.v0/internal/rands"

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const randChars = "`1234567890-=qwertyuiop[]asdfghjkl;'\\<zxcvbnm,./~!@#$%^&*()_+QWERTYUIOP{}ASDFGHJKL:\"|>ZXCVBNM<>?"

// New returns new random slice.
func New(n int) []byte {
	salt := make([]byte, n)
	for i := 0; i < n; i++ {
		salt[i] = randChars[rand.Intn(len(randChars))]
	}
	return salt
}
