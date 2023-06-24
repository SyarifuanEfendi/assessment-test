package model

import (
	"time"

	"gorm.io/gorm"
)

type Nasabah struct {
	No_rekening uint32    `gorm:"primaryKey" json:"no_rekening"`
	Nama        string    `gorm:"column:nama" json:"nama"`
	Nik         string    `gorm:"column:nik" json:"nik"`
	No_hp       string    `gorm:"column:no_hp" json:"no_hp"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func CreateNasabah(db *gorm.DB, Nasabah *Nasabah) (err error) {
	err = db.Create(Nasabah).Error
	if err != nil {
		return err
	}
	return nil
}

func ValidateNasabah(db *gorm.DB, nsb *Nasabah) (res Nasabah, err error) {
	var nasabah Nasabah

	nik := nsb.Nik
	no_hp := nsb.No_hp

	err = db.Where("nik = ? or no_hp = ?", nik, no_hp).Find(&nasabah).Error

	if err != nil {
		return nasabah, err
	}
	return nasabah, nil
}

func ValidateNoNasabah(db *gorm.DB, norek uint32) (res Nasabah, err error) {
	var nasabah Nasabah

	err = db.Where("no_rekening = ?", norek).Find(&nasabah).Error
	if err != nil {
		return nasabah, err
	}
	return nasabah, nil
}
