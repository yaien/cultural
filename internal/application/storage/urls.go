package storage

import (
	"fmt"

	"github.com/yaien/cultural/internal/lib/primitive"
)

type URLFunc func(filename string, variant ...int) string

func FileURL(filename string, variant ...int) string {
	if len(variant) > 0 {
		return fmt.Sprintf("/assets/dynamic/files/%s?variant=%d", filename, variant[0])
	}
	return fmt.Sprintf("/assets/dynamic/files/%s", filename)
}

// NewExternalURLFunc creates a FileURLFunc that generates URLs for files served from the server's external assets endpoint.
func NewExternalURLFunc(serverURL string, organizationID primitive.ID) URLFunc {
	return func(filename string, variant ...int) string {
		if len(variant) > 0 {
			return fmt.Sprintf("%s/assets/external/%d/%s?variant=%d", serverURL, organizationID, filename, variant[0])
		}
		return fmt.Sprintf("%s/assets/external/%d/%s", serverURL, organizationID, filename)
	}
}
