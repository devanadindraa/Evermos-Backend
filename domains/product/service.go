package product

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/devanadindraa/Evermos-Backend/domains/shop"
	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	fileutils "github.com/devanadindraa/Evermos-Backend/utils/file"
	"gorm.io/gorm"
)

type Service interface {
	AddProduct(ctx context.Context, input ProductReq) (res *ProductRes, err error)
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

func (s *service) AddProduct(ctx context.Context, input ProductReq) (res *ProductRes, err error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, err
	}

	UserID := token.Claims.ID

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var shop shop.Toko
	if err := s.db.WithContext(ctx).First(&shop, "id_user = ?", UserID).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, shop not found")
	}

	product := Product{
		IdToko:        shop.ID,
		NamaProduk:    input.NamaProduk,
		IdCategory:    input.IdCategory,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Stok:          input.Stok,
		Deskripsi:     input.Deskripsi,
		CreatedAtDate: time.Now(),
		UpdatedAtDate: time.Now(),
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var photos []Photo
	for _, file := range input.Photos {
		ext := filepath.Ext(file.Filename)

		filename, err := fileutils.GenerateMediaName(strconv.Itoa(int(product.ID)))
		if err != nil {
			tx.Rollback()
			return nil, apierror.NewWarn(http.StatusNotFound, fmt.Sprintf("error generating image name: %v", err))
		}

		filename = fmt.Sprintf("%s%s", filename, ext)
		path := filepath.Join("uploads", "products", filename)

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			tx.Rollback()
			return nil, apierror.NewWarn(http.StatusNotFound, fmt.Sprintf("failed to create directory: %v", err))
		}

		if err := fileutils.SaveMedia(ctx, file, path); err != nil {
			tx.Rollback()
			return nil, err
		}

		photos = append(photos, Photo{
			IdProduk:      product.ID,
			Url:           "/uploads/products/" + filename,
			CreatedAtDate: time.Now(),
			UpdatedAtDate: time.Now(),
		})
	}

	if len(photos) > 0 {
		if err := tx.Create(&photos).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	result := &ProductRes{
		ID:         int(product.ID),
		NamaProduk: product.NamaProduk,
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}
