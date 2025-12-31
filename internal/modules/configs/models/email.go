package models

type Email struct {
	Subject string `bson:"subject" json:"subject"`
	Body    string `bson:"body" json:"body"`
}
