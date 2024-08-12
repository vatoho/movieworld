package filmrepo

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"kinopoisk/app/entity"
	errorapp "kinopoisk/app/errors"
)

type FilmRepo interface {
	GetFilmsRepo(genre, country, producer string) ([]*entity.Film, error)
	GetFilmByIDRepo(filmID uint64) (*entity.Film, error)
	GetFilmsByActorRepo(ID uint64) ([]*entity.Film, error)
	GetSoonFilmsRepo(date string) ([]*entity.Film, error)
	GetFavouriteFilmsRepo(userID uint64) ([]*entity.Film, error)
	AddFavouriteFilmRepo(userID, filmID uint64) (bool, error)
	DeleteFavouriteFilmRepo(ID uint64) (bool, error)
	GetFilmActorsRepo(filmID uint64) ([]*entity.Actor, error)
	GetFilmGenresRepo(filmID uint64) ([]*entity.Genre, error)
	GetFilmInFavourites(filmID, userID uint64) (uint64, error)
}

type FilmRepoMySQL struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewFilmRepoMySQL(db *sql.DB, logger *zap.SugaredLogger) *FilmRepoMySQL {
	return &FilmRepoMySQL{
		db:     db,
		logger: logger,
	}
}

func (r *FilmRepoMySQL) GetFilmsRepo(genre, country, producer string) ([]*entity.Film, error) {
	var args []interface{}
	query := "SELECT f.id, f.name, f.description, f.duration, f.min_age, f.country, f.producer_name, f.date_of_release, f.num_of_marks, f.rating from films f"
	if genre != "" {
		query += " INNER JOIN film_genres fg ON f.id = fg.film_id INNER JOIN genres g ON g.id = fg.genre_id WHERE g.name = ?"
		args = append(args, genre)
	} else {
		query += " WHERE 1 = 1"
	}
	if country != "" {
		args = append(args, country)
		query += " AND f.country = ?"
	}
	if producer != "" {
		args = append(args, producer)
		query += " AND f.producer_name = ?"
	}
	rows, err := r.db.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.logger.Errorf("error in closing db rows in mysql")
		}
	}(rows)
	films := make([]*entity.Film, 0)
	for rows.Next() {
		film := &entity.Film{}
		err = rows.Scan(&film.ID, &film.Name, &film.Description, &film.Duration, &film.MinAge, &film.Country,
			&film.ProducerName, &film.DateOfRelease, &film.NumOfMarks, &film.Rating)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *FilmRepoMySQL) GetFilmByIDRepo(filmID uint64) (*entity.Film, error) {
	film := &entity.Film{}
	err := r.db.
		QueryRow("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE id = ?", filmID).
		Scan(&film.ID, &film.Name, &film.Description, &film.Duration, &film.MinAge, &film.Country, &film.ProducerName, &film.DateOfRelease, &film.NumOfMarks, &film.Rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return film, nil
}

func (r *FilmRepoMySQL) GetFilmsByActorRepo(id uint64) ([]*entity.Film, error) {
	films := []*entity.Film{}
	rows, err := r.db.Query(`SELECT f.id, f.name, f.description, f.duration, f.min_age, f.country, f.producer_name, f.date_of_release, f.num_of_marks, f.rating
FROM films f INNER JOIN actor_films af ON f.id = af.film_id INNER JOIN actors a ON a.id = af.actor_id WHERE a.id = ?`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.logger.Errorf("error in closing db rows")
		}
	}(rows)
	for rows.Next() {
		film := &entity.Film{}
		err = rows.Scan(&film.ID, &film.Name, &film.Description, &film.Duration, &film.MinAge, &film.Country, &film.ProducerName, &film.DateOfRelease, &film.NumOfMarks, &film.Rating)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *FilmRepoMySQL) GetSoonFilmsRepo(date string) ([]*entity.Film, error) {
	films := []*entity.Film{}
	rows, err := r.db.Query("SELECT id, name, description, duration, min_age, country, producer_name, date_of_release, num_of_marks, rating FROM films WHERE date_of_release > ?", date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.logger.Errorf("error in closing db rows")
		}
	}(rows)
	for rows.Next() {
		film := &entity.Film{}
		err = rows.Scan(&film.ID, &film.Name, &film.Description, &film.Duration, &film.MinAge, &film.Country, &film.ProducerName, &film.DateOfRelease, &film.NumOfMarks, &film.Rating)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *FilmRepoMySQL) GetFavouriteFilmsRepo(userID uint64) ([]*entity.Film, error) {
	films := []*entity.Film{}
	rows, err := r.db.Query("SELECT f.id, f.name, f.description, f.duration, f.min_age, f.country, f.producer_name, f.date_of_release, f.num_of_marks, f.rating FROM films f JOIN favourite_films ff on f.id = ff.film_id WHERE ff.user_id = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.logger.Errorf("error in closing db rows")
		}
	}(rows)
	for rows.Next() {
		film := &entity.Film{}
		err = rows.Scan(&film.ID, &film.Name, &film.Description, &film.Duration, &film.MinAge, &film.Country, &film.ProducerName, &film.DateOfRelease, &film.NumOfMarks, &film.Rating)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *FilmRepoMySQL) AddFavouriteFilmRepo(userID, filmID uint64) (bool, error) {
	_, err := r.db.Exec(
		"INSERT INTO favourite_films (`user_id`, `film_id`) VALUES (?, ?)",
		userID,
		filmID,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *FilmRepoMySQL) DeleteFavouriteFilmRepo(id uint64) (bool, error) {
	_, err := r.db.Exec(
		"DELETE FROM favourite_films WHERE id = ?",
		id,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *FilmRepoMySQL) GetFilmActorsRepo(filmID uint64) ([]*entity.Actor, error) {
	actors := []*entity.Actor{}
	rows, err := r.db.Query("SELECT a.id, a.name, a.surname, a.nationality, a.birthday FROM actors a INNER JOIN actor_films af ON a.id = af.actor_id WHERE af.film_id = ?", filmID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.logger.Errorf("error in closing db rows")
		}
	}(rows)
	for rows.Next() {
		actor := &entity.Actor{}
		err = rows.Scan(&actor.ID, &actor.Name, &actor.Surname, &actor.Nationality, &actor.Birthday)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}
	return actors, nil
}

func (r *FilmRepoMySQL) GetFilmGenresRepo(filmID uint64) ([]*entity.Genre, error) {
	genres := []*entity.Genre{}
	rows, err := r.db.Query("SELECT g.id, g.name FROM genres g INNER JOIN film_genres fg ON g.id = fg.genre_id WHERE fg.film_id = ?", filmID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			r.logger.Errorf("error in closing db rows")
		}
	}(rows)
	for rows.Next() {
		genre := &entity.Genre{}
		err = rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (r *FilmRepoMySQL) GetFilmInFavourites(filmID, userID uint64) (uint64, error) {
	var id uint64
	err := r.db.
		QueryRow("SELECT id from favourite_films WHERE user_id = ? AND film_id = ?", userID, filmID).
		Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errorapp.ErrorNoFilm
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
