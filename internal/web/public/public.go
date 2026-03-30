package public

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/web/public/routes"
)

func Register(mono *infrastructure.Monolith) error {
	routes.Register(mono)
	return nil

}
