package models

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	DB *gorm.DB
}

func (d *DBInstance) ConnectToDB(dbName string) error {
	if dbName == "" {
		return errors.New("DB name is empty")
	}

	dialect := mysql.Open(getGORMDatabaseUrl(dbName, map[string]string{
		"charset":   "utf8mb4",
		"parseTime": "True",
		"loc":       "Local",
	}))

	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	d.DB = db

	return nil
}

func (d *DBInstance) IsDBConnected() bool {
	return d.DB != nil
}

func (d *DBInstance) SyncDatabase() error {
	if !d.IsDBConnected() {
		return errors.New("DB is not initialized")
	}

	d.DB.AutoMigrate(&User{})

	return nil
}

func (d *DBInstance) CheckDBConnection(c *gin.Context) {
	if !d.IsDBConnected() {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorMessage{
			Message: "DB is not initialized",
		})
	} else {
		c.Next()
	}

}

func (d *DBInstance) MakeHTTPHandleFunc(f ApiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if statusCode, err := f(d, c); err != nil {
			c.JSON(statusCode, ErrorMessage{
				Message: err.Error(),
			})
		}
	}
}

// helpers
func getGORMDatabaseUrl(dbName string, params map[string]string) string {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), dbName)
	queryParams := url.Values{}

	for k, v := range params {
		queryParams.Add(k, v)
	}

	return dbUrl + "?" + queryParams.Encode()
}
