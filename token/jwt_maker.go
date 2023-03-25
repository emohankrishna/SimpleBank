package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: size must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (JWTMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		Payload:          *payload,
		RegisteredClaims: jwt.RegisteredClaims{},
	})
	signedString, err := jwtToken.SignedString([]byte(JWTMaker.secretKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func (JWTMaker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunction := func(unverifiedToken *jwt.Token) (interface{}, error) {
		if unverifiedToken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(JWTMaker.secretKey), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, keyFunction)
	if err != nil {
		return nil, ErrInvalidToken
	}
	/**

	if token.Valid {
		fmt.Println("You look nice today")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		fmt.Println("That's not even a token")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		fmt.Println("Timing is everything")
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

	*/
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		if claims.Payload.ExpiredAt.After(time.Now()) {
			return &claims.Payload, nil
		}
		return nil, ErrExpiredToken
	}
	return nil, ErrInvalidToken
}
