package category

import (
	"context"
	"time"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"gorm.io/gorm"
)

type Service interface {
	AddCategory(ctx context.Context, input CategoryReq) (res *Category, err error)
	GetAllCategory(ctx context.Context) (res []CategoryRes, err error)
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
	// Build user object
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
