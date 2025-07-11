package product

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/devanadindraa/Evermos-Backend/domains/category"
	"github.com/devanadindraa/Evermos-Backend/domains/shop"
	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	fileutils "github.com/devanadindraa/Evermos-Backend/utils/file"
	"gorm.io/gorm"
)

type Service interface {
	AddProduct(ctx context.Context, input ProductReq) (res *ProductRes, err error)
	GetProductByID(ctx context.Context, productID string) (res *ProductRes, err error)
	DeleteProduct(ctx context.Context, productID string) error
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
		IdToko:     shop.ID,
		NamaProduk: input.NamaProduk,
		Slug: func() string {
			if input.Slug != nil {
				return *input.Slug
			}
			return ""
		}(),
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
		NamaProduk: &product.NamaProduk,
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetProductByID(ctx context.Context, productID string) (res *ProductRes, err error) {

	var product Product
	if err := s.db.WithContext(ctx).Preload("Photos").First(&product, "id = ?", productID).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, product not found")
	}

	var shops shop.Toko
	if err := s.db.WithContext(ctx).First(&shops, "id = ?", product.IdToko).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, shop not found")
	}

	var categorys category.Category
	if err := s.db.WithContext(ctx).First(&categorys, "id = ?", product.IdCategory).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, category not found")
	}

	var photoURLs []string
	for _, p := range product.Photos {
		photoURLs = append(photoURLs, p.Url)
	}

	result := &ProductRes{
		ID:            int(product.ID),
		NamaProduk:    &product.NamaProduk,
		Slug:          &product.Slug,
		HargaReseller: &product.HargaReseller,
		HargaKonsumen: &product.HargaKonsumen,
		Stok:          &product.Stok,
		Deskripsi:     &product.Deskripsi,
		Shop: &shop.ShopRes{
			ID:       int(shops.ID),
			NamaToko: shops.NamaToko,
			UrlFoto:  shops.UrlFoto,
		},
		Category: &category.CategoryRes{
			ID:           int(categorys.ID),
			NamaCategory: categorys.NamaCategory,
		},
		Photos: photoURLs,
	}

	return result, nil
}

func (s *service) DeleteProduct(ctx context.Context, productID string) error {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return apierror.FromErr(err)
	}

	userID := token.Claims.ID
	isAdmin := token.Claims.IsAdmin

	var shop shop.Toko
	if !isAdmin {
		if err := s.db.WithContext(ctx).First(&shop, "id_user = ?", userID).Error; err != nil {
			return apierror.FromErr(err)
		}
	}

	var product Product
	if err := s.db.WithContext(ctx).
		Preload("Photos").
		First(&product, "id = ?", productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.NewWarn(http.StatusNotFound, "Product not found")
		}
		return apierror.FromErr(err)
	}

	if !isAdmin && product.IdToko != shop.ID {
		return apierror.NewWarn(http.StatusForbidden, "This product is not yours")
	}

	for _, photo := range product.Photos {
		_ = os.Remove("." + photo.Url)
	}

	if err := s.db.WithContext(ctx).Delete(&product).Error; err != nil {
		return fmt.Errorf("error deleting product details: %v", err)
	}

	return nil
}
