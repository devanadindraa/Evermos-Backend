package middlewares

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/devanadindraa/Evermos-Backend/utils/constants"
	"github.com/gofiber/fiber/v2"
)

func getRequestPayload(fiberCtx *fiber.Ctx) *constants.RequestPayload {

	// get raw body
	body := fiberCtx.Body()
	bodyStr := toValidJson(body)

	// copy header so it doest affect original header
	headers := make(map[string][]string)
	fiberCtx.Request().Header.VisitAll(func(k, v []byte) {
		key := string(k)
		value := string(v)
		headers[key] = append(headers[key], value)
	})

	// Mask sensitive headers
	for _, ignored := range IGNORED_HEADERS {
		for k := range headers {
			if strings.EqualFold(k, ignored) {
				for i := range headers[k] {
					headers[k][i] = "***********"
				}
			}
		}
	}

	queries := url.Values{}
	for key, val := range fiberCtx.Queries() {
		queries.Set(key, val)
	}

	return &constants.RequestPayload{
		Body:        bodyStr,
		QueryParams: queries,
		Headers:     headers,
	}
}

func toValidJson(data []byte) (res map[string]any) {
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil
	}

	return res
}
