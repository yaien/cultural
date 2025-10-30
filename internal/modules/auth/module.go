package auth

import "github.com/yaien/cultural/internal/infrastructure"

type Module struct {
}

func (m *Module) Init(mono *infrastructure.Monolith) error {
	return nil
}
