package provcity

import (
	"fmt"
	"net/http"

	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetProvinces(ctx *fiber.Ctx) error
	GetCitys(ctx *fiber.Ctx) error
	GetDetailProvince(ctx *fiber.Ctx) error
	GetDetailCity(ctx *fiber.Ctx) error
}

type handler struct {
	provcityService Provcity
}

func NewProvcityHandler(provcity Provcity) Handler {
	return &handler{
		provcityService: provcity,
	}
}

func (h *handler) GetProvinces(ctx *fiber.Ctx) error {
	result, err := h.provcityService.GetListProvince()
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to get data", result)
	return nil
}

func (h *handler) GetCitys(ctx *fiber.Ctx) error {
	provID := ctx.Params("prov_id")
	if provID == "" {
		respond.Error(ctx, fmt.Errorf("prov_id is required"))
		return nil
	}
	result, err := h.provcityService.GetListCity(provID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to get data", result)
	return nil
}

func (h *handler) GetDetailProvince(ctx *fiber.Ctx) error {
	provID := ctx.Params("prov_id")
	if provID == "" {
		respond.Error(ctx, fmt.Errorf("prov_id is required"))
		return nil
	}
	result, err := h.provcityService.GetProvinceByID(provID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to get data", result)
	return nil
}

func (h *handler) GetDetailCity(ctx *fiber.Ctx) error {
	cityID := ctx.Params("city_id")
	if cityID == "" {
		respond.Error(ctx, fmt.Errorf("prov_id is required"))
		return nil
	}
	result, err := h.provcityService.GetCityByID(cityID)
	if err != nil {
		respond.Error(ctx, err)
		return nil
	}

	respond.Success(ctx, http.StatusOK, "Succeed to get data", result)
	return nil
}
