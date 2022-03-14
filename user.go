package main

import "github.com/golang-jwt/jwt"

type User struct {
	tableName          struct{} `bun:"users" json:"-"`
	jwt.StandardClaims `bun:"-" json:"tokenOpts"`
	id                 string `bun:"id,notnull,unique,type:uuid,default:gen_random_uuid()" json:"id"`
	Email              string `bun:"email,notnull,unique" json:"email"`
	Name               string `bun:"username,notnull,unique" json:"name"`
	DisplayName        string `bun:"displayName" json:"displayName"`
	Password           string `bun:"password,notnull" json:"-"`
}

func (r *User) Id() string {
	return r.id
}
