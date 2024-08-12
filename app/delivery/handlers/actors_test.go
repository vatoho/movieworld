package handlers

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	actorusecase "kinopoisk/app/actors/usecase"
	"kinopoisk/app/entity"
	"kinopoisk/app/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()
	ctx := context.WithValue(context.Background(), middleware.MyLoggerKey, logger)
	testUseCase := actorusecase.NewMockActorUseCase(ctrl)
	testHandler := NewActorHandler(testUseCase)
	// usecase returns error
	testUseCase.EXPECT().GetActors().Return(nil, fmt.Errorf("error"))
	request := httptest.NewRequest(http.MethodGet, "/actors/", nil)
	respWriter := httptest.NewRecorder()
	testHandler.GetActors(respWriter, request.WithContext(ctx))
	resp := respWriter.Result()
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
		return
	}

	// usecase returns actors without error
	actors := []*entity.Actor{
		{
			ID:          1,
			Name:        "Sergey",
			Surname:     "Burunov",
			Nationality: "Russian",
			Birthday:    "2012-12-12",
		},
	}
	testUseCase.EXPECT().GetActors().Return(actors, nil)
	request = httptest.NewRequest(http.MethodGet, "/actors/", nil)
	respWriter = httptest.NewRecorder()
	testHandler.GetActors(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status %d, got status %d", http.StatusOK, resp.StatusCode)
		return
	}
}

func TestGetActorByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()
	testUseCase := actorusecase.NewMockActorUseCase(ctrl)
	testHandler := NewActorHandler(testUseCase)
	// bad id format
	request := httptest.NewRequest(http.MethodGet, "/actor/bad", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "bad"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request.WithContext(ctx))
	resp := respWriter.Result()
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got status %d", http.StatusBadRequest, resp.StatusCode)
		return
	}

	// usecase returns error
	var ID uint64 = 1
	testUseCase.EXPECT().GetActorByID(ID).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/actor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status %d, got status %d", http.StatusInternalServerError, resp.StatusCode)
		return
	}

	// usecase returns nil actor

	testUseCase.EXPECT().GetActorByID(ID).Return(nil, nil)
	request = httptest.NewRequest(http.MethodGet, "/actor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 404 {
		t.Errorf("expected status %d, got status %d", http.StatusNotFound, resp.StatusCode)
		return
	}

	// all is ok
	actor := &entity.Actor{
		ID:          1,
		Name:        "Sergey",
		Surname:     "Burunov",
		Nationality: "Russian",
		Birthday:    "2012-12-12",
	}
	testUseCase.EXPECT().GetActorByID(ID).Return(actor, nil)
	request = httptest.NewRequest(http.MethodGet, "/actor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"ACTOR_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetActorByID(respWriter, request.WithContext(ctx))
	resp = respWriter.Result()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body")
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status %d, got status %d", http.StatusOK, resp.StatusCode)
		return
	}
}
