package product

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	UpdateProduct(ctx context.Context, input UpdateProductReq, IdToko string) (res *ProductRes, err error)
	GetProducts(ctx context.Context, filter GetProductReq) ([]ProductRes, error)
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

func (s *service) UpdateProduct(ctx context.Context, input UpdateProductReq, productId string) (res *ProductRes, err error) {
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
	if err := tx.First(&shop, "id_user = ?", UserID).Error; err != nil {
		tx.Rollback()
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, shop not found")
	}

	var product Product
	if err := tx.First(&product, "id_toko = ? AND id = ?", shop.ID, productId).Error; err != nil {
		tx.Rollback()
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, product not found")
	}

	if input.NamaProduk != nil {
		product.NamaProduk = *input.NamaProduk
	}
	if input.Slug != nil {
		product.Slug = *input.Slug
	}
	if input.IdCategory != nil {
		product.IdCategory = uint(*input.IdCategory)
	}
	if input.HargaReseller != nil {
		product.HargaReseller = *input.HargaReseller
	}
	if input.HargaKonsumen != nil {
		product.HargaKonsumen = *input.HargaKonsumen
	}
	if input.Stok != nil {
		product.Stok = *input.Stok
	}
	if input.Deskripsi != nil {
		product.Deskripsi = *input.Deskripsi
	}
	product.UpdatedAtDate = time.Now()

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var photos []Photo
	if input.Photos != nil {
		var oldPhotos []Photo
		if err := tx.Where("id_produk = ?", product.ID).Find(&oldPhotos).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, p := range oldPhotos {
			fullPath := filepath.Clean("." + p.Url)
			if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
				tx.Rollback()
				return nil, fmt.Errorf("failed to delete old photo: %v", err)
			}
		}

		if err := tx.Where("id_produk = ?", product.ID).Delete(&Photo{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		for _, file := range *input.Photos {
			ext := filepath.Ext(file.Filename)
			filename, err := fileutils.GenerateMediaName(strconv.Itoa(int(product.ID)))
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("error generating image name: %v", err)
			}
			filename = fmt.Sprintf("%s%s", filename, ext)
			path := filepath.Join("uploads", "products", filename)

			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create directory: %v", err)
			}

			if err := fileutils.SaveMedia(ctx, file, path); err != nil {
				tx.Rollback()
				return nil, err
			}

			photos = append(photos, Photo{
				IdProduk: product.ID,
				Url:      "/uploads/products/" + filename,
			})
		}

		if len(photos) > 0 {
			if err := tx.Create(&photos).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &ProductRes{
		ID:         int(product.ID),
		NamaProduk: &product.NamaProduk,
	}, nil
}

func (s *service) GetProducts(ctx context.Context, filter GetProductReq) ([]ProductRes, error) {
	var products []Product

	db := s.db.WithContext(ctx).Model(&Product{})

	if filter.Keyword != "" {
		db = db.Where("nama_produk LIKE ?", "%"+filter.Keyword+"%")
	}

	if filter.CategoryID != nil {
		db = db.Where("id_category = ?", *filter.CategoryID)
	}

	if filter.TokoID != nil {
		db = db.Where("id_toko = ?", *filter.TokoID)
	}

	if filter.MinHarga != nil {
		db = db.Where("CAST(harga_konsumen AS UNSIGNED) >= ?", *filter.MinHarga)
	}

	if filter.MaxHarga != nil {
		db = db.Where("CAST(harga_konsumen AS UNSIGNED) <= ?", *filter.MaxHarga)
	}

	if filter.StartCreatedAt != nil {
		db = db.Where("created_at_date >= ?", *filter.StartCreatedAt)
	}
	if filter.EndCreatedAt != nil {
		db = db.Where("created_at_date <= ?", *filter.EndCreatedAt)
	}
	if filter.StartUpdatedAt != nil {
		db = db.Where("updated_at_date >= ?", *filter.StartUpdatedAt)
	}
	if filter.EndUpdatedAt != nil {
		db = db.Where("updated_at_date <= ?", *filter.EndUpdatedAt)
	}

	orderStr := fmt.Sprintf("%s %s", filter.OrderBy, strings.ToUpper(filter.SortOrder))
	db = db.Order(orderStr)

	offset := (filter.Page - 1) * filter.Limit
	db = db.Limit(int(filter.Limit)).Offset(int(offset))

	db = db.Preload("Photos")

	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}

	var result []ProductRes
	for _, p := range products {
		p := p
		res := ProductRes{
			ID:            int(p.ID),
			NamaProduk:    &p.NamaProduk,
			Slug:          &p.Slug,
			IdCategory:    &p.IdCategory,
			HargaReseller: &p.HargaReseller,
			HargaKonsumen: &p.HargaKonsumen,
			Stok:          &p.Stok,
			Deskripsi:     &p.Deskripsi,
		}
		result = append(result, res)
	}

	return result, nil
}
