package shop

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service interface {
	GetMyShop(ctx context.Context) (res *ShopRes, err error)
	GetShopByID(ctx context.Context, shopID string) (res *ShopRes, err error)
	UpdateMyShop(ctx context.Context, input UpdateShopReq, shopID string) (Toko, error)
}

type service struct {
	authConfig config.Auth
	db         *gorm.DB
}

func NewService(config *config.Config, db *gorm.DB) Service {
	return &service{
		authConfig: config.Auth,
		db:         db,
	}
}

func (s *service) GetMyShop(ctx context.Context) (res *ShopRes, err error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, apierror.FromErr(err)
	}
	userID := token.Claims.ID

	var shop Toko
	if err := s.db.WithContext(ctx).First(&shop, "id_user = ?", userID).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, you don't have a shop")
	}

	IdUser := int(shop.IdUser)
	result := &ShopRes{
		ID:       int(shop.ID),
		NamaToko: shop.NamaToko,
		UrlFoto:  shop.UrlFoto,
		IdUser:   &IdUser,
	}

	return result, nil
}

func (s *service) GetShopByID(ctx context.Context, shopID string) (res *ShopRes, err error) {

	var shop Toko
	if err := s.db.WithContext(ctx).First(&shop, "id = ?", shopID).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, shop not found")
	}

	result := &ShopRes{
		ID:       int(shop.ID),
		NamaToko: shop.NamaToko,
		UrlFoto:  shop.UrlFoto,
	}

	return result, nil
}

func (s *service) UpdateMyShop(ctx context.Context, input UpdateShopReq, shopID string) (Toko, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return Toko{}, apierror.FromErr(err)
	}
	userID := token.Claims.ID
	isAdmin := token.Claims.IsAdmin

	var shop Toko
	if !isAdmin {
		if err := s.db.WithContext(ctx).First(&shop, "id = ? AND id_user = ?", shopID, userID).Error; err != nil {
			return Toko{}, apierror.NewWarn(http.StatusNotFound, "Failed, shop not found / this is not your shop")
		}
	} else {
		if err := s.db.WithContext(ctx).First(&shop, "id = ?", shopID).Error; err != nil {
			return Toko{}, apierror.NewWarn(http.StatusNotFound, "Failed, shop not found")
		}
	}

	if input.NamaToko != nil {
		shop.NamaToko = *input.NamaToko
	}

	if input.UrlFoto != nil {
		if shop.UrlFoto != "" {
			oldPath := filepath.Join(".", shop.UrlFoto)
			if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
				return Toko{}, fmt.Errorf("failed to remove old photo: %w", err)
			}
		}

		filename := fmt.Sprintf("shop_%d_%d_%d", shop.ID, time.Now().Unix(), shop.IdUser)
		savePath := filepath.Join("uploads", "shops", filename)

		fiberCtx := ctx.Value(contextUtil.FiberCtxKey).(*fiber.Ctx)
		if err := fiberCtx.SaveFile(input.UrlFoto, savePath); err != nil {
			return Toko{}, apierror.FromErr(err)
		}

		shop.UrlFoto = "/" + savePath
	}

	shop.UpdatedAtDate = time.Now()

	if err := s.db.WithContext(ctx).Save(&shop).Error; err != nil {
		return Toko{}, apierror.FromErr(err)
	}

	return shop, nil
}

func (s *service) GetAllShop(ctx context.Context) (res []ShopRes, err error) {

	var shops []Toko
	if err := s.db.WithContext(ctx).Find(&shops).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	var result []ShopRes

	for _, cat := range shops {
		IdUser := int(cat.IdUser)
		result = append(result, ShopRes{
			ID:       int(cat.ID),
			NamaToko: cat.NamaToko,
			UrlFoto:  cat.UrlFoto,
			IdUser:   &IdUser,
		})
	}

	return result, nil
}
