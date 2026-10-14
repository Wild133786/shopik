package password_test

import (
	"encoding/base64"
	"fmt"
	"shopik/services/user/internal/password"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashFormat(t *testing.T) {
	const (
		wantTime            = 1
		wantMemory          = 64 * 1024
		wantThreads         = 4
		wantKeyLen          = 32
		wantSaltLen         = 16
		wantVersion         = 19
		wantHashParamsCount = 6
	)
	pass := "oleg14213"
	hashPassStr, err := password.Hash(pass)
	require.NoError(t, err, "password hash error")
	hashStrSlice := strings.Split(hashPassStr, "$")
	require.Equal(t, wantHashParamsCount, len(hashStrSlice))
	require.Equal(t, "argon2id", hashStrSlice[1])
	var (
		version,
		m,
		time,
		p int
	)
	fmt.Sscanf(hashStrSlice[2], "v=%d", &version)
	fmt.Sscanf(hashStrSlice[3], "m=%d,t=%d,p=%d", &m, &time, &p)
	hash, err := base64.RawStdEncoding.DecodeString(hashStrSlice[5])
	require.NoError(t, err)
	salt, err := base64.RawStdEncoding.DecodeString(hashStrSlice[4])
	require.NoError(t, err)
	assert.Equal(t, wantVersion, version)
	assert.Equal(t, wantMemory, m)
	assert.Equal(t, wantTime, time)
	assert.Equal(t, wantThreads, p)
	assert.Equal(t, wantKeyLen, len(hash))
	assert.Equal(t, wantSaltLen, len(salt))
}

func TestHashNotIdempotentHashAndSalt(t *testing.T) {
	const wantHashParamsCount = 6
	pass := "oleg14213"

	firstHashPassStr, err := password.Hash(pass)
	require.NoError(t, err)
	firstHashStrSlice := strings.Split(firstHashPassStr, "$")
	require.Equal(t, wantHashParamsCount, len(firstHashStrSlice))

	secondHashPassStr, err := password.Hash(pass)
	require.NoError(t, err)
	secondHashStrSlice := strings.Split(secondHashPassStr, "$")
	require.Equal(t, wantHashParamsCount, len(secondHashStrSlice))

	assert.NotEqual(t, firstHashStrSlice[4], secondHashStrSlice[4])
	assert.NotEqual(t, firstHashStrSlice[5], secondHashStrSlice[5])
}

func TestVerifySuccess(t *testing.T) {
	pass := "oleg228"
	computedHash := "$argon2id$v=19$m=65536,t=1,p=4$6gGS5CiokEYOQrmfqHki1A$mLlUyoa3xBcUU3WiPEYnnH8l2Gi7C4qYYQugJcXuPiA"
	hashPassStr, err := password.Hash(pass)
	require.NoError(t, err)
	ok, err := password.Verify(pass, hashPassStr)
	require.NoError(t, err)
	require.Equal(t, true, ok)
	ok2, err := password.Verify(pass, computedHash)
	require.NoError(t, err)
	require.Equal(t, ok2, ok)
}

func TestVerify_InvalidHash(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
	}{
		{"empty string", ""},
		{"not a hash", "garbage"},
		{"too few parts", "$argon2id$v=19$m=65536,t=1,p=4$c2FsdA"},
		{"wrong algorithm", "$argon2i$v=19$m=65536,t=1,p=4$c2FsdA$aGFzaA"},
		{"wrong version", "$argon2id$v=18$m=65536,t=1,p=4$c2FsdA$aGFzaA"},
		{"broken params", "$argon2id$v=19$m=abc,t=1,p=4$c2FsdA$aGFzaA"},
		{"zero iterations", "$argon2id$v=19$m=65536,t=0,p=4$c2FsdA$aGFzaA"},
		{"zero threads", "$argon2id$v=19$m=65536,t=1,p=0$c2FsdA$aGFzaA"},
		{"invalid base64 salt", "$argon2id$v=19$m=65536,t=1,p=4$!!!$aGFzaA"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := password.Verify("any-pass", tc.encoded)
			assert.Error(t, err)
			assert.False(t, ok)
		})
	}
}
