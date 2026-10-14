// package service

// import (
// 	"crypto/rand"
// 	"crypto/subtle"

// 	"golang.org/x/crypto/argon2"
// )

// func hashPassword(password string) (hash, salt []byte, err error) {
// 	salt = make([]byte, 16)
// 	if _, err = rand.Read(salt); err != nil {
// 		return nil, nil, err
// 	}
// 	hash = argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
// 	return hash, salt, nil
// }
// func verifyPassword(password string, hash, salt []byte) bool {
// 	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
// 	return subtle.ConstantTimeCompare(hash, computedHash) == 1
// }