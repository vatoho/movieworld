package actorusecase_test

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	actorrepo "kinopoisk/app/actors/repo/mysql"
	actorusecase "kinopoisk/app/actors/usecase"
	"kinopoisk/app/entity"
	"reflect"
	"testing"
)

func TestGetActorByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can not create mock")
	}
	defer db.Close()
	dbRepo := actorrepo.NewActorRepoMySQL(db, zap.NewNop().Sugar())
	testUsecase := actorusecase.NewActorUseCaseStruct(dbRepo)

	var id uint64 = 1
	mock.
		ExpectQuery("SELECT id, name, surname, nationality, birthday FROM actors WHERE").
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = testUsecase.GetActorByID(id)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	mock.
		ExpectQuery("SELECT id, name, surname, nationality, birthday FROM actors WHERE").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	actor, err := testUsecase.GetActorByID(id)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if actor != nil {
		t.Errorf("unexpected non nil actor: %v", actor)
		return
	}

	expectedActor := &entity.Actor{
		ID:          1,
		Name:        "Titanic",
		Surname:     "fvegfvreggev",
		Nationality: "frwfer",
		Birthday:    "2000-12-12",
	}
	expectedActors := []*entity.Actor{expectedActor}
	rows := sqlmock.NewRows([]string{"id", "name", "surname", "nationality", "birthday"})
	for _, currentActor := range expectedActors {
		rows = rows.AddRow(currentActor.ID, currentActor.Name, currentActor.Surname, currentActor.Nationality, currentActor.Birthday)
	}
	mock.
		ExpectQuery("SELECT id, name, surname, nationality, birthday FROM actors WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	actor, err = testUsecase.GetActorByID(id)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedActor, actor) {
		t.Errorf("wrong result: expected %v, got %v", expectedActor, actor)
		return
	}

}
