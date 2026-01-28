package entity

type Client struct {
	ClientID     string `gorm:"primaryKey"`
	ClientSecret string
	ClientName   string
	Scope        string
}
