package trx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/devanadindraa/Evermos-Backend/domains/address"
	"github.com/devanadindraa/Evermos-Backend/domains/product"
	apierror "github.com/devanadindraa/Evermos-Backend/utils/api-error"
	"github.com/devanadindraa/Evermos-Backend/utils/config"
	contextUtil "github.com/devanadindraa/Evermos-Backend/utils/context"
	"gorm.io/gorm"
)

type Service interface {
	AddTrx(ctx context.Context, input TrxReq) (res *Trx, err error)
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
