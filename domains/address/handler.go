package address

import (
	"context"
	"net/http"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	AddAddress(ctx *fiber.Ctx) error
	GetMyAddress(ctx *fiber.Ctx) error
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

func (h *handler) AddAddress(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	var input AddressReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return nil
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	res, err := h.service.AddAddress(reqCtx, input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to POST data", res)
	return nil
}

func (h *handler) GetMyAddress(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	result, err := h.service.GetMyAddress(reqCtx)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", result)
	return nil
}
