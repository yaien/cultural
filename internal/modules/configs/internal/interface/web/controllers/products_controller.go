package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/coderror"
	"github.com/yaien/cultural/internal/label"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/dashboard"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/products"
	"github.com/yaien/cultural/internal/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductsController struct {
	products      *store.Products
	presentations *store.Presentations
	files         *store.Files
}

func NewProductsController(products *store.Products, presentations *store.Presentations, files *store.Files) *ProductsController {
	return &ProductsController{products, presentations, files}
}

func (c *ProductsController) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)
	prs, err := c.products.GetByOrganizationID(ctx, config.OrganizationID)
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
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, err)
		return
	}

	opts := &store.CreateProductOptions{
		OrganizationID: config.OrganizationID,
		Name:           r.PostForm.Get("name"),
	}

	product, err := c.products.Create(ctx, opts)
	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	w.Header().Set("HX-Location", "/dashboard/products/"+product.ID.Hex())
	w.WriteHeader(http.StatusOK)
}

func (c *ProductsController) Show(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid product id: %w", err))
		return
	}

	product, err := c.products.GetByIDAndOrganizationID(ctx, productID, config.OrganizationID)
	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	var presentation *store.Presentation
	if pid, err := primitive.ObjectIDFromHex(r.URL.Query().Get("presentation")); err == nil {
		for _, p := range product.Presentations {
			if p.ID == pid {
				presentation = p
				break
			}
		}

		if presentation == nil {
			WriteHTMLErr(w, coderror.Newf("presentation_not_found", "presentation with id %s not found", pid.Hex()))
			return
		}
	} else if len(product.Presentations) > 0 {
		presentation = product.Presentations[0]
	}

	if r.Header.Get("HX-Target") == products.PresentationsID {
		_ = templ.Join(
			products.Presentations(product, presentation),
			products.Pictures(product, presentation, products.SWAPOOB),
		).Render(ctx, w)

		return
	}

	_ = products.Show(product, presentation).Render(r.Context(), w)
}

func (c *ProductsController) CreatePresentation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid product id: %w", err))
		return
	}

	opts := &store.CreatePresentationOptions{
		OrganizationID: config.OrganizationID,
		ProductID:      productID,
	}

	product, presentation, err := c.presentations.Create(ctx, opts)
	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	_ = templ.Join(
		products.Presentations(product, presentation),
		products.Pictures(product, presentation, products.SWAPOOB),
	).Render(ctx, w)
}

func (c *ProductsController) UpdatePresentation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid product id: %w", err))
		return
	}

	presentationID, err := primitive.ObjectIDFromHex(r.PathValue("presentationId"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid presentation id: %w", err))
		return
	}

	name := r.PostFormValue("name")
	quantity, _ := strconv.Atoi(r.PostFormValue("quantity"))
	price, _ := strconv.ParseFloat(r.PostFormValue("price"), 64)

	product, presentation, err := c.presentations.Update(ctx, &store.UpdatePresentationOptions{
		PresentationID: presentationID,
		ProductID:      productID,
		OrganizationID: config.OrganizationID,
		Name:           name,
		Quantity:       quantity,
		Price:          price,
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
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid product id: %w", err))
		return
	}

	presentationID, err := primitive.ObjectIDFromHex(r.PathValue("presentationId"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid presentation id: %w", err))
		return
	}

	opts := &store.DeletePresentationOptions{
		ID:             presentationID,
		ProductID:      productID,
		OrganizationID: config.OrganizationID,
	}

	product, err := c.presentations.Delete(ctx, opts)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("error deleting presentation: %w", err))
		return
	}

	var presentation *store.Presentation
	if len(product.Presentations) > 0 {
		presentation = product.Presentations[0]
	}

	_ = templ.Join(
		products.Presentations(product, presentation),
		products.Pictures(product, presentation, products.SWAPOOB),
	).Render(ctx, w)

}

func (c *ProductsController) UploadPresentationFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	productID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid product id: %w", err))
		return
	}

	presentationID, err := primitive.ObjectIDFromHex(r.PathValue("presentationId"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid presentation id: %w", err))
		return
	}

	file, fileheader, err := r.FormFile("file")
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("error retrieving file: %w", err))
		return
	}

	opts := &store.UploadFileOptions{
		PresentationID: presentationID,
		ProductID:      productID,
		OrganizationID: config.OrganizationID,
		Name:           fileheader.Filename,
		Size:           fileheader.Size,
		ContentType:    fileheader.Header.Get("Content-Type"),
		Data:           file,
	}

	product, presentation, err := c.files.Upload(ctx, opts)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("error adding presentation picture: %w", err))
		return
	}

	_ = templ.Join(
		products.Pictures(product, presentation),
		dashboard.Toast("Archivo subido", dashboard.Primary),
	).Render(ctx, w)

}
