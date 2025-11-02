package landing

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/routes"
)

type Module struct {
}

func (m *Module) Init(mono *infrastructure.Monolith) error {
	routes.Register(mono)
	return nil

}
