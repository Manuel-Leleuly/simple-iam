package initializers

import (
	"os"
	"time"

	"github.com/Manuel-Leleuly/simple-iam/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func getGORMDatabaseUrl(params map[string]string) string {
	var dbUrl = os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_URL") + ")/" + os.Getenv("DB_NAME")
	queryParams := ""

	for k, v := range params {
		selectedParam := k + "=" + v
		if queryParams == "" {
			queryParams = "?" + selectedParam
		} else {
			queryParams += "&" + selectedParam
		}
	}

	return dbUrl + queryParams
}

func ConnectToDB() {
	dialect := mysql.Open(getGORMDatabaseUrl(map[string]string{
		"charset":   "utf8mb4",
		"parseTime": "True",
		"loc":       "Local",
	}))

	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	DB = db
}

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
