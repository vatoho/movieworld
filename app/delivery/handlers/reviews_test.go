package handlers

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"kinopoisk/app/dto"
	"kinopoisk/app/entity"
	errorapp "kinopoisk/app/errors"
	"kinopoisk/app/middleware"
	reviewusecase "kinopoisk/app/reviews/usecase"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorReader struct{}

func (er *errorReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestGetReviewsForFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := reviewusecase.NewMockReviewUseCase(ctrl)
	testHandler := NewReviewHandler(testUseCase)
	// bad film id
	request := httptest.NewRequest(http.MethodGet, "/review/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "bad_id"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.GetReviewsForFilm(respWriter, request.WithContext(ctx))
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
	testUseCase.EXPECT().GetFilmReviews(filmID, logger).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodGet, "/review/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetReviewsForFilm(respWriter, request.WithContext(ctx))
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
	reviews := []*entity.Review{
		{
			ID:      1,
			Mark:    8,
			Comment: "good film",
			Author: &entity.User{
				ID:       1,
				Username: "vasyan",
			},
		},
	}
	testUseCase.EXPECT().GetFilmReviews(filmID, logger).Return(reviews, nil)
	request = httptest.NewRequest(http.MethodGet, "/review/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.GetReviewsForFilm(respWriter, request.WithContext(ctx))
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

func TestAddReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := reviewusecase.NewMockReviewUseCase(ctrl)
	testHandler := NewReviewHandler(testUseCase)
	// bad film id
	request := httptest.NewRequest(http.MethodPost, "/review/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "bad_id"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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

	// bad user in context
	request = httptest.NewRequest(http.MethodPost, "/review/1", nil)
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, "bad user")
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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

	// can not read request body
	request = httptest.NewRequest(http.MethodPost, "/review/1", &errorReader{})
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID:       1,
		Username: "vasya",
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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
	}

	// error in unmarshalling request body
	request = httptest.NewRequest(http.MethodPost, "/review/1", strings.NewReader("{"))
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID:       1,
		Username: "vasya",
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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
	}

	request = httptest.NewRequest(http.MethodPost, "/review/1", strings.NewReader(`{"comment":"v film"}`))
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, &entity.User{
		ID:       1,
		Username: "vasya",
	})
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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
	if resp.StatusCode != 422 {
		t.Errorf("expected status %d, got status %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}

	var filmID uint64 = 1
	author := &entity.User{
		ID:       1,
		Username: "vasya",
	}
	newReview := &dto.ReviewDTO{
		Mark:    10,
		Comment: "very interesting film",
	}
	testUseCase.EXPECT().NewReview(newReview, filmID, author, logger).Return(nil, fmt.Errorf("error"))
	request = httptest.NewRequest(http.MethodPost, "/review/1", strings.NewReader(`{"mark": 10,"comment":"very interesting film"}`))
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, author)
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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
	}

	// no film with such id
	testUseCase.EXPECT().NewReview(newReview, filmID, author, logger).Return(nil, errorapp.ErrorNoFilm)
	request = httptest.NewRequest(http.MethodPost, "/review/1", strings.NewReader(`{"mark": 10,"comment":"very interesting film"}`))
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, author)
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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
	}

	// film has been already reviewed
	testUseCase.EXPECT().NewReview(newReview, filmID, author, logger).Return(nil, nil)
	request = httptest.NewRequest(http.MethodPost, "/review/1", strings.NewReader(`{"mark": 10,"comment":"very interesting film"}`))
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, author)
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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
	}

	// all is ok
	addedReview := &entity.Review{}
	testUseCase.EXPECT().NewReview(newReview, filmID, author, logger).Return(addedReview, nil)
	request = httptest.NewRequest(http.MethodPost, "/review/1", strings.NewReader(`{"mark": 10,"comment":"very interesting film"}`))
	request = mux.SetURLVars(request, map[string]string{"FILM_ID": "1"})
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, author)
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter = httptest.NewRecorder()
	testHandler.AddReview(respWriter, request.WithContext(ctx))
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
	}
}

func TestDeleteReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewNop().Sugar()

	testUseCase := reviewusecase.NewMockReviewUseCase(ctrl)
	testHandler := NewReviewHandler(testUseCase)

	// bad user in context
	request := httptest.NewRequest(http.MethodDelete, "/review/1", nil)
	request = mux.SetURLVars(request, map[string]string{"REVIEW_ID": "1"})
	ctx := request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, "bad user")
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	respWriter := httptest.NewRecorder()
	testHandler.DeleteReview(respWriter, request.WithContext(ctx))
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

	// bad review id
	author := &entity.User{
		ID:       1,
		Username: "vasya",
	}
	request = httptest.NewRequest(http.MethodDelete, "/review/bad_id", nil)
	request = mux.SetURLVars(request, map[string]string{"REVIEW_ID": "bad_id"})
	respWriter = httptest.NewRecorder()
	ctx = request.Context()
	ctx = context.WithValue(ctx, middleware.MyUserKey, author)
	ctx = context.WithValue(ctx, middleware.MyLoggerKey, logger)
	testHandler.DeleteReview(respWriter, request.WithContext(ctx))
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
}
