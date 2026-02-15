package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views"
)

type ProductsController struct {
	app *application.Application
}

func NewProductsController(app *application.Application) *ProductsController {
	return &ProductsController{app: app}
}

func (c *ProductsController) Index(w http.ResponseWriter, r *http.Request) {
	views.Products(w, r)
}
