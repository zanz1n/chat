package testutils

import "math/rand/v2"

const (
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charset_lower = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func RandomString(size int) string {
	return string(RandomBytes(size))
}

func RandomStringLower(size int) string {
	return string(RandomBytesLower(size))
}

func RandomBytes(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return b
}

func RandomBytesLower(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = charset_lower[rand.IntN(len(charset_lower))]
	}
	return b
}
