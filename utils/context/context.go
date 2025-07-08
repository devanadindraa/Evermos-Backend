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
	combined := newCombinerCtx(
		newStopperCtx(context.Background()),
		newStopperCtx(ctx),
	)
	fiberCtx.Locals("ctx", combined)
	return fiberCtx
}

func GetTokenClaims(ctx context.Context) (constants.Token, error) {
	tokenVal := ctx.Value(tokenKey)
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
	reqId := ctx.Value(requestIdKey)
	reqIdUUID, ok := reqId.(uuid.UUID)
	if !ok {
		return nil
	}
	return &reqIdUUID
}

func setRequestId(ctx context.Context, requestId uuid.UUID) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}
