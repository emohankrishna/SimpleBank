package token

import (
	"testing"
	"time"

	"github.com/emohankrishna/Simplebank/db/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	marker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomString(10)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	tokenSignedString, err := marker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenSignedString)
	payload, err := marker.VerifyToken(tokenSignedString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	marker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	tokenSignedString, err := marker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenSignedString)
	payload, err := marker.VerifyToken(tokenSignedString)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}

func TestInvalidJWTTokenAlgoNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		Payload: *payload,
	})
	tokenString, err := jwtToken.SignedString([]byte(jwt.UnsafeAllowNoneSignatureType))
	require.NoError(t, err)
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	payload, err = maker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
