package entity

import "time"

type Customer struct {
	CifNo      int `gorm:"primaryKey"`
	NameKana   string
	NameKanji  string
	BirthDate  time.Time
	Prefecture string
	City       string
	Town       string
	Street     string
	Building   string
	Room       string
	Email      string
	Phone      string
	CreatedAt  time.Time
}
