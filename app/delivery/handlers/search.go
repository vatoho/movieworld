package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"kinopoisk/app/delivery"
	"kinopoisk/app/middleware"
	searchusecase "kinopoisk/app/search/usecase"
	"log"
	"net/http"
	"unicode/utf8"
)

type SearchHandler struct {
	searchUseCases searchusecase.SearchUseCase
}

func NewSearchHandler(searchUseCases searchusecase.SearchUseCase) *SearchHandler {
	return &SearchHandler{
		searchUseCases: searchUseCases,
	}
}

func (sh *SearchHandler) MakeSearch(w http.ResponseWriter, r *http.Request) {
	logger, err := middleware.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		middleware.WriteNoLoggerResponse(w)
	}
	vars := mux.Vars(r)
	searchData := vars["DATA"]
	logger.Infof("search for: %s", searchData)
	if utf8.RuneCountInString(searchData) < 3 {
		logger.Errorf("at least 3 letters in query required, but there was passed %s", searchData)
		delivery.WriteResponse(logger, w, []byte(`{"message":"at least 3 letters in query required"}`), http.StatusBadRequest)
		return
	}
	result, err := sh.searchUseCases.MakeSearch(searchData, logger)
	if err != nil {
		logger.Errorf("error in making search: %s", err)
		delivery.WriteResponse(logger, w, []byte(`{"message":"search error"}`), http.StatusInternalServerError)
		return
	}
	if len(result.Actors) == 0 && len(result.Films) == 0 {
		logger.Infof("no films and actors for search %s", searchData)
		delivery.WriteResponse(logger, w, []byte(`{"message":"search found nothing"}`), http.StatusNotFound)
		return
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		logger.Errorf("error in json coding result: %s", err)
		delivery.WriteResponse(logger, w, []byte(`{"message":"internal error"}`), http.StatusInternalServerError)
		return
	}
	delivery.WriteResponse(logger, w, resultJSON, http.StatusOK)
}
