package label

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Color struct {
	ID    primitive.ObjectID `bson:"_id"`
	Value string             `bson:"value"`
	Tag   string             `bson:"tag"`
}

type Colors []*Color

func NewColor(colors []*Color) (*Color, error) {
loop:
	for index := range 100 {
		tag := fmt.Sprintf("color-%d", len(colors)+1+index)
		for _, color := range colors {
			if color.Tag == tag {
				continue loop
			}
		}

		color := &Color{
			ID:    primitive.NewObjectID(),
			Tag:   tag,
			Value: "#000000",
		}

		return color, nil
	}

	return nil, fmt.Errorf("failed to generate unique key for color")
}
