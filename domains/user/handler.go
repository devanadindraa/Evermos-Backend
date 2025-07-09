package user

import (
	"net/http"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Login(ctx *fiber.Ctx)
	VerifyToken(ctx *fiber.Ctx)
	Logout(ctx *fiber.Ctx)
	Register(ctx *fiber.Ctx)
	UpdateProfile(ctx *fiber.Ctx)
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

func (h *handler) Login(ctx *fiber.Ctx) {
	var input LoginReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	res, err := h.service.Login(ctx.Context(), input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, "Succeed to POST data", res)
}

func (h *handler) VerifyToken(ctx *fiber.Ctx) {
	respond.Success(ctx, http.StatusOK, "Succeed to POST data", VerifyTokenRes{TokenVerified: true})
}

func (h *handler) Logout(ctx *fiber.Ctx) {

	token, err := contextUtil.GetTokenClaims(ctx.Context())
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	input := LogoutReq{
		Token:   token.Token,
		Expires: token.Claims.ExpiresAt.Time,
	}

	res, err := h.service.Logout(ctx.Context(), input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, "Succeed to POST data", res)
}

func (h *handler) Register(ctx *fiber.Ctx) {

	var input RegisterReq
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	res, err := h.service.Register(ctx.Context(), input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusCreated, "Succeed to POST data", res)
}

func (h *handler) UpdateProfile(ctx *fiber.Ctx) {
	var input UpdateProfileReq

	// Bind JSON ke struct input
	if err := ctx.BodyParser(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	// Validasi input menggunakan validator
	if err := h.validate.Struct(input); err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	// Panggil service untuk update profile
	res, err := h.service.UpdateProfile(ctx.Context(), input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	// Respon sukses
	respond.Success(ctx, http.StatusOK, "Succeed to UPDATE data", res)
}
