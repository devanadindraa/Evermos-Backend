package trx

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
	AddTrx(ctx *fiber.Ctx) error
	GetTrxByID(ctx *fiber.Ctx) error
	GetTrx(ctx *fiber.Ctx) error
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

func (h *handler) AddTrx(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	var input TrxReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return nil
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	res, err := h.service.AddTrx(reqCtx, input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to POST data", res)
	return nil
}

func (h *handler) GetTrxByID(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	trxID := ctx.Params("id")
	if trxID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}
	result, err := h.service.GetTrxByID(reqCtx, trxID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", result)
	return nil
}

func (h *handler) GetTrx(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)
	filter, err := common.GetMetaData(ctx, h.validate, "created_at_date", "updated_at_date")
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	result, err := h.service.GetTrx(reqCtx, filter)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET all trx", result)
	return nil
}
