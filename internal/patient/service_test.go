package patient

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSearchScopesResultsToAuthenticatedHospital(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	database, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, DriverName: "postgres"}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "patients" WHERE hospital_id = $1 AND national_id = $2`)).
		WithArgs(uint(7), "123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "hospital_id", "national_id"}).AddRow(1, 7, "123"))

	patients, err := NewService(database, nil).Search(SearchRequest{NationalID: "123"}, 7)
	require.NoError(t, err)
	require.Len(t, patients, 1)
	require.Equal(t, uint(7), patients[0].HospitalID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchRejectsInvalidDate(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	database, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, DriverName: "postgres"}), &gorm.Config{})
	require.NoError(t, err)

	_, err = NewService(database, nil).Search(SearchRequest{DateOfBirth: "not-a-date"}, 7)
	require.Error(t, err)
}
