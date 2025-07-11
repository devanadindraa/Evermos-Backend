package address

import (
	"context"
	"fmt"
	"net/http"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/common"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	AddAddress(ctx *fiber.Ctx) error
	GetMyAddress(ctx *fiber.Ctx) error
	GetAddressByID(ctx *fiber.Ctx) error
	DeleteAddress(ctx *fiber.Ctx) error
	UpdateAddress(ctx *fiber.Ctx) error
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

	filter, err := common.GetMetaData(ctx, h.validate, "judul_alamat", "created_at_date", "updated_at_date")
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}
	result, err := h.service.GetMyAddress(reqCtx, filter)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", result)
	return nil
}

func (h *handler) GetAddressByID(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	addressID := ctx.Params("id")
	if addressID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}
	result, err := h.service.GetAddressByID(reqCtx, addressID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", result)
	return nil
}

func (h *handler) DeleteAddress(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	categoryID := ctx.Params("id")
	if categoryID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}
	err := h.service.DeleteAddress(reqCtx, categoryID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to DELETE data", nil)
	return nil
}

func (h *handler) UpdateAddress(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	addressID := ctx.Params("id")
	if addressID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}

	var input UpdateAddressReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return nil
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	res, err := h.service.UpdateAddress(reqCtx, input, addressID)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to PUT data", res)
	return nil
}
