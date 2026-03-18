package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/products"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (c *ProductsController) CreateModal(w http.ResponseWriter, r *http.Request) {
	_ = products.Create().Render(r.Context(), w)
}

func (c *ProductsController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, err)
		return
	}

	product, err := c.app.CreateProduct(ctx, commands.CreateProductRequest{
		OrganizationID: config.OrganizationID,
		Name:           r.PostForm.Get("name"),
	})

	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	w.Header().Set("HX-Location", "/dashboard/products/"+product.ID.Hex())
	w.WriteHeader(http.StatusOK)
}

func (c *ProductsController) Show(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("invalid product id: %w", err)))
		return
	}

	product, err := c.app.GetProductByID(ctx, productID, config.OrganizationID)
	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	_ = products.Show(product).Render(r.Context(), w)
}
