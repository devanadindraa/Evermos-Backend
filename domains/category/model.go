package category

import "time"

type Category struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	NamaCategory  string    `json:"nama_category"`
	CreatedAtDate time.Time `gorm:"autoCreateTime"`
	UpdatedAtDate time.Time `gorm:"autoUpdateTime"`
}

func (Category) TableName() string {
	return "category"
}
