package category

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"gorm.io/gorm"
)

type Service interface {
	AddCategory(ctx context.Context, input CategoryReq) (res *Category, err error)
	GetAllCategory(ctx context.Context) (res []CategoryRes, err error)
	GetCategoryByID(ctx context.Context, categoryID string) (res *CategoryRes, err error)
	DeleteCategory(ctx context.Context, categoryID string) error
	UpdateCategory(ctx context.Context, input CategoryReq, categoryID string) (res *CategoryRes, err error)
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

func (s *service) AddCategory(ctx context.Context, input CategoryReq) (res *Category, err error) {
	category := Category{
		NamaCategory:  input.NamaCategory,
		CreatedAtDate: time.Now(),
		UpdatedAtDate: time.Now(),
	}

	// Insert into DB
	if err := s.db.WithContext(ctx).Create(&category).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	return &category, nil
}

func (s *service) GetAllCategory(ctx context.Context) (res []CategoryRes, err error) {

	var categories []Category
	if err := s.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	var result []CategoryRes

	for _, cat := range categories {
		result = append(result, CategoryRes{
			ID:           int(cat.ID),
			NamaCategory: cat.NamaCategory,
		})
	}

	return result, nil
}

func (s *service) GetCategoryByID(ctx context.Context, categoryID string) (res *CategoryRes, err error) {
	var category Category

	if err := s.db.WithContext(ctx).Where("id = ?", categoryID).First(&category).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	result := &CategoryRes{
		ID:           int(category.ID),
		NamaCategory: category.NamaCategory,
	}

	return result, nil
}

func (s *service) DeleteCategory(ctx context.Context, categoryID string) error {
	var category Category

	if err := s.db.WithContext(ctx).First(&category, "id = ?", categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.NewWarn(http.StatusNotFound, "Category not found")
		}
		return apierror.FromErr(err)
	}
	if err := s.db.WithContext(ctx).Where("id = ?", categoryID).Delete(&Category{}).Error; err != nil {
		return fmt.Errorf("error deleting product capital details: %v", err)
	}

	return nil
}

func (s *service) UpdateCategory(ctx context.Context, input CategoryReq, categoryID string) (res *CategoryRes, err error) {

	var category Category
	if err := s.db.WithContext(ctx).First(&category, "id = ?", categoryID).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	category.NamaCategory = input.NamaCategory
	category.UpdatedAtDate = time.Now()

	if err := s.db.WithContext(ctx).Save(&category).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	result := &CategoryRes{
		ID:           int(category.ID),
		NamaCategory: input.NamaCategory,
	}

	return result, nil
}
