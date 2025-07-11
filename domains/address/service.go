package address

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/devanadindraa/Evermos-Backend/domains/user"
	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/devanadindraa/Evermos-Backend/utils/constants"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"gorm.io/gorm"
)

type Service interface {
	AddAddress(ctx context.Context, input AddressReq) (res *Address, err error)
	GetMyAddress(ctx context.Context, filter *constants.FilterReq) ([]AddressRes, error)
	GetAddressByID(ctx context.Context, addressID string) (res *AddressRes, err error)
	DeleteAddress(ctx context.Context, addressID string) error
	UpdateAddress(ctx context.Context, input UpdateAddressReq, addressID string) (Address, error)
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

func (s *service) GetMyAddress(ctx context.Context, filter *constants.FilterReq) ([]AddressRes, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, apierror.FromErr(err)
	}
	userID := token.Claims.ID

	var addresses []Address

	query := s.db.WithContext(ctx).Where("id_user = ?", userID)

	if filter.Keyword != "" {
		query = query.Where("judul_alamat LIKE ?", "%"+filter.Keyword+"%")
	}

	if filter.OrderBy != "" && filter.SortOrder != "" {
		query = query.Order(fmt.Sprintf("%s %s", filter.OrderBy, filter.SortOrder))
	}

	offset := (filter.Page - 1) * filter.Limit
	query = query.Limit(int(filter.Limit)).Offset(int(offset))

	if err := query.Find(&addresses).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, you don't have a address")
	}

	var result []AddressRes
	for _, addr := range addresses {
		result = append(result, AddressRes{
			ID:           int(addr.ID),
			JudulAlamat:  addr.JudulAlamat,
			NamaPenerima: addr.NamaPenerima,
			NoTelp:       addr.NoTelp,
			DetailAlamat: addr.DetailAlamat,
		})
	}

	return result, nil
}

func (s *service) GetAddressByID(ctx context.Context, addressID string) (res *AddressRes, err error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return &AddressRes{}, apierror.FromErr(err)
	}

	userID := token.Claims.ID
	isAdmin := token.Claims.IsAdmin

	var address Address
	if !isAdmin {
		if err := s.db.WithContext(ctx).First(&address, "id = ?", addressID).Error; err != nil {
			return nil, apierror.NewWarn(http.StatusNotFound, "Failed, address not found")
		}
		if address.IdUser != uint(userID) {
			return nil, apierror.NewWarn(http.StatusNotFound, "Failed, this address is not yours")
		}
	} else {
		if err := s.db.WithContext(ctx).First(&address, "id = ?", addressID).Error; err != nil {
			return nil, apierror.NewWarn(http.StatusNotFound, "Failed, address not found")
		}
	}

	result := &AddressRes{
		ID:           int(address.ID),
		JudulAlamat:  address.JudulAlamat,
		NamaPenerima: address.NamaPenerima,
		NoTelp:       address.NoTelp,
		DetailAlamat: address.DetailAlamat,
	}

	return result, nil
}

func (s *service) DeleteAddress(ctx context.Context, addressID string) error {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return apierror.FromErr(err)
	}

	userID := token.Claims.ID
	isAdmin := token.Claims.IsAdmin

	var address Address

	if !isAdmin {
		if err := s.db.WithContext(ctx).First(&address, "id = ?", addressID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apierror.NewWarn(http.StatusNotFound, "Address not found")
			}
			return apierror.FromErr(err)
		}
		if address.IdUser != uint(userID) {
			return apierror.NewWarn(http.StatusNotFound, "This address is not yours")
		}
	} else {
		if err := s.db.WithContext(ctx).First(&address, "id = ?", addressID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apierror.NewWarn(http.StatusNotFound, "Address not found")
			}
			return apierror.FromErr(err)
		}
	}
	if err := s.db.WithContext(ctx).Where("id = ?", addressID).Delete(&Address{}).Error; err != nil {
		return fmt.Errorf("error deleting address details: %v", err)
	}

	return nil
}

func (s *service) UpdateAddress(ctx context.Context, input UpdateAddressReq, addressID string) (Address, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return Address{}, apierror.FromErr(err)
	}
	userID := token.Claims.ID
	isAdmin := token.Claims.IsAdmin

	var address Address
	if !isAdmin {
		if err := s.db.WithContext(ctx).First(&address, "id = ?", addressID).Error; err != nil {
			return Address{}, apierror.NewWarn(http.StatusNotFound, "Failed, shop not found")
		}
		if address.IdUser != uint(userID) {
			return Address{}, apierror.NewWarn(http.StatusNotFound, "Failed, this is not your shop")
		}
	} else {
		if err := s.db.WithContext(ctx).First(&address, "id = ?", addressID).Error; err != nil {
			return Address{}, apierror.NewWarn(http.StatusNotFound, "Failed, shop not found")
		}
	}

	if input.NoTelp != nil {
		address.NoTelp = *input.NoTelp
	}

	if input.JudulAlamat != nil {
		address.JudulAlamat = *input.JudulAlamat
	}

	if input.NamaPenerima != nil {
		address.NamaPenerima = *input.NamaPenerima
	}

	if input.DetailAlamat != nil {
		address.DetailAlamat = *input.DetailAlamat
	}

	address.UpdatedAtDate = time.Now()

	if err := s.db.WithContext(ctx).Save(&address).Error; err != nil {
		return Address{}, apierror.FromErr(err)
	}

	return address, nil
}
