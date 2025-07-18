package user

import (
	"context"
	"net/http"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Login(ctx *fiber.Ctx) error
	VerifyToken(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	UpdateProfile(ctx *fiber.Ctx) error
	GetProfile(ctx *fiber.Ctx) error
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

func (h *handler) Login(ctx *fiber.Ctx) error {
	var input LoginReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return nil
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	res, err := h.service.Login(ctx.Context(), input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to POST data", res)
	return nil
}

func (h *handler) VerifyToken(ctx *fiber.Ctx) error {
	respond.Success(ctx, http.StatusOK, "Succeed to POST data", VerifyTokenRes{TokenVerified: true})
	return nil
}

func (h *handler) Logout(ctx *fiber.Ctx) error {

	token, err := contextUtil.GetTokenClaims(ctx.Context())
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	input := LogoutReq{
		Token:   token.Token,
		Expires: token.Claims.ExpiresAt.Time,
	}

	res, err := h.service.Logout(ctx.Context(), input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to POST data", res)
	return nil
}

func (h *handler) Register(ctx *fiber.Ctx) error {
	var input RegisterReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return nil
	}

	if h.validate == nil {
		respond.Error(ctx, apierror.NewError(http.StatusInternalServerError, "Validator is nil"))
		return nil
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	if h.service == nil {
		respond.Error(ctx, apierror.NewError(http.StatusInternalServerError, "Service is nil"))
		return nil
	}

	res, err := h.service.Register(ctx.Context(), input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return nil
	}

	respond.Success(ctx, http.StatusCreated, "Succeed to POST data", res)
	return nil
}

func (h *handler) UpdateProfile(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)

	var input UpdateProfileReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return nil
	}

	res, err := h.service.UpdateProfile(reqCtx, input)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to PUT data", res)
	return nil
}

func (h *handler) GetProfile(ctx *fiber.Ctx) error {
	reqCtx := ctx.Locals("ctx").(context.Context)

	res, err := h.service.GetProfile(reqCtx)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to GET data", res)
	return nil
}
