package storage

// import (
// 	"database/sql"
// 	"fmt"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/owjoel/client-factpack/apps/clients/config"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// type DatabaseTestSuite struct {
// 	suite.Suite
// 	mockDB *sql.DB
// 	mock   sqlmock.Sqlmock
// 	db     *gorm.DB
// }

// // set up a database mock for the test suite
// func (suite *DatabaseTestSuite) SetupSuite() {
// 	mockDB, mock, err := sqlmock.New()
// 	assert.NoError(suite.T(), err, "Should create a mock database successfully")

// 	mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(
// 		sqlmock.NewRows([]string{"VERSION()"}).AddRow("test-version"),
// 	)

// 	suite.mockDB = mockDB
// 	suite.mock = mock

// 	suite.db, err = gorm.Open(mysql.New(mysql.Config{
// 		Conn: mockDB,
// 	}), &gorm.Config{})
// 	assert.NoError(suite.T(), err, "Should open a mock GORM database")
// }

// func (suite *DatabaseTestSuite) TearDownSuite() {
// 	suite.mockDB.Close() // close mock database after all tests have completed
// }

// func (suite *DatabaseTestSuite) TestGetDSN() {
// 	config.DBUser = "testuser"
// 	config.DBPassword = "testpass"
// 	config.DBHost = "127.0.0.1"
// 	config.DBPort = "3306"
// 	config.DBName = "testdb"

// 	expectedDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
// 		"testuser",
// 		"testpass",
// 		"127.0.0.1",
// 		"3306",
// 		"testdb",
// 	)

// 	suite.Equal(expectedDSN, GetDSN(), "DSN should be formatted correctly")
// }

// func (suite *DatabaseTestSuite) TestGetInstance() {
// 	db = &storage{Client: &ClientStorage{suite.db}}
// 	instance := GetInstance()

// 	suite.Equal(db, instance, "GetInstance should return the same storage instance")
// }

// func (suite *DatabaseTestSuite) TestInit() {
// 	mockDB, mock, err := sqlmock.New()
// 	assert.NoError(suite.T(), err, "Should create a mock database successfully")

// 	mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(
// 		sqlmock.NewRows([]string{"VERSION()"}).AddRow("test-version"),
// 	)

// 	defer mockDB.Close()

// 	_db, err := gorm.Open(mysql.New(mysql.Config{
// 		Conn: mockDB,
// 	}), &gorm.Config{})

// 	assert.NoError(suite.T(), err, "Should open a mock GORM database")

// 	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))

// 	db = &storage{Client: &ClientStorage{_db}}

// 	suite.NotNil(db, "Storage instance should be initialized")
// 	suite.NotNil(db.Client, "Client interface should not be nil")
// }

// func TestDatabaseSuite(t *testing.T) {
// 	suite.Run(t, new(DatabaseTestSuite))
// }
