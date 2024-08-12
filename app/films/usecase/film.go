package filmusecase

import (
	"errors"
	"kinopoisk/app/entity"
	errorapp "kinopoisk/app/errors"
	filmrepo "kinopoisk/app/films/repo/mysql"
	"sync"
	"time"
)

type FilmUseCase interface {
	GetFilms(genre, country, producer string) ([]*entity.Film, error)
	GetFilmByID(filmID uint64) (*entity.Film, error)
	GetFilmsByActor(ID uint64) ([]*entity.Film, error)
	GetSoonFilms() ([]*entity.Film, error)
	GetFavouriteFilms(userID uint64) ([]*entity.Film, error)
	AddFavouriteFilm(userID, filmID uint64) (bool, error)
	DeleteFavouriteFilm(userID, filmID uint64) (bool, error)
	GetFilmActors(filmID uint64) ([]*entity.Actor, error)
	GetFilmGenres(filmID uint64) ([]*entity.Genre, error)
}

type FilmUseCaseStruct struct {
	mu       *sync.RWMutex
	FilmRepo filmrepo.FilmRepo
}

func NewFilmUseCaseStruct(filmRepo filmrepo.FilmRepo) *FilmUseCaseStruct {
	return &FilmUseCaseStruct{
		mu:       &sync.RWMutex{},
		FilmRepo: filmRepo,
	}
}

func (f *FilmUseCaseStruct) GetFilms(genre, country, producer string) ([]*entity.Film, error) {
	f.mu.RLock()
	films, err := f.FilmRepo.GetFilmsRepo(genre, country, producer)
	f.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	return films, nil
}

func (f *FilmUseCaseStruct) GetFilmByID(filmID uint64) (*entity.Film, error) {
	f.mu.RLock()
	film, err := f.FilmRepo.GetFilmByIDRepo(filmID)
	f.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	if film == nil {
		return nil, nil
	}
	return film, nil
}

func (f *FilmUseCaseStruct) GetFilmsByActor(id uint64) ([]*entity.Film, error) {
	f.mu.RLock()
	films, err := f.FilmRepo.GetFilmsByActorRepo(id)
	f.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	return films, nil
}

func (f *FilmUseCaseStruct) GetSoonFilms() ([]*entity.Film, error) {
	f.mu.RLock()
	currentDate := time.Now().Format("2006-01-02")
	f.mu.RUnlock()
	films, err := f.FilmRepo.GetSoonFilmsRepo(currentDate)
	if err != nil {
		return nil, err
	}
	return films, nil
}

func (f *FilmUseCaseStruct) GetFavouriteFilms(userID uint64) ([]*entity.Film, error) {
	f.mu.RLock()
	films, err := f.FilmRepo.GetFavouriteFilmsRepo(userID)
	f.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	return films, nil
}

func (f *FilmUseCaseStruct) AddFavouriteFilm(userID, filmID uint64) (bool, error) {
	f.mu.RLock()
	film, err := f.FilmRepo.GetFilmByIDRepo(filmID)
	f.mu.RUnlock()
	if err != nil {
		return false, err
	}
	if film == nil {
		return false, errorapp.ErrorNoFilm
	}
	f.mu.RLock()
	_, err = f.FilmRepo.GetFilmInFavourites(filmID, userID)
	f.mu.RUnlock()
	if err == nil {
		return false, nil
	}
	if err != nil {
		if !errors.Is(err, errorapp.ErrorNoFilm) {
			return false, err
		}
	}
	f.mu.Lock()
	wasAdded, err := f.FilmRepo.AddFavouriteFilmRepo(userID, filmID)
	f.mu.Unlock()
	if err != nil {
		return false, err
	}
	return wasAdded, nil
}

func (f *FilmUseCaseStruct) DeleteFavouriteFilm(userID, filmID uint64) (bool, error) {
	f.mu.RLock()
	ID, err := f.FilmRepo.GetFilmInFavourites(filmID, userID)
	f.mu.RUnlock()
	if errors.Is(err, errorapp.ErrorNoFilm) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	f.mu.Lock()
	wasDeleted, err := f.FilmRepo.DeleteFavouriteFilmRepo(ID)
	f.mu.Unlock()
	if err != nil {
		return false, err
	}
	return wasDeleted, nil
}

func (f *FilmUseCaseStruct) GetFilmActors(filmID uint64) ([]*entity.Actor, error) {
	f.mu.RLock()
	film, err := f.GetFilmByID(filmID)
	f.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	if film == nil {
		return nil, errorapp.ErrorNoFilm
	}
	f.mu.RLock()
	actors, err := f.FilmRepo.GetFilmActorsRepo(filmID)
	f.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	return actors, nil
}

func (f *FilmUseCaseStruct) GetFilmGenres(filmID uint64) ([]*entity.Genre, error) {
	f.mu.RLock()
	film, err := f.GetFilmByID(filmID)
	f.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	if film == nil {
		return nil, errorapp.ErrorNoFilm
	}
	f.mu.RLock()
	genres, err := f.FilmRepo.GetFilmGenresRepo(filmID)
	f.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	return genres, nil
}
