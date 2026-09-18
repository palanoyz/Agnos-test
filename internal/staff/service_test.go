package staff

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newTestService(t *testing.T) (*Service, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	database, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	}), &gorm.Config{})
	require.NoError(t, err)
	return NewService(database), mock
}

func TestCreateStaff(t *testing.T) {
	service, mock := newTestService(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE name = $1 AND "hospitals"."name" = $2 ORDER BY "hospitals"."id" LIMIT $3`)).
		WithArgs("Hospital A", "Hospital A", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "hospitals" ("name","created_at","updated_at") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs("Hospital A", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs" WHERE username = $1 AND hospital_id = $2 ORDER BY "staffs"."id" LIMIT $3`)).
		WithArgs("alice", 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "hospital_id"}))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "staffs" ("username","password_hash","hospital_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs("alice", sqlmock.AnyArg(), 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	staffMember, err := service.Create(CreateRequest{
		Username: "alice",
		Password: "secret",
		Hospital: "Hospital A",
	})

	require.NoError(t, err)
	require.Equal(t, "alice", staffMember.Username)
	require.NotEmpty(t, staffMember.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(staffMember.PasswordHash), []byte("secret")))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStaffRejectsDuplicateInSameHospital(t *testing.T) {
	service, mock := newTestService(t)
	request := CreateRequest{Username: "alice", Password: "secret", Hospital: "Hospital A"}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE name = $1 AND "hospitals"."name" = $2 ORDER BY "hospitals"."id" LIMIT $3`)).
		WithArgs("Hospital A", "Hospital A", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "hospitals" ("name","created_at","updated_at") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs("Hospital A", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs" WHERE username = $1 AND hospital_id = $2 ORDER BY "staffs"."id" LIMIT $3`)).
		WithArgs("alice", 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "hospital_id"}))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "staffs" ("username","password_hash","hospital_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs("alice", sqlmock.AnyArg(), 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	_, err := service.Create(request)
	require.NoError(t, err)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE name = $1 AND "hospitals"."name" = $2 ORDER BY "hospitals"."id" LIMIT $3`)).
		WithArgs("Hospital A", "Hospital A", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Hospital A"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs" WHERE username = $1 AND hospital_id = $2 ORDER BY "staffs"."id" LIMIT $3`)).
		WithArgs("alice", 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "hospital_id"}).AddRow(1, "alice", "hash", 1))
	_, err = service.Create(request)
	require.ErrorIs(t, err, ErrStaffAlreadyExists)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStaffRejectsMissingFields(t *testing.T) {
	service, _ := newTestService(t)

	_, err := service.Create(CreateRequest{Username: "alice", Password: "secret"})
	require.ErrorIs(t, err, ErrInvalidInput)
}
