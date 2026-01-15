package tester

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"go-banking-api/entity"
	"go-banking-api/infrastructure/database"
)

type DBSQLiteSuite struct {
	suite.Suite
	DB     *gorm.DB
	DBName string
}

func (suite *DBSQLiteSuite) SetupSuite() {
	suite.DBName = fmt.Sprintf("%s.unittest.sqlite", suite.T().Name())

	err := os.Setenv("DB_NAME", suite.DBName)
	suite.Assert().Nil(err)
	db, err := database.NewDatabaseSQLFactory(database.InstanceSQLite)
	suite.Assert().Nil(err)
	suite.DB = db

	for _, model := range entity.NewDomains() {
		err := suite.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

func (suite *DBSQLiteSuite) TearDownSuite() {
	err := os.Remove(suite.DBName)
	suite.Assert().Nil(err)
	err = os.Unsetenv("DB_NAME")
	suite.Assert().Nil(err)
}
