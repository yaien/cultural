package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/products"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type ProductsController struct {
	app *application.Application
}

func NewProductsController(app *application.Application) *ProductsController {
	return &ProductsController{app: app}
}

func (c *ProductsController) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)
	prs, err := c.app.GetProducts(ctx, config.OrganizationID)
	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	_ = products.Page(prs).Render(r.Context(), w)
}
