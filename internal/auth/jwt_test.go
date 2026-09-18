package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken("test-secret", 10, 20)
	require.NoError(t, err)

	claims, err := ParseToken("test-secret", token)
	require.NoError(t, err)
	require.Equal(t, uint(10), claims.StaffID)
	require.Equal(t, uint(20), claims.HospitalID)
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	token, err := GenerateToken("test-secret", 10, 20)
	require.NoError(t, err)

	_, err = ParseToken("wrong-secret", token)
	require.ErrorIs(t, err, ErrInvalidToken)
}
