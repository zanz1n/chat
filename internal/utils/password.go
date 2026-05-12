package utils

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var HashPasswordCost = ""

var bcryptCost = 12

func init() {
	if HashPasswordCost != "" {
		nc, err := strconv.Atoi(HashPasswordCost)
		if err != nil {
			return
		}
		bcryptCost = nc
	}
}

func HashPassword(passwd string) []byte {
	h, err := bcrypt.GenerateFromPassword([]byte(passwd), bcryptCost)
	if err != nil {
		panic(fmt.Errorf("hash password failed: %w", err))
	}
	return h
}

func CheckPasswordHash(hash []byte, passwd string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(passwd))
	return err == nil
}
