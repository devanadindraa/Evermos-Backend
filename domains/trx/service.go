package trx

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/devanadindraa/Evermos-Backend/domains/address"
	"github.com/devanadindraa/Evermos-Backend/domains/category"
	"github.com/devanadindraa/Evermos-Backend/domains/product"
	"github.com/devanadindraa/Evermos-Backend/domains/shop"
	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	"github.com/devanadindraa/Evermos-Backend/utils/constants"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"gorm.io/gorm"
)

type Service interface {
	AddTrx(ctx context.Context, input TrxReq) (res *Trx, err error)
	GetTrxByID(ctx context.Context, trxID string) (*TrxRes, error)
	GetTrx(ctx context.Context, filter *constants.FilterReq) (*PaginatedTrxRes, error)
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

func (s *service) AddTrx(ctx context.Context, input TrxReq) (*Trx, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, err
	}
	userID := uint(token.Claims.ID)

	var trx *Trx

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var address address.Address
		if err := tx.Where("id = ? AND id_user = ?", input.AlamatKirim, userID).First(&address).Error; err != nil {
			return fmt.Errorf("invalid address")
		}

		var totalHarga int

		for _, item := range input.DetailTrx {
			var produk product.Product
			if err := tx.First(&produk, item.ProdukId).Error; err != nil {
				return fmt.Errorf("ID product %d not found: %w", item.ProdukId, err)
			}

			harga, err := strconv.Atoi(produk.HargaKonsumen)
			if err != nil {
				return fmt.Errorf("ID product price %d not valid: %w", item.ProdukId, err)
			}
			totalHarga += harga * item.Kuantitas
		}

		kodeInvoice := fmt.Sprintf("INV-%d", time.Now().Unix())
		trx = &Trx{
			IdUser:           userID,
			MethodBayar:      input.MethodBayar,
			AlamatPengiriman: uint(input.AlamatKirim),
			KodeInvoice:      kodeInvoice,
			HargaTotal:       totalHarga,
			CreatedAtDate:    time.Now(),
			UpdatedAtDate:    time.Now(),
		}

		if err := tx.Create(trx).Error; err != nil {
			return fmt.Errorf("failed to save transaction: %w", err)
		}

		for _, item := range input.DetailTrx {
			var produk product.Product
			if err := tx.First(&produk, item.ProdukId).Error; err != nil {
				return fmt.Errorf("ID product %d not found during insert log: %w", item.ProdukId, err)
			}

			harga, _ := strconv.Atoi(produk.HargaKonsumen)

			logProduk := LogProduk{
				IdProduk:      produk.ID,
				NamaProduk:    produk.NamaProduk,
				Slug:          produk.Slug,
				HargaReseller: produk.HargaReseller,
				HargaKonsumen: produk.HargaKonsumen,
				Deskripsi:     produk.Deskripsi,
				IdToko:        produk.IdToko,
				IdCategory:    produk.IdCategory,
				CreatedAtDate: time.Now(),
				UpdatedAtDate: time.Now(),
			}

			if err := tx.Create(&logProduk).Error; err != nil {
				return fmt.Errorf("failed to insert product log: %w", err)
			}

			detail := DetailTrx{
				IdTrx:         trx.ID,
				IdLogProduk:   logProduk.ID,
				IdToko:        logProduk.IdToko,
				Kuantitas:     item.Kuantitas,
				HargaTotal:    harga * item.Kuantitas,
				CreatedAtDate: time.Now(),
				UpdatedAtDate: time.Now(),
			}

			if err := tx.Create(&detail).Error; err != nil {
				return fmt.Errorf("failed to insert transaction details: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, apierror.FromErr(err)
	}

	return trx, nil
}

func (s *service) GetTrxByID(ctx context.Context, trxID string) (*TrxRes, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, err
	}
	userID := uint(token.Claims.ID)
	isAdmin := token.Claims.IsAdmin

	var trx Trx
	if !isAdmin {
		if err := s.db.WithContext(ctx).
			First(&trx, "id = ? AND id_user = ?", trxID, userID).Error; err != nil {
			return nil, apierror.NewWarn(http.StatusNotFound, "Failed, trx not found")
		}
	} else {
		if err := s.db.WithContext(ctx).
			First(&trx, "id = ?", trxID).Error; err != nil {
			return nil, apierror.NewWarn(http.StatusNotFound, "Failed, trx not found")
		}
	}

	var addresss address.Address
	if err := s.db.WithContext(ctx).
		First(&addresss, "id = ?", trx.AlamatPengiriman).Error; err != nil {
		return nil, apierror.NewWarn(http.StatusNotFound, "Failed, trx not found")
	}

	var details []DetailTrx
	if err := s.db.WithContext(ctx).
		Where("id_trx = ?", trx.ID).
		Find(&details).Error; err != nil {
		return nil, fmt.Errorf("failed to get detail trx: %w", err)
	}

	var detailResList []DetailTrxRes

	for _, d := range details {
		var logProduk LogProduk
		if err := s.db.WithContext(ctx).First(&logProduk, "id = ?", d.IdLogProduk).Error; err != nil {
			continue
		}

		var shops shop.Toko
		s.db.WithContext(ctx).First(&shops, logProduk.IdToko)

		var categorys category.Category
		s.db.WithContext(ctx).First(&categorys, logProduk.IdCategory)

		var photos []product.Photo
		s.db.WithContext(ctx).Where("id_produk = ?", logProduk.IdProduk).Find(&photos)

		var photoURLs []string
		for _, p := range photos {
			photoURLs = append(photoURLs, p.Url)
		}

		productRes := &product.ProductRes{
			ID:            int(logProduk.IdProduk),
			NamaProduk:    &logProduk.NamaProduk,
			Slug:          &logProduk.Slug,
			HargaReseller: &logProduk.HargaReseller,
			HargaKonsumen: &logProduk.HargaKonsumen,
			Deskripsi:     &logProduk.Deskripsi,
			Category: &category.CategoryRes{
				ID:           int(categorys.ID),
				NamaCategory: categorys.NamaCategory,
			},
			Photos: photoURLs,
		}

		detailRes := DetailTrxRes{
			Product: *productRes,
			Toko: &shop.ShopRes{
				ID:       int(logProduk.IdToko),
				NamaToko: shops.NamaToko,
				UrlFoto:  shops.UrlFoto,
			},
			Kuantitas:  d.Kuantitas,
			HargaTotal: d.HargaTotal,
		}

		detailResList = append(detailResList, detailRes)
	}

	res := &TrxRes{
		ID:          int(trx.ID),
		HargaTotal:  trx.HargaTotal,
		KodeInvoice: trx.KodeInvoice,
		MethodBayar: trx.MethodBayar,
		AlamatKirim: &address.AddressRes{
			ID:           int(trx.AlamatPengiriman),
			JudulAlamat:  addresss.JudulAlamat,
			NamaPenerima: addresss.NamaPenerima,
			NoTelp:       addresss.NoTelp,
			DetailAlamat: addresss.DetailAlamat,
		},
		DetailTrx: detailResList,
	}

	return res, nil
}

func (s *service) GetTrx(ctx context.Context, filter *constants.FilterReq) (*PaginatedTrxRes, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, err
	}
	userID := uint(token.Claims.ID)
	isAdmin := token.Claims.IsAdmin

	var totalData int64
	var trxs []Trx
	var db *gorm.DB
	if !isAdmin {
		db = s.db.Model(&Trx{}).Where("id_user = ?", userID)
	} else {
		db = s.db.Model(&Trx{})
	}
	if err := db.Count(&totalData).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := db.
		Order(fmt.Sprintf("%s %s", filter.OrderBy, filter.SortOrder)).
		Limit(int(filter.Limit)).
		Offset(int(offset)).
		Find(&trxs).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	var trxResponses []TrxRes

	for _, trx := range trxs {
		var addresss address.Address
		if err := s.db.WithContext(ctx).
			First(&addresss, "id = ?", trx.AlamatPengiriman).Error; err != nil {
			continue
		}

		var details []DetailTrx
		if err := s.db.WithContext(ctx).
			Where("id_trx = ?", trx.ID).
			Find(&details).Error; err != nil {
			continue
		}

		var detailResList []DetailTrxRes

		for _, d := range details {
			var logProduk LogProduk
			if err := s.db.WithContext(ctx).First(&logProduk, "id = ?", d.IdLogProduk).Error; err != nil {
				continue
			}

			var shops shop.Toko
			s.db.WithContext(ctx).First(&shops, logProduk.IdToko)

			var categorys category.Category
			s.db.WithContext(ctx).First(&categorys, logProduk.IdCategory)

			var photos []product.Photo
			s.db.WithContext(ctx).Where("id_produk = ?", logProduk.IdProduk).Find(&photos)

			var photoURLs []string
			for _, p := range photos {
				photoURLs = append(photoURLs, p.Url)
			}

			productRes := &product.ProductRes{
				ID:            int(logProduk.IdProduk),
				NamaProduk:    &logProduk.NamaProduk,
				Slug:          &logProduk.Slug,
				HargaReseller: &logProduk.HargaReseller,
				HargaKonsumen: &logProduk.HargaKonsumen,
				Deskripsi:     &logProduk.Deskripsi,
				Category: &category.CategoryRes{
					ID:           int(categorys.ID),
					NamaCategory: categorys.NamaCategory,
				},
				Photos: photoURLs,
			}

			detailRes := DetailTrxRes{
				Product: *productRes,
				Toko: &shop.ShopRes{
					ID:       int(logProduk.IdToko),
					NamaToko: shops.NamaToko,
					UrlFoto:  shops.UrlFoto,
				},
				Kuantitas:  d.Kuantitas,
				HargaTotal: d.HargaTotal,
			}

			detailResList = append(detailResList, detailRes)
		}

		trxResponses = append(trxResponses, TrxRes{
			ID:          int(trx.ID),
			HargaTotal:  trx.HargaTotal,
			KodeInvoice: trx.KodeInvoice,
			MethodBayar: trx.MethodBayar,
			AlamatKirim: &address.AddressRes{
				ID:           int(trx.AlamatPengiriman),
				JudulAlamat:  addresss.JudulAlamat,
				NamaPenerima: addresss.NamaPenerima,
				NoTelp:       addresss.NoTelp,
				DetailAlamat: addresss.DetailAlamat,
			},
			DetailTrx: detailResList,
		})
	}

	return &PaginatedTrxRes{
		Data:  trxResponses,
		Page:  int(filter.Page),
		Limit: int(filter.Limit),
	}, nil
}
