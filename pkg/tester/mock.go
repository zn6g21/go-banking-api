package tester

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-banking-api/pkg/logger"
)

func MockDB() (mock sqlmock.Sqlmock, mockGormDB *gorm.DB) {
	mockDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		logger.Fatal(err.Error())
	}

	mockGormDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "mock_db",
		DriverName:                "mysql",
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	return mock, mockGormDB
}
