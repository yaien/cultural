package models

type Email struct {
	Subject  string `bson:"subject"`
	Body     string `bson:"body"`
	Category string `bson:"category"`
}
