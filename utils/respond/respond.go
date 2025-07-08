package respond

import (
	"fmt"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/gofiber/fiber/v2"
)

func Error(ctx *fiber.Ctx, err error) {
	apiErrors := apierror.GetApiErrors(err)

	resp := ApiModel[*string]{
		Status:  false,
		Message: "Failed to process request",
		Errors:  apiErrors.Messages,
		Data:    nil,
	}

	ctx.Set("error", err.Error())
	ctx.Status(apiErrors.Code).JSON(resp)
}

func Success(ctx *fiber.Ctx, code int, message string, data any) {
	ctx.Set("error", "")

	resp := ApiModel[any]{
		Status:  true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}

	ctx.Status(code).JSON(resp)
}

func Data(ctx *fiber.Ctx, param DataParam) {
	ctx.Set("error", "")

	if param.Data == nil {
		ctx.Status(param.Code)
		return
	}

	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", param.Filename))
	ctx.Set("Content-Type", param.MimeType)
	ctx.Status(param.Code).Send(param.Data)
}
