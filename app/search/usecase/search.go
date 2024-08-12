package searchusecase

import (
	"go.uber.org/zap"
	"kinopoisk/app/entity"
	searchrepo "kinopoisk/app/search/repo/mysql"
	"sync"
)

type SearchUseCase interface {
	MakeSearch(inputStr string, logger *zap.SugaredLogger) (*entity.SearchResult, error)
}

type SearchUseCaseStruct struct {
	mu         *sync.RWMutex
	searchRepo searchrepo.SearchRepo
}

func NewSearchUseCaseStruct(searchRepo searchrepo.SearchRepo) *SearchUseCaseStruct {
	return &SearchUseCaseStruct{
		mu:         &sync.RWMutex{},
		searchRepo: searchRepo,
	}
}

func (sr *SearchUseCaseStruct) MakeSearch(inputStr string, logger *zap.SugaredLogger) (*entity.SearchResult, error) {
	films, err := sr.searchRepo.MakeSearchFilms(inputStr)
	if err != nil {
		logger.Errorf("error in search films in db: %s", err)
		return nil, err
	}
	actors, err := sr.searchRepo.MakeSearchActors(inputStr)
	if err != nil {
		logger.Errorf("error in search actors in db: %s", err)
		return nil, err
	}
	return &entity.SearchResult{
		Films:  films,
		Actors: actors,
	}, nil
}
