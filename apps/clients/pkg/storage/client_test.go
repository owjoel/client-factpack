package storage

// import (
// 	"database/sql"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// type ClientStorageTestSuite struct {
// 	suite.Suite
// 	mockDB *sql.DB
// 	mock   sqlmock.Sqlmock
// 	db     *gorm.DB
// 	store  *ClientStorage
// }

// func (suite *ClientStorageTestSuite) SetupSuite() {
// 	mockDB, mock, err := sqlmock.New()
// 	assert.NoError(suite.T(), err, "Mock database should initialize without error")

// 	mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(
// 		sqlmock.NewRows([]string{"VERSION()"}).AddRow("test-version"),
// 	)

// 	suite.mockDB = mockDB
// 	suite.mock = mock

// 	suite.db, err = gorm.Open(mysql.New(mysql.Config{
// 		Conn: mockDB,
// 	}), &gorm.Config{})
// 	assert.NoError(suite.T(), err, "Mock GORM database should initialize without error")

// 	suite.store = &ClientStorage{suite.db}
// }

// func (suite *ClientStorageTestSuite) TearDownSuite() {
// 	suite.mockDB.Close()
// }

// func (suite *ClientStorageTestSuite) TestCreateClient() {
// 	suite.mock.ExpectBegin()
// 	suite.mock.ExpectExec(".*").
// 		// created_at, updated_at, deleted_at
// 		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "DonaldTrump", 30, "USA", "Active").
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	suite.mock.ExpectCommit()

// 	client := &model.Client{Name: "DonaldTrump", Age: 30, Nationality: "USA", Status: "Active"}
// 	err := suite.store.Create(client)
// 	suite.NoError(err, "should not return an error")

// 	err = suite.mock.ExpectationsWereMet()
// 	suite.NoError(err, "All expectations should be met")
// }

// func (suite *ClientStorageTestSuite) TestGetClient() {
// 	suite.mock.ExpectQuery(".*").
// 		WithArgs(1, sqlmock.AnyArg()). // since it expects id, limit
// 		WillReturnRows(sqlmock.NewRows([]string{"name", "age", "nationality", "status", "created_at", "updated_at", "deleted_at"}).
// 			AddRow("Donald Trump", 30, "USA", "Active", time.Now(), time.Now(), nil)) // NOTE: will be deleted after every test

// 	client, err := suite.store.Get(1)
// 	suite.NoError(err, "Get should not return an error")
// 	suite.NotNil(client, "Client should not be nil")
// 	suite.Equal("Donald Trump", client.Name)
// 	suite.Equal(uint(30), client.Age)
// 	suite.Equal("USA", client.Nationality)
// 	suite.Equal("Active", client.Status)

// 	err = suite.mock.ExpectationsWereMet()
// 	suite.NoError(err, "All expectations should be met")
// }

// func (suite *ClientStorageTestSuite) TestUpdateClient() {
// 	suite.mock.ExpectBegin()
// 	suite.mock.ExpectExec(".*").
// 		WithArgs(1, sqlmock.AnyArg(), "Elon Musk", 28, "Canada", 1).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	suite.mock.ExpectCommit()

// 	client := &model.Client{ID: 1, Name: "Elon Musk", Age: 28, Nationality: "Canada"}
// 	err := suite.store.Update(client)
// 	suite.NoError(err, "Update should not return an error")

// 	err = suite.mock.ExpectationsWereMet()
// 	suite.NoError(err, "All expectations should be met")
// }

// func TestClientStorageSuite(t *testing.T) {
// 	suite.Run(t, new(ClientStorageTestSuite))
// }
