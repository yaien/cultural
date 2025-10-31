package shared

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IDToStr(id any) string {
	switch v := id.(type) {
	case primitive.ObjectID:
		return v.Hex()
	case *primitive.ObjectID:
		if v != nil {
			return v.Hex()
		}
		return ""
	default:
		return ""
	}
}
