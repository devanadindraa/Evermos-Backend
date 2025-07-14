package product

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/common"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	AddProduct(ctx *fiber.Ctx) error
	GetProductByID(ctx *fiber.Ctx) error
	DeleteProduct(ctx *fiber.Ctx) error
	UpdateProduct(ctx *fiber.Ctx) error
	GetProducts(ctx *fiber.Ctx) error
}

type handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service, validate *validator.Validate) Handler {
	return &handler{
		service:  service,
		validate: validate,
	}
}

func (h *handler) AddProduct(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)

	var req ProductReq
	if err := ctx.BodyParser(&req); err != nil {
		respond.Error(ctx, err)
		return nil
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		respond.Error(ctx, fmt.Errorf("failed to parse multipart form: %v", err))
		return nil
	}

	files := form.File["photos"]
	req.Photos = append([]*multipart.FileHeader{}, files...)

	result, err := h.service.AddProduct(reqCtx, req)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to POST product", result)
	return nil
}

func (h *handler) GetProductByID(ctx *fiber.Ctx) error {
	productID := ctx.Params("id")
	if productID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}
	result, err := h.service.GetProductByID(ctx.Context(), productID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", result)
	return nil
}

func (h *handler) DeleteProduct(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	productID := ctx.Params("id")
	if productID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}
	err := h.service.DeleteProduct(reqCtx, productID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to DELETE data", nil)
	return nil
}

func (h *handler) UpdateProduct(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	productID := ctx.Params("id")
	if productID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		respond.Error(ctx, fmt.Errorf("failed to parse multipart form: %v", err))
		return nil
	}

	var input UpdateProductReq
	if v := form.Value["nama_produk"]; len(v) > 0 {
		input.NamaProduk = &v[0]
	}
	if v := form.Value["slug"]; len(v) > 0 {
		input.Slug = &v[0]
	}
	if v := form.Value["id_category"]; len(v) > 0 {
		id, err := strconv.Atoi(v[0])
		if err == nil {
			idCategory := int(id)
			input.IdCategory = &idCategory
		}
	}
	if v := form.Value["harga_reseller"]; len(v) > 0 {
		input.HargaReseller = &v[0]
	}
	if v := form.Value["harga_konsumen"]; len(v) > 0 {
		input.HargaKonsumen = &v[0]
	}
	if v := form.Value["stok"]; len(v) > 0 {
		stok, err := strconv.Atoi(v[0])
		if err == nil {
			input.Stok = &stok
		}
	}
	if v := form.Value["deskripsi"]; len(v) > 0 {
		input.Deskripsi = &v[0]
	}

	files := form.File["photos"]
	if len(files) > 0 {
		input.Photos = &files
	}

	if err := h.validate.Struct(input); err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	res, err := h.service.UpdateProduct(reqCtx, input, productID)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to PUT data", res)
	return nil
}

func (h *handler) GetProducts(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)

	meta, err := common.GetMetaData(ctx, h.validate, "id", "nama_produk", "harga_konsumen", "created_at_date")
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	var req GetProductReq
	req.FilterReq = meta

	if ctx.Query("category_id") != "" {
		categoryID, err := strconv.Atoi(ctx.Query("category_id"))
		if err == nil {
			cid := uint(categoryID)
			req.CategoryID = &cid
		}
	}

	if ctx.Query("toko_id") != "" {
		tokoID, err := strconv.Atoi(ctx.Query("toko_id"))
		if err == nil {
			tid := uint(tokoID)
			req.TokoID = &tid
		}
	}

	if ctx.Query("min_harga") != "" {
		minHarga, err := strconv.Atoi(ctx.Query("min_harga"))
		if err == nil {
			req.MinHarga = &minHarga
		}
	}

	if ctx.Query("max_harga") != "" {
		maxHarga, err := strconv.Atoi(ctx.Query("max_harga"))
		if err == nil {
			req.MaxHarga = &maxHarga
		}
	}

	// Call service
	res, err := h.service.GetProducts(reqCtx, req)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Success", res)
	return nil
}
