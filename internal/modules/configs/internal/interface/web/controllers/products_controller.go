package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/dashboard"
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

	var presentation *models.Presentation
	if pid, err := primitive.ObjectIDFromHex(r.URL.Query().Get("presentation")); err == nil {
		for _, p := range product.Presentations {
			if p.ID == pid {
				presentation = p
				break
			}
		}

		if presentation == nil {
			WriteHTMLErr(w, &models.Error{Code: "presentation_not_found", Err: fmt.Errorf("presentation with id %s not found", pid.Hex())})
			return
		}
	} else if len(product.Presentations) > 0 {
		presentation = product.Presentations[0]
	}

	if r.Header.Get("HX-Target") == products.PresentationsID {
		_ = templ.Join(
			products.Presentations(product, presentation),
			products.Pictures(presentation, products.SWAPOOB),
		).Render(ctx, w)

		return
	}

	_ = products.Show(product, presentation).Render(r.Context(), w)
}

func (c *ProductsController) CreatePresentation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("invalid product id: %w", err)))
		return
	}

	product, presentation, err := c.app.CreateProductPresentation(ctx, productID, config.OrganizationID)
	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	_ = templ.Join(
		products.Presentations(product, presentation),
		products.Pictures(presentation, products.SWAPOOB),
	).Render(ctx, w)
}

func (c *ProductsController) UpdatePresentation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("invalid product id: %w", err)))
		return
	}

	presentationID, err := primitive.ObjectIDFromHex(r.PathValue("presentationId"))
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("invalid presentation id: %w", err)))
		return
	}

	name := r.PostFormValue("name")
	quantity, _ := strconv.Atoi(r.PostFormValue("quantity"))
	price, _ := strconv.ParseFloat(r.PostFormValue("price"), 64)

	product, presentation, err := c.app.UpdateProductPresentation(ctx, commands.UpdateProductPresentationRequest{
		ID:             presentationID,
		Name:           name,
		Quantity:       quantity,
		Price:          price,
		ProductID:      productID,
		OrganizationID: config.OrganizationID,
	})

	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("error updating presentation: %w", err))
		return
	}

	_ = templ.Join(
		products.Presentations(product, presentation),
		dashboard.Toast("Presentación de producto guardada", dashboard.Primary),
	).Render(ctx, w)
}

func (c *ProductsController) DeletePresentation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("invalid product id: %w", err)))
		return
	}

	presentationID, err := primitive.ObjectIDFromHex(r.PathValue("presentationId"))
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("invalid presentation id: %w", err)))
		return
	}

	err = c.app.DeleteProductPresentation(ctx, commands.DeleteProductPresentationRequest{
		ID:             presentationID,
		ProductID:      productID,
		OrganizationID: config.OrganizationID,
	})

	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("error deleting presentation: %w", err))
		return
	}

	product, err := c.app.GetProductByID(ctx, productID, config.OrganizationID)
	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	var presentation *models.Presentation
	if len(product.Presentations) > 0 {
		presentation = product.Presentations[0]
	}

	_ = templ.Join(
		products.Presentations(product, presentation),
		products.Pictures(presentation, products.SWAPOOB),
	).Render(ctx, w)

}
