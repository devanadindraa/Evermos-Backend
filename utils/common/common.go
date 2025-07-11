package common

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/constants"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetMetaData(ctx *fiber.Ctx, validate *validator.Validate, allowedColumns ...string) (res *constants.FilterReq, err error) {
	limit, err := strconv.Atoi(ctx.Query(constants.QUERY_PARAMS_LIMIT, "20"))
	if err != nil {
		return nil, apierror.NewWarn(http.StatusBadRequest, "Limit must be a number")
	}

	page, err := strconv.Atoi(ctx.Query(constants.QUERY_PARAMS_PAGE, "1"))
	if err != nil {
		return nil, apierror.NewWarn(http.StatusBadRequest, "Page must be a number")
	}

	orderBy := ctx.Query(constants.QUERY_PARAMS_ORDER_BY, allowedColumns[0])
	if !slices.Contains(allowedColumns, orderBy) {
		return nil, apierror.NewWarn(http.StatusBadRequest, fmt.Sprintf("Order by column '%s' is not allowed!", orderBy))
	}

	sortOrder := strings.ToLower(ctx.Query(constants.QUERY_PARAMS_SORT_ORDER, "asc"))
	keyword := ctx.Query(constants.QUERY_PARAMS_KEYWORD)
	if keyword == "" {
		keyword = ctx.Query("nama")
	}
	if keyword == "" {
		keyword = ctx.Query("search")
	}
	if keyword == "" {
		keyword = ctx.Query("q")
	}
	if keyword == "" {
		keyword = ctx.Query("judul_alamat")
	}

	startCreatedAtStr := ctx.Query(constants.QUERY_PARAMS_START_CREATED_AT)
	endCreatedAtStr := ctx.Query(constants.QUERY_PARAMS_END_CREATED_AT)
	startUpdatedAtStr := ctx.Query(constants.QUERY_PARAMS_START_UPDATED_AT)
	endUpdatedAtStr := ctx.Query(constants.QUERY_PARAMS_END_UPDATED_AT)

	var startCreatedAt, endCreatedAt, startUpdatedAt, endUpdatedAt *time.Time

	if startCreatedAtStr != "" {
		temp, err := time.Parse(time.RFC3339Nano, startCreatedAtStr)
		if err != nil {
			return nil, apierror.NewWarn(http.StatusBadRequest, err.Error())
		}
		startCreatedAt = &temp
	}

	if endCreatedAtStr != "" {
		temp, err := time.Parse(time.RFC3339Nano, endCreatedAtStr)
		if err != nil {
			return nil, apierror.NewWarn(http.StatusBadRequest, err.Error())
		}
		endCreatedAt = &temp
	}

	if startUpdatedAtStr != "" {
		temp, err := time.Parse(time.RFC3339Nano, startUpdatedAtStr)
		if err != nil {
			return nil, apierror.NewWarn(http.StatusBadRequest, err.Error())
		}
		startUpdatedAt = &temp
	}

	if endUpdatedAtStr != "" {
		temp, err := time.Parse(time.RFC3339Nano, endUpdatedAtStr)
		if err != nil {
			return nil, apierror.NewWarn(http.StatusBadRequest, err.Error())
		}
		endUpdatedAt = &temp
	}

	res = &constants.FilterReq{
		Limit:          int64(limit),
		Page:           int64(page),
		OrderBy:        orderBy,
		Keyword:        keyword,
		SortOrder:      sortOrder,
		StartCreatedAt: startCreatedAt,
		EndCreatedAt:   endCreatedAt,
		StartUpdatedAt: startUpdatedAt,
		EndUpdatedAt:   endUpdatedAt,
	}

	err = validate.Struct(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
