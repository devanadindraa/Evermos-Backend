package context

import (
	"context"
	"net/http"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/constants"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func FiberWithCtx(fiberCtx *fiber.Ctx, ctx context.Context) *fiber.Ctx {
	ctxWithFiber := context.WithValue(ctx, FiberCtxKey, fiberCtx)

	combined := newCombinerCtx(
		newStopperCtx(context.Background()),
		newStopperCtx(ctxWithFiber),
	)

	fiberCtx.Locals("ctx", combined)
	return fiberCtx
}

func GetTokenClaims(ctx context.Context) (constants.Token, error) {
	val := ctx.Value(FiberCtxKey)
	if val == nil {
		return constants.Token{}, apierror.NewWarn(http.StatusInternalServerError, "Fiber context not found")
	}

	fiberCtx, ok := val.(*fiber.Ctx)
	if !ok {
		return constants.Token{}, apierror.NewWarn(http.StatusInternalServerError, "Invalid fiber context type")
	}

	tokenVal := fiberCtx.Locals("token")
	token, ok := tokenVal.(constants.Token)
	if !ok {
		return constants.Token{}, apierror.NewWarn(http.StatusInternalServerError, "Can't get token claims")
	}
	return token, nil
}

func SetTokenClaims(ctx context.Context, token constants.Token) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func GetRequestId(ctx context.Context) *uuid.UUID {
	reqId := ctx.Value(RequestIdKey)
	reqIdUUID, ok := reqId.(uuid.UUID)
	if !ok {
		return nil
	}
	return &reqIdUUID
}

func SetRequestId(ctx context.Context, requestId uuid.UUID) context.Context {
	return context.WithValue(ctx, RequestIdKey, requestId)
}
