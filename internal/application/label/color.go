package label

import (
	"fmt"

	"github.com/yaien/cultural/internal/lib/primitive"
)

type Color struct {
	ID    primitive.UUID
	Value string
	Tag   string
}

type Colors []*Color

func NewColor(colors []*Color) (*Color, error) {
loop:
	for index := range 100 {
		tag := fmt.Sprintf("tag-%s", primitive.ID(len(colors)+1+index))
		for _, color := range colors {
			if color.Tag == tag {
				continue loop
			}
		}

		color := &Color{
			ID:    primitive.NewUUID(),
			Tag:   tag,
			Value: "#000000",
		}

		return color, nil
	}

	return nil, fmt.Errorf("failed to generate unique key for color")
}
