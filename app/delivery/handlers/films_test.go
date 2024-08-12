package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"kinopoisk/app/entity"
	errorapp "kinopoisk/app/errors"
	filmusecase "kinopoisk/app/films/usecase"
	"kinopoisk/app/middleware"
)

func TestGetFilms(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := filmusecase.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)
	// unknown query params
	request := httptest.NewRequest(http.MethodGet, "/films?bad_param=bp", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.GetFilms(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 400 {
		t.Errorf("expected status %d, got status %d", http.StatusBadRequest, resp.StatusCode)
		return
	}

	// usecase returns error
	testUseCase.EXPECT().GetFilms("Drama", "", "").Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/films?genre=Drama", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilms(respWriter, request.WithContext(ctx))
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

	// all is ok
	films := []*entity.Film{
		{
			ID:            1,
			Name:          "Green mile",
			Description:   "interesting",
			Duration:      184,
			MinAge:        12,
			Country:       "USA",
			ProducerName:  "Ivan",
			DateOfRelease: "2012-12-12",
		},
	}
	testUseCase.EXPECT().GetFilms("Drama", "", "").Return(films, nil)
	request = httptest.NewRequest(http.MethodGet, "/films?genre=Drama", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilms(respWriter, request.WithContext(ctx))
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

func TestGetFilmByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()
	testUseCase := filmusecase.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)
	// bad film id
	request := httptest.NewRequest(http.MethodGet, "/film/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "bad_id"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 400 {
		t.Errorf("expected status %d, got status %d", http.StatusBadRequest, resp.StatusCode)
		return
	}

	// usecase returns error
	var filmID uint64 = 1
	testUseCase.EXPECT().GetFilmByID(filmID).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/film/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request.WithContext(ctx))
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

	// usecase returns nil film
	testUseCase.EXPECT().GetFilmByID(filmID).Return(nil, nil)
	request = httptest.NewRequest(http.MethodGet, "/film/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request.WithContext(ctx))
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
	film := &entity.Film{

		ID:            1,
		Name:          "Green mile",
		Description:   "interesting",
		Duration:      184,
		MinAge:        12,
		Country:       "USA",
		ProducerName:  "Ivan",
		DateOfRelease: "2012-12-12",
	}
	testUseCase.EXPECT().GetFilmByID(filmID).Return(film, nil)
	request = httptest.NewRequest(http.MethodGet, "/film/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFilmByID(respWriter, request.WithContext(ctx))
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

func TestGeyFavouriteFilms(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := filmusecase.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)
	// bad user in context
	request := httptest.NewRequest(http.MethodGet, "/films/favourite", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, "bad user")
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.GetFavouriteFilms(respWriter, request.WithContext(ctx))
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

	// usecase returns error
	var userID uint64 = 1
	testUseCase.EXPECT().GetFavouriteFilms(userID).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/films/favourite", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID: 1,
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFavouriteFilms(respWriter, request.WithContext(ctx))
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

	// all is ok
	films := []*entity.Film{
		{
			ID:            1,
			Name:          "Green mile",
			Description:   "interesting",
			Duration:      184,
			MinAge:        12,
			Country:       "USA",
			ProducerName:  "Ivan",
			DateOfRelease: "2012-12-12",
		},
	}

	testUseCase.EXPECT().GetFavouriteFilms(userID).Return(films, nil)
	request = httptest.NewRequest(http.MethodGet, "/films/favourite", nil)
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID: 1,
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetFavouriteFilms(respWriter, request.WithContext(ctx))
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

func TestAddFavouriteFilm(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := filmusecase.NewMockFilmUseCase(ctrl)
	testHandler := NewFilmHandler(testUseCase)

	// bad user in context
	request := httptest.NewRequest(http.MethodPost, "/films/favourite/1", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, "bad user")
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.AddFavouriteFilm(respWriter, request.WithContext(ctx))
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

	// bad film id
	request = httptest.NewRequest(http.MethodPost, "/films/favourite/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "bad_id"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID: 1,
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFavouriteFilm(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 400 {
		t.Errorf("expected status %d, got status %d", http.StatusBadRequest, resp.StatusCode)
		return
	}

	// no film with such id
	var (
		userID uint64 = 1
		filmID uint64 = 1
	)
	testUseCase.EXPECT().AddFavouriteFilm(userID, filmID).Return(false, errorapp.ErrorNoFilm)
	request = httptest.NewRequest(http.MethodPost, "/films/favourite/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID: 1,
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFavouriteFilm(respWriter, request.WithContext(ctx))
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

	// usecase returns error
	testUseCase.EXPECT().AddFavouriteFilm(userID, filmID).Return(false, fmt.Errorf("internal error"))
	request = httptest.NewRequest(http.MethodPost, "/films/favourite/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID: 1,
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFavouriteFilm(respWriter, request.WithContext(ctx))
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

	// film was already in favourites
	testUseCase.EXPECT().AddFavouriteFilm(userID, filmID).Return(false, nil)
	request = httptest.NewRequest(http.MethodPost, "/films/favourite/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID: 1,
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFavouriteFilm(respWriter, request.WithContext(ctx))
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

	// all is ok
	testUseCase.EXPECT().AddFavouriteFilm(userID, filmID).Return(true, nil)
	request = httptest.NewRequest(http.MethodPost, "/films/favourite/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID: 1,
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddFavouriteFilm(respWriter, request.WithContext(ctx))
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
