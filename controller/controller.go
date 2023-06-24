package controller

import (
	"math/rand"
	"net/http"
	"strconv"

	"assesment-test/config"
	"assesment-test/logger"
	"assesment-test/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	Db *gorm.DB
}

func New() *Server {
	db := database.InitDB()
	db.AutoMigrate(&model.Nasabah{})
	db.AutoMigrate(&model.Mutasi{})
	db.AutoMigrate(&model.Saldo{})
	return &Server{Db: db}
}

func (repository *Server) CreateNasabah(c echo.Context) error {
	no_rek := rand.Uint32()
	nasabah := model.Nasabah{No_rekening: no_rek}
	c.Bind(&nasabah)
	
	// Cek Nik dan No HP
	res, err := model.ValidateNasabah(repository.Db, &nasabah)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if res.Nik != "" || res.No_hp != "" {
		logger.LoggerWarn("Nik atau No HP Sudah Ada")
		return c.JSON(http.StatusBadRequest, map[string]string{"remark": "Data Sudah Ada"})
	}
	
	// Create Nasabah
	err = model.CreateNasabah(repository.Db, &nasabah)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	logger.LoggerInfo("Create Nasabah Berhasil")

	// Create Saldo
	saldo := model.Saldo{No_Rekening: no_rek}
	err = model.CreateSaldo(repository.Db, &saldo)
	if err != nil {
		return err
	}
	logger.LoggerInfo("Create Saldo Berhasil")

	return c.JSON(http.StatusOK, nasabah)
}

func (repository *Server) CreateTabung(c echo.Context) error {
	mutasi := model.Mutasi{Kode_transaksi: "C"}
	c.Bind(&mutasi)

	// Cek No rekening
	nasabah, err := model.ValidateNoNasabah(repository.Db, mutasi.No_rekening)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if nasabah.No_rekening == 0 {
		logger.LoggerWarn("Rekening Tidak Ada")
		return c.JSON(http.StatusBadRequest, map[string]string{"remark": "Rekening Tidak Ada"})
	}

	// Cek Saldo
	saldo, err := model.CekSaldoForUpdate(repository.Db, mutasi.No_rekening)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Create Mutasi
	err = model.CreateMutasi(repository.Db, &mutasi)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	logger.LoggerInfo("Create Mutasi Berhasil")

	// Hitung Saldo
	logger.LoggerInfo("Proses Hitung Saldo")
	res, err := model.SumMutasi(repository.Db, &mutasi)
	saldo1 := res.Total_Credit - res.Total_Debit
	saldo.Saldo = saldo1

	// Update Saldo
	err = model.UpdateSaldo(repository.Db, &saldo)
	if err != nil {
		return err
	}
	logger.LoggerInfo("Update Saldo Berhasil")

	return c.JSON(http.StatusOK, map[string]float64{"Saldo": saldo.Saldo})
}

func (repository *Server) CreateTarik(c echo.Context) error {
	mutasi := model.Mutasi{Kode_transaksi: "D"}
	c.Bind(&mutasi)

	// Cek No rekening
	nasabah, err := model.ValidateNoNasabah(repository.Db, mutasi.No_rekening)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if nasabah.No_rekening == 0 {
		logger.LoggerWarn("Rekening Tidak Ada")
		return c.JSON(http.StatusBadRequest, map[string]string{"remark": "Rekening Tidak Ada"})
	}

	// Cek Saldo
	saldo, err := model.CekSaldoForUpdate(repository.Db, mutasi.No_rekening)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	nominal := mutasi.Nominal
	if saldo.Saldo-nominal < 0 {
		logger.LoggerWarn("Saldo Tidak Cukup")
		return c.JSON(http.StatusBadRequest, map[string]string{"remark": "Saldo Tidak Cukup"})
	}

	// Create Mutasi
	err = model.CreateMutasi(repository.Db, &mutasi)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	logger.LoggerInfo("Create Mutasi Success")

	// Hitung Saldo
	res, err := model.SumMutasi(repository.Db, &mutasi)
	saldo1 := res.Total_Credit - res.Total_Debit
	saldo.Saldo = saldo1

	// Update Saldo
	err = model.UpdateSaldo(repository.Db, &saldo)
	if err != nil {
		return err
	}
	logger.LoggerInfo("Update Saldo Success")

	return c.JSON(http.StatusOK, map[string]float64{"Saldo": saldo.Saldo})
}

func (repository *Server) CekSaldo(c echo.Context) error {
	norek := c.Param("no_rekening")
	u, err := strconv.ParseUint(norek, 0, 32)
	saldo := model.Saldo{}
	c.Bind(&saldo)

	// Cek No rekening
	nasabah, err := model.ValidateNoNasabah(repository.Db, uint32(u))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if nasabah.No_rekening == 0 {
		logger.LoggerWarn("Rekening Tidak Ada")
		return c.JSON(http.StatusBadRequest, map[string]string{"remark": "Rekening Tidak Ada"})
	}

	// Cek Saldo
	res, err := model.CekSaldo(repository.Db, &saldo, uint32(u))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, map[string]float64{"Saldo": res.Saldo})
}

func (repository *Server) CekMutasi(c echo.Context) error {
	norek := c.Param("no_rekening")
	u, err := strconv.ParseUint(norek, 0, 32)
	mutasi := model.Mutasi{}
	c.Bind(&mutasi)

	// Cek No rekening
	nasabah, err := model.ValidateNoNasabah(repository.Db, uint32(u))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if nasabah.No_rekening == 0 {
		logger.LoggerWarn("Rekening Tidak Ada")
		return c.JSON(http.StatusBadRequest, map[string]string{"remark": "Rekening Tidak Ada"})
	}

	// Cek Mutasi
	res, err := model.CekMutasi(repository.Db, &mutasi, uint32(u))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, res)
}
