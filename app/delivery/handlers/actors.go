package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	actorusecase "kinopoisk/app/actors/usecase"
	"kinopoisk/app/delivery"
	"kinopoisk/app/middleware"
	"log"
	"net/http"
	"strconv"
)

type ActorHandler struct {
	ActorUseCases actorusecase.ActorUseCase
}

func NewActorHandler(actorUseCases actorusecase.ActorUseCase) *ActorHandler {
	return &ActorHandler{
		ActorUseCases: actorUseCases,
	}
}

func (ah *ActorHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	logger, err := middleware.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		middleware.WriteNoLoggerResponse(w)
	}
	actors, err := ah.ActorUseCases.GetActors()
	if err != nil {
		errText := fmt.Sprintf(`{"message": "internal server error: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	actorsJSON, err := json.Marshal(actors)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in coding actors: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	delivery.WriteResponse(logger, w, actorsJSON, http.StatusOK)
}

func (ah *ActorHandler) GetActorByID(w http.ResponseWriter, r *http.Request) {
	logger, err := middleware.GetLoggerFromContext(r.Context())
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		middleware.WriteNoLoggerResponse(w)
	}
	vars := mux.Vars(r)
	actorID := vars["ACTOR_ID"]
	actorIDInt, err := strconv.ParseUint(actorID, 10, 64)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "bad format of actor id: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusBadRequest)
		return
	}
	actor, err := ah.ActorUseCases.GetActorByID(actorIDInt)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "internal server error: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	if actor == nil {
		errText := fmt.Sprintf(`{"message": "actor with ID %d is not found"}`, actorIDInt)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusNotFound)
		return
	}
	actorJSON, err := json.Marshal(actor)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in coding actor: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	delivery.WriteResponse(logger, w, actorJSON, http.StatusOK)
}
