package filmusecase_test

import (
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"kinopoisk/app/entity"
	errorapp "kinopoisk/app/errors"
	filmrepo "kinopoisk/app/films/repo/mysql"
	filmusecase "kinopoisk/app/films/usecase"
	"reflect"
	"testing"
)

func TestGetFilms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can not create mock")
	}
	defer db.Close()
	dbRepo := filmrepo.NewFilmRepoMySQL(db, zap.NewNop().Sugar())
	testUsecase := filmusecase.NewFilmUseCaseStruct(dbRepo)

	// какая то ошибка базы данных
	genre := "drama"
	mock.
		ExpectQuery("SELECT f.id, f.name, f.description, f.duration, f.min_age, f.country, f.producer_name, f.date_of_release, f.num_of_marks, f.rating from films f INNER JOIN film_genres fg ON f.id = fg.film_id INNER JOIN genres g ON g.id = fg.genre_id WHERE").
		WithArgs(genre).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = testUsecase.GetFilms("drama", "", "")
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	expectedFilm := &entity.Film{
		ID:            1,
		Name:          "Titanic",
		Description:   "fvegfvreggev",
		Duration:      212,
		MinAge:        6,
		Country:       "USA",
		ProducerName:  "dwdw",
		DateOfRelease: "2012-12-12",
		NumOfMarks:    1,
		Rating:        7.0,
	}
	expectedFilms := []*entity.Film{expectedFilm}
	rows := sqlmock.NewRows([]string{"id", "name", "description", "duration", "min_age", "country", "producer_name", "date_of_release", "num_of_marks", "rating"})
	for _, currentFilm := range expectedFilms {
		rows = rows.AddRow(currentFilm.ID, currentFilm.Name, currentFilm.Description, currentFilm.Duration, currentFilm.MinAge,
			currentFilm.Country, currentFilm.ProducerName, currentFilm.DateOfRelease, currentFilm.NumOfMarks, currentFilm.Rating)
	}
	mock.
		ExpectQuery("SELECT f.id, f.name, f.description, f.duration, f.min_age, f.country, f.producer_name, f.date_of_release, f.num_of_marks, f.rating from films f INNER JOIN film_genres fg ON f.id = fg.film_id INNER JOIN genres g ON g.id = fg.genre_id WHERE").
		WithArgs(genre).
		WillReturnRows(rows)

	films, err := testUsecase.GetFilms("drama", "", "")
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedFilms, films) {
		t.Errorf("wrong result: expected %v, got %v", expectedFilms, films)
		return
	}

}

func TestGetFilmActors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can not create mock")
	}
	defer db.Close()
	dbRepo := filmrepo.NewFilmRepoMySQL(db, zap.NewNop().Sugar())
	testUsecase := filmusecase.NewFilmUseCaseStruct(dbRepo)

	// какая то ошибка базы данных
	var id uint64 = 1
	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(id).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = testUsecase.GetFilmActors(id)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// нет фильма с таким айди
	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err = testUsecase.GetFilmActors(id)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !errors.Is(err, errorapp.ErrorNoFilm) {
		t.Errorf("wrong error: expected %s, got %s", errorapp.ErrorNoFilm, err)
		return
	}

	// фильм есть, ошибка при запросе актеров
	expectedFilm := &entity.Film{
		ID:            1,
		Name:          "Titanic",
		Description:   "fvegfvreggev",
		Duration:      212,
		MinAge:        6,
		Country:       "USA",
		ProducerName:  "dwdw",
		DateOfRelease: "2012-12-12",
		NumOfMarks:    1,
		Rating:        7.0,
	}
	expectedFilms := []*entity.Film{expectedFilm}
	rows := sqlmock.NewRows([]string{"id", "name", "description", "duration", "min_age", "country", "producer_name", "date_of_release", "num_of_marks", "rating"})
	for _, currentFilm := range expectedFilms {
		rows = rows.AddRow(currentFilm.ID, currentFilm.Name, currentFilm.Description, currentFilm.Duration, currentFilm.MinAge,
			currentFilm.Country, currentFilm.ProducerName, currentFilm.DateOfRelease, currentFilm.NumOfMarks, currentFilm.Rating)
	}
	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(id).
		WillReturnRows(rows)
	mock.
		ExpectQuery("SELECT a.id, a.name, a.surname, a.nationality, a.birthday FROM actors a INNER JOIN actor_films af ON a.id = af.actor_id WHERE").
		WithArgs(expectedFilms[0].ID).
		WillReturnError(fmt.Errorf("db error"))

	_, err = testUsecase.GetFilmActors(id)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// все ок
	expectedActors := []*entity.Actor{
		{
			ID:          1,
			Name:        "Ivan",
			Surname:     "Ivanov",
			Nationality: "Russia",
			Birthday:    "2000-12-12",
		},
	}
	actorsRows := sqlmock.NewRows([]string{"id", "name", "surname", "nationality", "birthday"})
	for _, currentActor := range expectedActors {
		actorsRows = actorsRows.AddRow(currentActor.ID, currentActor.Name, currentActor.Surname, currentActor.Nationality, currentActor.Birthday)
	}
	rows = sqlmock.NewRows([]string{"id", "name", "description", "duration", "min_age", "country", "producer_name", "date_of_release", "num_of_marks", "rating"})
	for _, currentFilm := range expectedFilms {
		rows = rows.AddRow(currentFilm.ID, currentFilm.Name, currentFilm.Description, currentFilm.Duration, currentFilm.MinAge,
			currentFilm.Country, currentFilm.ProducerName, currentFilm.DateOfRelease, currentFilm.NumOfMarks, currentFilm.Rating)
	}
	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(id).
		WillReturnRows(rows)
	mock.
		ExpectQuery("SELECT a.id, a.name, a.surname, a.nationality, a.birthday FROM actors a INNER JOIN actor_films af ON a.id = af.actor_id WHERE").
		WithArgs(expectedFilms[0].ID).
		WillReturnRows(actorsRows)

	_, err = testUsecase.GetFilmActors(id)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

}

