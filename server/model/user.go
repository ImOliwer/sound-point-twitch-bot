package model

type User struct {
	tableName struct{} `bun:"users" json:"-"`
	ID        uint64   `bun:"id,pk,notnull,unique" json:"id"`
	Points    uint64   `bun:"points,notnull,default:0" json:"points"`
}
