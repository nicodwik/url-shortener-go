package entity

import "time"

type Redirection struct {
	ShortUrl  string     `gorm:"primarykey" json:"short_url"`
	LongUrl   string     `json:"long_url"`
	UserId    string     `json:"user_id"`
	Status    string     `json:"status"`
	HitCount  int        `json:"hit_count"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
