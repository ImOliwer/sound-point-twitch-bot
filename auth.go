package main

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserToken struct {
	tableName struct{} `bun:"user_tokens" json:"-"`
	UserID    string   `bun:"user_id,notnull,type:uuid" json:"-"`
	Token     string   `bun:"token,notnull" json:"-"`
	CreatedAt string   `bun:"created_at,notnull,type:timestamp" json:"-"`
}

func (r *Application) GenerateJWT(item jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, item)
	return token.SignedString(r.Settings.Secrets.AccountAuth.Private)
}

func (r *Application) VerifyJWT(token string, of *User) (*jwt.Token, error) {
	parsed, parseErr := jwt.ParseWithClaims(token, of, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("incorrect signing method: %v", token.Header["alg"])
		}
		return r.Settings.Secrets.AccountAuth.Public, nil
	})
	if parseErr != nil {
		return nil, parseErr
	}
	return parsed, nil
}

func HashPassword(content string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(content), 14)
	return string(bytes), err
}

func ComparePassword(hashed string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

func (r *Application) TokenFor(user User) string {
	token, err := r.GenerateJWT(user)

	if err != nil {
		return ""
	}

	_, insErr := r.DbClient.NewInsert().Model(&UserToken{
		UserID: user.Id(),
		Token:  token,
	}).Exec(context.Background())

	if insErr != nil {
		return ""
	}

	return token
}
