package middlewares

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devanadindraa/Evermos-Backend/domains/user"
	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/devanadindraa/Evermos-Backend/utils/constants"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"github.com/devanadindraa/Evermos-Backend/utils/logger"
	"github.com/devanadindraa/Evermos-Backend/utils/respond"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type Middlewares interface {
	AddRequestId(ctx *fiber.Ctx) error
	Logging(ctx *fiber.Ctx) error
	BasicAuth(ctx *fiber.Ctx) error
	JWT(ctx *fiber.Ctx) error
	Recover(ctx *fiber.Ctx) error
	RateLimiter(ctx *fiber.Ctx) error
}

type middlewares struct {
	conf        *config.Config
	rateLimiter *rate.Limiter
	userService user.Service
}

// Constructor untuk middlewares
func NewMiddlewares(conf *config.Config, userService user.Service) Middlewares {
	return &middlewares{
		conf:        conf,
		rateLimiter: rate.NewLimiter(rate.Limit(conf.RateLimiter.Rps), conf.RateLimiter.Bursts),
		userService: userService,
	}
}

func (m *middlewares) AddRequestId(ctx *fiber.Ctx) error {
	requestId := uuid.New()
	contextUtil.FiberWithCtx(ctx, contextUtil.SetRequestId(context.Background(), requestId))
	ctx.Set("Request-Id", requestId.String())
	return ctx.Next()
}

func (m *middlewares) Logging(ctx *fiber.Ctx) error {
	start := time.Now()
	reqPayload := getRequestPayload(ctx)

	ctx.Next()

	logPayload := logger.LogPayload{
		Method:         ctx.Method(),
		Path:           ctx.Path(),
		StatusCode:     ctx.Response().StatusCode(),
		Took:           time.Since(start),
		RequestPayload: reqPayload,
	}

	ctxVal := ctx.Locals("ctx")
	ctxReal, ok := ctxVal.(context.Context)
	if !ok {
		ctxReal = context.Background()
	}

	var err error
	errAny := ctx.Locals("error")
	if errAny != nil {
		if castErr, ok := errAny.(error); ok {
			err = castErr
		}
	}

	logger.Log(ctxReal, logPayload, err)
	return nil
}

func (m *middlewares) BasicAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		respond.Error(ctx, apierror.Unauthorized())
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	// Decode base64(username:password)
	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
	if err != nil {
		respond.Error(ctx, apierror.Unauthorized())
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	// Split username and password
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		respond.Error(ctx, apierror.Unauthorized())
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	username := pair[0]
	password := pair[1]

	// Hash and compare securely
	usernameHash := sha256.Sum256([]byte(username))
	passwordHash := sha256.Sum256([]byte(password))
	expectedUsernameHash := sha256.Sum256([]byte(m.conf.Auth.Basic.Username))
	expectedPasswordHash := sha256.Sum256([]byte(m.conf.Auth.Basic.Password))

	usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
	passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

	if usernameMatch && passwordMatch {
		return ctx.Next()
	}

	respond.Error(ctx, apierror.Unauthorized())
	return ctx.SendStatus(fiber.StatusUnauthorized)
}

func (m *middlewares) JWT(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		authHeader = ctx.Get("Auth")
	}

	authorizationSplit := strings.Split(authHeader, " ")
	if len(authorizationSplit) < 2 {
		respond.Error(ctx, apierror.Unauthorized())
		return nil
	}

	tokenStr := authorizationSplit[1]
	claims := constants.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.conf.Auth.JWT.SecretKey), nil
	})
	if err != nil || !token.Valid {
		respond.Error(ctx, apierror.Unauthorized())
		return nil
	}

	err = m.userService.ValidateToken(ctx.Context(), tokenStr)
	if err != nil {
		respond.Error(ctx, apierror.Unauthorized())
		return nil
	}

	ctx.Locals("token", constants.Token{
		Token:  tokenStr,
		Claims: claims,
	})

	return ctx.Next()
}

func (m *middlewares) Recover(ctx *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			respond.Error(ctx, apierror.NewError(http.StatusInternalServerError, fmt.Sprintf("Panic : %v", r)))
		}
	}()
	return ctx.Next()
}

func (m *middlewares) RateLimiter(ctx *fiber.Ctx) error {
	if !m.rateLimiter.Allow() {
		respond.Error(ctx, apierror.NewWarn(http.StatusTooManyRequests, "Too many request"))
		return ctx.SendStatus(fiber.StatusTooManyRequests)
	}
	return ctx.Next()
}
