package shop

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/devanadindraa/Evermos-Backend/utils/common"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetMyShop(ctx *fiber.Ctx) error
	GetShopByID(ctx *fiber.Ctx) error
	UpdateMyShop(ctx *fiber.Ctx) error
	GetAllShop(ctx *fiber.Ctx) error
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

func (h *handler) GetMyShop(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	result, err := h.service.GetMyShop(reqCtx)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", result)
	return nil
}

func (h *handler) GetShopByID(ctx *fiber.Ctx) error {
	shopID := ctx.Params("id_toko")
	if shopID == "" {
		respond.Error(ctx, fmt.Errorf("id_toko is required"))
		return nil
	}
	result, err := h.service.GetShopByID(ctx.Context(), shopID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", result)
	return nil
}

func (h *handler) UpdateMyShop(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	shopID := ctx.Params("id_toko")

	form, err := ctx.MultipartForm()
	if err != nil {
		respond.Error(ctx, fmt.Errorf("failed to parse multipart form: %v", err))
		return nil
	}

	namaToko := form.Value["nama_toko"]
	var foto *multipart.FileHeader
	if len(form.File["photo"]) > 0 {
		foto = form.File["photo"][0]
	}

	req := UpdateShopReq{
		NamaToko: nil,
		UrlFoto:  nil,
	}

	if len(namaToko) > 0 {
		req.NamaToko = &namaToko[0]
	}
	if foto != nil {
		req.UrlFoto = foto
	}

	result, err := h.service.UpdateMyShop(reqCtx, req, shopID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to UPDATE shop", result)
	return nil
}

func (h *handler) GetAllShop(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	filter, err := common.GetMetaData(ctx, h.validate, "nama_toko", "created_at_date", "updated_at_date")
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	result, err := h.service.GetAllShop(reqCtx, filter)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET all shop", result)
	return nil
}
