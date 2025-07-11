package product

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	AddProduct(ctx *fiber.Ctx) error
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
