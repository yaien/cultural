package storage

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLFunc func(filename string, variant ...int) string

func FileURL(filename string, variant ...int) string {
	if len(variant) > 0 {
		return fmt.Sprintf("/assets/dynamic/files/%s?variant=%d", filename, variant[0])
	}
	return fmt.Sprintf("/assets/dynamic/files/%s", filename)
}

// NewExternalURLFunc creates a FileURLFunc that generates URLs for files served from the server's external assets endpoint.
func NewExternalURLFunc(serverURL string, organizationID primitive.ObjectID) URLFunc {
	return func(filename string, variant ...int) string {
		if len(variant) > 0 {
			return fmt.Sprintf("%s/assets/external/%s/%s?variant=%d", serverURL, organizationID.Hex(), filename, variant[0])
		}
		return fmt.Sprintf("%s/assets/external/%s/%s", serverURL, organizationID.Hex(), filename)
	}
}
