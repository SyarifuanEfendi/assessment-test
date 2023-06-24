package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Saldo struct {
	No_Rekening uint32    `gorm:"primaryKey" json:"no_rekening"`
	Saldo       float64   `gorm:"column:saldo" json:"saldo"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func CekSaldo(db *gorm.DB, sld *Saldo, norek uint32) (res Saldo, err error) {
	var saldo Saldo
	err = db.Where("no_rekening = ?", norek).Find(&saldo).Error
	if err != nil {
		return saldo, err
	}
	return saldo, nil
}

func CekSaldoForUpdate(db *gorm.DB, norek uint32) (res Saldo, err error) {
	var saldo Saldo
	err = db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("no_rekening = ?", norek).Find(&saldo).Error
	if err != nil {
		return saldo, err
	}
	return saldo, nil
}

func CreateSaldo(db *gorm.DB, sld *Saldo) (err error) {
	err = db.Create(sld).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateSaldo(db *gorm.DB, sld *Saldo) (err error) {
	err = db.Save(sld).Where("no_rekenig = ?", sld.No_Rekening).Error
	if err != nil {
		return err
	}
	return nil
}
