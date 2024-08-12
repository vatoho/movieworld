package searchrepo

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"kinopoisk/app/entity"
)

type SearchRepo interface {
	MakeSearchFilms(inputStr string) ([]*entity.Film, error)
	MakeSearchActors(inputStr string) ([]*entity.Actor, error)
}

type SearchRepoMySQL struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewSearchRepoMySQL(db *sql.DB, logger *zap.SugaredLogger) *SearchRepoMySQL {
	return &SearchRepoMySQL{
		db:     db,
		logger: logger,
	}
}

func (sr *SearchRepoMySQL) MakeSearchFilms(inputStr string) ([]*entity.Film, error) {
	rows, err := sr.db.Query(`
    SELECT id, name, description, duration, min_age, country, 
           producer_name, date_of_release, num_of_marks, rating 
    FROM films
    WHERE LOWER(name) LIKE LOWER(CONCAT('%', ?, '%'))
       OR LOWER(producer_name) LIKE LOWER(CONCAT('%', ?, '%'))`, inputStr, inputStr)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sr.logger.Errorf("no films found for keyword: %s", err)
			return []*entity.Film{}, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			sr.logger.Errorf("error in closing rows: %s", err)
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

func (sr *SearchRepoMySQL) MakeSearchActors(inputStr string) ([]*entity.Actor, error) {
	rows, err := sr.db.Query(`
    SELECT id, name, surname, nationality, birthday 
    FROM actors
    WHERE LOWER(CONCAT(name, ' ', surname)) LIKE LOWER(CONCAT('%', ?, '%'))
       OR LOWER(CONCAT(surname, ' ', name)) LIKE LOWER(CONCAT('%', ?, '%'))`, inputStr, inputStr)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sr.logger.Errorf("no actors found for keyword: %s", err)
			return []*entity.Actor{}, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			sr.logger.Errorf("error in closing rows: %s", err)
		}
	}(rows)
	actors := make([]*entity.Actor, 0)
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
