package category

import (
	"net/http"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	AddCategory(ctx *fiber.Ctx) error
	GetAllCategory(ctx *fiber.Ctx) error
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
