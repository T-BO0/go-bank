package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(9)

	passwordHash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, passwordHash)

	wrongPassword := RandomString(9)
	err = CheckPassword(wrongPassword, passwordHash)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	passwordHash2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, passwordHash2)
	require.NotEqual(t, passwordHash, passwordHash2)
}
