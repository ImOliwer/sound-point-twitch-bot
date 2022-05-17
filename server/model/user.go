package model

type User struct {
	tableName struct{} `bun:"users" json:"-"`
	ID        int64    `bun:"id,notnull,unique" json:"id"`
	Points    int64    `bun:"points,notnull,default:0" json:"points"`
}
