package category

import (
	"fmt"
	"net/http"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	AddCategory(ctx *fiber.Ctx) error
	GetAllCategory(ctx *fiber.Ctx) error
	GetCategoryByID(ctx *fiber.Ctx) error
	DeleteCategory(ctx *fiber.Ctx) error
	UpdateCategory(ctx *fiber.Ctx) error
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

func (h *handler) AddCategory(ctx *fiber.Ctx) error {
	var input CategoryReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return nil
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	res, err := h.service.AddCategory(ctx.Context(), input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to POST data", res)
	return nil
}

func (h *handler) GetAllCategory(ctx *fiber.Ctx) error {

	res, err := h.service.GetAllCategory(ctx.Context())
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", res)
	return nil
}

func (h *handler) GetCategoryByID(ctx *fiber.Ctx) error {
	categoryID := ctx.Params("id")
	if categoryID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}
	result, err := h.service.GetCategoryByID(ctx.Context(), categoryID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", result)
	return nil
}

func (h *handler) DeleteCategory(ctx *fiber.Ctx) error {
	categoryID := ctx.Params("id")
	if categoryID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}
	err := h.service.DeleteCategory(ctx.Context(), categoryID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to DELETE data", nil)
	return nil
}

func (h *handler) UpdateCategory(ctx *fiber.Ctx) error {
	categoryID := ctx.Params("id")
	if categoryID == "" {
		respond.Error(ctx, fmt.Errorf("id is required"))
		return nil
	}

	var input CategoryReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return nil
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	res, err := h.service.UpdateCategory(ctx.Context(), input, categoryID)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to PUT data", res)
	return nil
}
