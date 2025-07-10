package user

import (
	"context"
	"net/http"
	"time"

	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/devanadindraa/Evermos-Backend/utils/constants"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Service interface {
	Login(ctx context.Context, input LoginReq) (res *LoginRes, err error)
	Logout(ctx context.Context, input LogoutReq) (res *LogoutRes, err error)
	ValidateToken(ctx context.Context, token string) (err error)
	Register(ctx context.Context, input RegisterReq) (res *User, err error)
	UpdateProfile(ctx context.Context, input UpdateProfileReq) (res *User, err error)
	GetProfile(ctx context.Context) (*User, error)
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

func (s *service) Login(ctx context.Context, input LoginReq) (*LoginRes, error) {

	var err error
	var user User
	if err = s.db.WithContext(ctx).Where("notelp = ?", input.Notelp).First(&user).Error; err == nil {
		if !comparePassword(user.KataSandi, input.KataSandi) {
			return nil, apierror.NewWarn(http.StatusUnauthorized, ErrInvalidCredentials)
		}
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NewWarn(http.StatusUnauthorized, ErrInvalidCredentials)
		}
		return nil, err
	}

	expirationTime := time.Now().Add(s.authConfig.JWT.ExpireIn)
	claims := &constants.JWTClaims{
		ID:     int(user.ID),
		NoTelp: input.Notelp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.authConfig.JWT.SecretKey))
	if err != nil {
		return nil, apierror.FromErr(err)
	}

	return &LoginRes{
		Token:   tokenString,
		Expires: expirationTime,
	}, nil
}

func (s *service) Logout(ctx context.Context, input LogoutReq) (res *LogoutRes, err error) {

	invalidToken := InvalidToken(input)

	err = s.db.WithContext(ctx).Create(&invalidToken).Error

	if err != nil {
		return nil, err
	}

	return &LogoutRes{
		LoggedOut: true,
	}, nil
}

func (s *service) ValidateToken(ctx context.Context, token string) (err error) {

	var invalidToken InvalidToken
	if err := s.db.WithContext(ctx).Where("token = ?", token).First(&invalidToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	return nil
}

func (s *service) Register(ctx context.Context, input RegisterReq) (res *User, err error) {

	// Hash password
	hashedPassword, err := hashPassword(input.KataSandi)
	if err != nil {
		return nil, apierror.FromErr(err)
	}

	parsedDate, err := ParseDateFromPointer(input.TanggalLahir, "02/01/2006")
	if err != nil {
		return nil, apierror.FromErr(err)
	}

	// Build user object
	user := User{
		Nama:          input.Nama,
		KataSandi:     hashedPassword,
		Notelp:        input.NoTelp,
		TanggalLahir:  parsedDate,
		Pekerjaan:     GetStringOrDefault(input.Pekerjaan, ""),
		Email:         input.Email,
		IdProvinsi:    GetStringOrDefault(input.IdProvinsi, ""),
		IdKota:        GetStringOrDefault(input.IdKota, ""),
		IsAdmin:       GetBoolOrDefault(input.IsAdmin, false),
		CreatedAtDate: time.Now(),
		UpdatedAtDate: time.Now(),
	}

	// Insert into DB
	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	return &user, nil
}

func (s *service) UpdateProfile(ctx context.Context, input UpdateProfileReq) (res *User, err error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, apierror.FromErr(err)
	}
	userID := token.Claims.ID

	var user User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	if input.Nama != nil {
		user.Nama = *input.Nama
	}
	if input.KataSandi != nil {
		hashedPassword, err := hashPassword(*input.KataSandi)
		if err != nil {
			return nil, apierror.FromErr(err)
		}
		user.KataSandi = hashedPassword
	}
	if input.NoTelp != nil {
		user.Notelp = *input.NoTelp
	}
	if input.TanggalLahir != nil {
		parsedDate, err := ParseDateFromPointer(input.TanggalLahir, "02/01/2006")
		if err != nil {
			return nil, apierror.FromErr(err)
		}
		user.TanggalLahir = parsedDate
	}
	if input.Pekerjaan != nil {
		user.Pekerjaan = *input.Pekerjaan
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.IdProvinsi != nil {
		user.IdProvinsi = *input.IdProvinsi
	}
	if input.IdKota != nil {
		user.IdKota = *input.IdKota
	}
	if input.IsAdmin != nil {
		user.IsAdmin = *input.IsAdmin
	}

	user.UpdatedAtDate = time.Now()

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	return &user, nil
}

func (s *service) GetProfile(ctx context.Context) (*User, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, apierror.FromErr(err)
	}

	userID := token.Claims.ID

	var user User
	if err := s.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	user.KataSandi = ""

	return &user, nil
}
