package context

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func FiberWithCtx(fiberCtx *fiber.Ctx, ctx context.Context) *fiber.Ctx {
	combined := newCombinerCtx(
		newStopperCtx(context.Background()),
		newStopperCtx(ctx),
	)
	fiberCtx.Locals("ctx", combined)
	return fiberCtx
}
