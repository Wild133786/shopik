package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func Hash(password string) (result string, err error) {
	const (
		time            = 1
		memory          = 64 * 1024
		threads         = 4
		keyLen          = 32
		saltLen         = 16
		hashParamsCount = 6
	)
	salt := make([]byte, saltLen)
	if _, err = rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)
	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	result = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, threads, saltBase64, hashBase64)
	return result, nil

}

func Verify(password string, encodedHash string) (bool, error) {
	const (
		wantSaltLen         = 16
		wantHashParamsCount = 6
		wantTime            = 1
		wantMemory          = 64 * 1024
		wantThreads         = 4
		wantAlg             = "argon2id"
		wantVersion         = 19
		wantKeyLen          = 32
	)
	hashStrSlice := strings.Split(encodedHash, "$")
	if len(hashStrSlice) != wantHashParamsCount {
		return false, fmt.Errorf("not enought hash+salt param, want: %d, got: %d", wantHashParamsCount, len(hashStrSlice))
	}
	var (
		version,
		m,
		t,
		p,
		keyLen int
	)
	alg := hashStrSlice[1]
	if alg != wantAlg {
		return false, fmt.Errorf("ivalid alg, want: %s, got: %s", wantAlg, alg)
	}
	fmt.Sscanf(hashStrSlice[2], "v=%d", &version)
	if version != wantVersion {
		return false, fmt.Errorf("invalid argon2id version, want: %d, got: %d", wantVersion, version)
	}
	fmt.Sscanf(hashStrSlice[3], "m=%d,t=%d,p=%d", &m, &t, &p)
	if m != wantMemory {
		return false, fmt.Errorf("invalid memory param, want: %d, got: %d", wantMemory, m)
	} else if t != wantTime {
		return false, fmt.Errorf("invalid time param, want: %d, got: %d", wantTime, t)
	} else if p != wantThreads {
		return false, fmt.Errorf("invalid threads param, want: %d, got: %d", wantThreads, p)
	}
	salt, err := base64.RawStdEncoding.DecodeString(hashStrSlice[4])
	if err != nil {
		return false, err
	} else if len(salt) != wantSaltLen {
		return false, fmt.Errorf("invalid salt len, want: %d, got: %d", wantSaltLen, len(salt))
	}
	expecctedHash, err := base64.RawStdEncoding.DecodeString(hashStrSlice[5])
	if err != nil {
		return false, err
	}
	keyLen = len(expecctedHash)
	if keyLen != wantKeyLen {
		return false, fmt.Errorf("invalid byte hahs len, want: %d, got: %d", wantKeyLen, keyLen)
	}
	computedHash := argon2.IDKey([]byte(password), salt, uint32(t), uint32(m), uint8(p), uint32(keyLen))
	return subtle.ConstantTimeCompare(expecctedHash, computedHash) == 1, nil
}
