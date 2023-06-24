package main

import (
	"assesment-test/api/controller"
	"assesment-test/api/config"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func makeLogEntry(c echo.Context) *log.Entry {
    if c == nil {
        return log.WithFields(log.Fields{
            "at": time.Now().Format("2006-01-02 15:04:05"),
        })
    }

    return log.WithFields(log.Fields{
        "at":     time.Now().Format("2006-01-02 15:04:05"),
        "method": c.Request().Method,
        "uri":    c.Request().URL.String(),
        "ip":     c.Request().RemoteAddr,
    })
}
func middlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        makeLogEntry(c).Info("incoming request")
        return next(c)
    }
}

func errorHandler(err error, c echo.Context) {
    report, ok := err.(*echo.HTTPError)
    if ok {
        report.Message = fmt.Sprintf("http error %d - %v", report.Code, report.Message)
    } else {
        report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    makeLogEntry(c).Error(report.Message)
    c.HTML(report.Code, report.Message.(string))
}

func main() {
	e := echo.New()
	e.Use(middlewareLogging)
	e.HTTPErrorHandler = errorHandler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	database.InitDB()
	controller := controller.New()
	e.POST("/daftar", controller.CreateNasabah)
	e.POST("/tabung", controller.CreateTabung)
	e.POST("/tarik", controller.CreateTarik)
	e.GET("/saldo/:no_rekening", controller.CekSaldo)
	e.GET("/mutasi/:no_rekening", controller.CekMutasi)
	lock := make(chan error)
    go func(lock chan error) { lock <- e.Start(":1323") }(lock)

    time.Sleep(1 * time.Millisecond)
    makeLogEntry(nil).Warning("application started ")

    err := <-lock
    if err != nil {
        makeLogEntry(nil).Panic("failed to start application")
    }
	// e.Logger.Fatal(e.Start(":1323"))
}
