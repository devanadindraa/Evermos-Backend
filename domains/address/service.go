package address

import (
	"context"
	"net/http"
	"time"

	"github.com/devanadindraa/Evermos-Backend/domains/user"
	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"gorm.io/gorm"
)

type Service interface {
	AddAddress(ctx context.Context, input AddressReq) (res *Address, err error)
	GetMyAddress(ctx context.Context) (res []AddressRes, err error)
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

func (s *service) AddAddress(ctx context.Context, input AddressReq) (res *Address, err error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return &Address{}, apierror.FromErr(err)
	}

	userID := token.Claims.ID

	var user user.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	address := Address{
		IdUser:       uint(userID),
		JudulAlamat:  input.JudulAlamat,
		NamaPenerima: input.NamaPenerima,
		NoTelp: func() string {
			if input.NoTelp != nil {
				return *input.NoTelp
			}
			return user.Notelp
		}(),
		DetailAlamat:  input.DetailAlamat,
		CreatedAtDate: time.Now(),
		UpdatedAtDate: time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(&address).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	return &address, nil
}

func (s *service) GetMyAddress(ctx context.Context) (res []AddressRes, err error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, apierror.FromErr(err)
	}
	userID := token.Claims.ID

	var address []Address
	if err := s.db.WithContext(ctx).Find(&address, "id_user = ?", userID).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, you don't have a address")
	}

	var result []AddressRes

	for _, cat := range address {
		result = append(result, AddressRes{
			ID:           int(cat.ID),
			JudulAlamat:  cat.JudulAlamat,
			NamaPenerima: cat.NamaPenerima,
			NoTelp:       cat.NoTelp,
			DetailAlamat: cat.DetailAlamat,
		})
	}

	return result, nil
}
