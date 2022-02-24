package main

type User struct {
	Id          string `pg:"id,nopk,notnull,unique,type:uuid,default:gen_random_uuid()"`
	Email       string `pg:"email,nopk,notnull,unique"`
	Name        string `pg:"username,nopk,notnull,unique"`
	DisplayName string `pg:"displayName,nopk"`
}
