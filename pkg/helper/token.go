package helper

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type GenerateAccessTokenArgs struct {
	UserID, TokenDuration uint
	SecretKey             string
}

func GenerateAccessToken(args *GenerateAccessTokenArgs) (string, error) {
	willExpiredAt := time.Now().Add(time.Duration(args.TokenDuration) * time.Minute)

	claims := jwt.MapClaims{}
	claims["user_id"] = args.UserID
	claims["issued_at"] = time.Now()
	claims["exp"] = willExpiredAt.Unix()

	signature := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := signature.SignedString([]byte(args.SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
