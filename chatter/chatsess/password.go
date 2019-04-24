package chatsess

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

//NewPassword take a password string and returns a hashed password
func NewPassword(p string) string {
	salt := make([]byte, 10)
	rand.Read(salt)
	return password(p, salt)

}

//CheckPassword checks to see if a password matches its hash?
func CheckPassword(p, h string) bool {
	s := strings.Split(h, "_")[0]
	salt, _ := hex.DecodeString(s)

	nh := password(p, salt)
	return h == nh

}

func password(p string, s []byte) string { //s is salt
	key, _ := scrypt.Key([]byte(p), s, 32768, 8, 1, 32)
	return fmt.Sprintf("%x_%x", s, key)

}
