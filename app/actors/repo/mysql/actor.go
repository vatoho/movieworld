package actorrepo

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"kinopoisk/app/entity"
)

type ActorRepo interface {
	GetActorByIDRepo(ID uint64) (*entity.Actor, error)
	GetActorsRepo() ([]*entity.Actor, error)
}

type ActorRepoMySQL struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewActorRepoMySQL(db *sql.DB, logger *zap.SugaredLogger) *ActorRepoMySQL {
	return &ActorRepoMySQL{
		db:     db,
		logger: logger,
	}
}

func (r *ActorRepoMySQL) GetActorByIDRepo(id uint64) (*entity.Actor, error) {
	actor := &entity.Actor{}
	err := r.db.
		QueryRow("SELECT id, name, surname, nationality, birthday FROM actors WHERE id = ?", id).
		Scan(&actor.ID, &actor.Name, &actor.Surname, &actor.Nationality, &actor.Birthday)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return actor, nil
}

func (r *ActorRepoMySQL) GetActorsRepo() ([]*entity.Actor, error) {
	actors := []*entity.Actor{}
	rows, err := r.db.Query("SELECT id, name, surname, nationality, birthday FROM actors")
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