func TestAddFavouriteFilm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can not create mock")
	}
	defer db.Close()
	dbRepo := filmrepo.NewFilmRepoMySQL(db, zap.NewNop().Sugar())
	testUsecase := filmusecase.NewFilmUseCaseStruct(dbRepo)

	// какая то ошибка базы данных
	var userID uint64 = 1
	var filmID uint64 = 1
	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(filmID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = testUsecase.AddFavouriteFilm(userID, filmID)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(filmID).
		WillReturnError(sql.ErrNoRows)

	_, err = testUsecase.AddFavouriteFilm(userID, filmID)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !errors.Is(err, errorapp.ErrorNoFilm) {
		t.Errorf("wrong error: expected %s, got %s", errorapp.ErrorNoFilm, err)
		return
	}

	expectedFilm := &entity.Film{
		ID:            1,
		Name:          "Titanic",
		Description:   "fvegfvreggev",
		Duration:      212,
		MinAge:        6,
		Country:       "USA",
		ProducerName:  "dwdw",
		DateOfRelease: "2012-12-12",
		NumOfMarks:    1,
		Rating:        7.0,
	}
	expectedFilms := []*entity.Film{expectedFilm}
	rows := sqlmock.NewRows([]string{"id", "name", "description", "duration", "min_age", "country", "producer_name", "date_of_release", "num_of_marks", "rating"})
	for _, currentFilm := range expectedFilms {
		rows = rows.AddRow(currentFilm.ID, currentFilm.Name, currentFilm.Description, currentFilm.Duration, currentFilm.MinAge,
			currentFilm.Country, currentFilm.ProducerName, currentFilm.DateOfRelease, currentFilm.NumOfMarks, currentFilm.Rating)
	}
	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(filmID).
		WillReturnRows(rows)
	mock.
		ExpectQuery("SELECT id from favourite_films WHERE").
		WithArgs(userID, filmID).
		WillReturnError(fmt.Errorf("db error"))

	_, err = testUsecase.AddFavouriteFilm(userID, filmID)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	rows = sqlmock.NewRows([]string{"id", "name", "description", "duration", "min_age", "country", "producer_name", "date_of_release", "num_of_marks", "rating"})
	for _, currentFilm := range expectedFilms {
		rows = rows.AddRow(currentFilm.ID, currentFilm.Name, currentFilm.Description, currentFilm.Duration, currentFilm.MinAge,
			currentFilm.Country, currentFilm.ProducerName, currentFilm.DateOfRelease, currentFilm.NumOfMarks, currentFilm.Rating)
	}
	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(filmID).
		WillReturnRows(rows)

	idRow := sqlmock.NewRows([]string{"id"})
	idRow = idRow.AddRow(filmID)
	mock.
		ExpectQuery("SELECT id from favourite_films WHERE").
		WithArgs(userID, filmID).
		WillReturnRows(idRow)

	wasAdded, err := testUsecase.AddFavouriteFilm(userID, filmID)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if wasAdded {
		t.Errorf("wrong result for wasAdded: expected %v, got %v", !wasAdded, wasAdded)
		return
	}

	rows = sqlmock.NewRows([]string{"id", "name", "description", "duration", "min_age", "country", "producer_name", "date_of_release", "num_of_marks", "rating"})
	for _, currentFilm := range expectedFilms {
		rows = rows.AddRow(currentFilm.ID, currentFilm.Name, currentFilm.Description, currentFilm.Duration, currentFilm.MinAge,
			currentFilm.Country, currentFilm.ProducerName, currentFilm.DateOfRelease, currentFilm.NumOfMarks, currentFilm.Rating)
	}
	mock.
		ExpectQuery("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE").
		WithArgs(filmID).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT id from favourite_films WHERE").
		WithArgs(userID, filmID).
		WillReturnError(sql.ErrNoRows)
	mock.
		ExpectExec(`INSERT INTO favourite_films`).
		WithArgs(userID, filmID).
		WillReturnError(fmt.Errorf("db error"))
	_, err = testUsecase.AddFavouriteFilm(userID, filmID)
	if err := mock.ExpectationsWereMet(); err != nil { // nolint govet
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
