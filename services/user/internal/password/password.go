package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (result string, err error) {
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

func ComparePassword(password string, encodedHash string) (bool, error) { //logs needed
	const (
		wantHashParamsCount = 6
		wantAlg             = "argon2id"
		wantVersion         = 19
		minTime             = 1
		maxTime             = 10
		minMemory           = 10 * 1024
		maxMemory           = 1024 * 1024
		minThreads          = 1
		maxThreads          = 16
		minSaltLen          = 16
		maxSaltLen          = 32
		wantKeyLen          = 32
	)
	hashStrSlice := strings.Split(encodedHash, "$")
	if len(hashStrSlice) != wantHashParamsCount {
		return false, fmt.Errorf("not enough hash+salt param, want: %d, got: %d", wantHashParamsCount, len(hashStrSlice))
	}
	var (
		version,
		keyLen int
		m uint32
		t uint32
		p uint8
	)
	alg := hashStrSlice[1]
	if alg != wantAlg {
		return false, fmt.Errorf("invalid alg, want: %s, got: %s", wantAlg, alg)
	}
	_, err := fmt.Sscanf(hashStrSlice[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != wantVersion {
		return false, fmt.Errorf("invalid argon2id version, want: %d, got: %d", wantVersion, version)
	}
	_, err = fmt.Sscanf(hashStrSlice[3], "m=%d,t=%d,p=%d", &m, &t, &p)
	if err != nil {
		return false, err
	}
	if m < minMemory || m > maxMemory {
		return false, fmt.Errorf("invalid memory param: %d", m)
	} else if t < minTime || t > maxTime {
		return false, fmt.Errorf("invalid time param: %d", t)
	} else if p < minThreads || p > maxThreads {
		return false, fmt.Errorf("invalid threads param: %d", p)
	}
	salt, err := base64.RawStdEncoding.DecodeString(hashStrSlice[4])
	if err != nil {
		return false, err
	} else if len(salt) < minSaltLen || len(salt) > maxSaltLen {
		return false, fmt.Errorf("invalid salt len: %d", len(salt))
	}
	expecctedHash, err := base64.RawStdEncoding.DecodeString(hashStrSlice[5])
	if err != nil {
		return false, err
	}
	keyLen = len(expecctedHash)
	if keyLen != wantKeyLen {
		return false, fmt.Errorf("invalid byte hash len, want: %d, got: %d", wantKeyLen, keyLen)
	}
	computedHash := argon2.IDKey([]byte(password), salt, t, m, p, uint32(keyLen))
	return subtle.ConstantTimeCompare(expecctedHash, computedHash) == 1, nil
}
