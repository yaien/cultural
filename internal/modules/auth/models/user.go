package models

import "time"

type User struct {
	ID        any       `bson:"_id" json:"id"`
	Email     string    `bson:"email" json:"email"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Authorization struct {
	Provider     string `bson:"provider" json:"provider"`
	UID          string `bson:"uid" json:"uid"`
	UserID       any    `bson:"userId" json:"userId"`
	AccessToken  string `bson:"accessToken" json:"accessToken"`
	RefreshToken string `bson:"refreshToken" json:"refreshToken"`
}
