package model

import (
	"time"

	"gorm.io/gorm"
)

type Mutasi struct {
	Id             int       `gorm:"primaryKey" json:"id"`
	No_rekening    uint32    `gorm:"column:no_rekening" json:"no_rekening"`
	Nominal        float64   `gorm:"column:nominal" json:"nominal"`
	Kode_transaksi string    `gorm:"column:kode_transaksi" json:"kode_transaksi"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type CekSaldoMutasi struct {
	Total_Debit  float64
	Total_Credit float64
}

func CreateMutasi(db *gorm.DB, mts *Mutasi) (err error) {
	err = db.Create(mts).Error
	if err != nil {
		return err
	}
	return nil
}

func CekMutasi(db *gorm.DB, mts *Mutasi, norek uint32) (_ []Mutasi, err error) {
	var mutasi []Mutasi
	err = db.Where("no_rekening = ?", norek).Find(&mutasi).Error
	if err != nil {
		return mutasi, err
	}
	return mutasi, nil
}

func SumMutasi(db *gorm.DB, mts *Mutasi) (res CekSaldoMutasi, err error) {
	var ceksaldo CekSaldoMutasi
	norek := mts.No_rekening
	err = db.Model(&Mutasi{}).Select(" no_rekening,sum(case when kode_transaksi='D' then nominal else 0 end) as total_debit,sum(case when kode_transaksi='C' then nominal else 0 end) as total_credit").Where("no_rekening = ?", norek).Group("no_rekening").Find(&ceksaldo).Error
	if err != nil {
		return ceksaldo, err
	}
	return ceksaldo, nil
}
