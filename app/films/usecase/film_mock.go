// Code generated by MockGen. DO NOT EDIT.
// Source: app/films/usecase/film.go

// Package filmusecase is a generated GoMock package.
package filmusecase

import (
	entity "kinopoisk/app/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFilmUseCase is a mock of FilmUseCase interface.
type MockFilmUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockFilmUseCaseMockRecorder
}

// MockFilmUseCaseMockRecorder is the mock recorder for MockFilmUseCase.
type MockFilmUseCaseMockRecorder struct {
	mock *MockFilmUseCase
}

// NewMockFilmUseCase creates a new mock instance.
func NewMockFilmUseCase(ctrl *gomock.Controller) *MockFilmUseCase {
	mock := &MockFilmUseCase{ctrl: ctrl}
	mock.recorder = &MockFilmUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilmUseCase) EXPECT() *MockFilmUseCaseMockRecorder {
	return m.recorder
}

// AddFavouriteFilm mocks base method.
func (m *MockFilmUseCase) AddFavouriteFilm(userID, filmID uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFavouriteFilm", userID, filmID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddFavouriteFilm indicates an expected call of AddFavouriteFilm.
func (mr *MockFilmUseCaseMockRecorder) AddFavouriteFilm(userID, filmID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFavouriteFilm", reflect.TypeOf((*MockFilmUseCase)(nil).AddFavouriteFilm), userID, filmID)
}

// DeleteFavouriteFilm mocks base method.
func (m *MockFilmUseCase) DeleteFavouriteFilm(userID, filmID uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFavouriteFilm", userID, filmID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFavouriteFilm indicates an expected call of DeleteFavouriteFilm.
func (mr *MockFilmUseCaseMockRecorder) DeleteFavouriteFilm(userID, filmID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFavouriteFilm", reflect.TypeOf((*MockFilmUseCase)(nil).DeleteFavouriteFilm), userID, filmID)
}

// GetFavouriteFilms mocks base method.
func (m *MockFilmUseCase) GetFavouriteFilms(userID uint64) ([]*entity.Film, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavouriteFilms", userID)
	ret0, _ := ret[0].([]*entity.Film)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavouriteFilms indicates an expected call of GetFavouriteFilms.
func (mr *MockFilmUseCaseMockRecorder) GetFavouriteFilms(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavouriteFilms", reflect.TypeOf((*MockFilmUseCase)(nil).GetFavouriteFilms), userID)
}

// GetFilmActors mocks base method.
func (m *MockFilmUseCase) GetFilmActors(filmID uint64) ([]*entity.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmActors", filmID)
	ret0, _ := ret[0].([]*entity.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmActors indicates an expected call of GetFilmActors.
func (mr *MockFilmUseCaseMockRecorder) GetFilmActors(filmID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmActors", reflect.TypeOf((*MockFilmUseCase)(nil).GetFilmActors), filmID)
}

// GetFilmByID mocks base method.
func (m *MockFilmUseCase) GetFilmByID(filmID uint64) (*entity.Film, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmByID", filmID)
	ret0, _ := ret[0].(*entity.Film)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmByID indicates an expected call of GetFilmByID.
func (mr *MockFilmUseCaseMockRecorder) GetFilmByID(filmID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmByID", reflect.TypeOf((*MockFilmUseCase)(nil).GetFilmByID), filmID)
}

// GetFilmGenres mocks base method.
func (m *MockFilmUseCase) GetFilmGenres(filmID uint64) ([]*entity.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmGenres", filmID)
	ret0, _ := ret[0].([]*entity.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmGenres indicates an expected call of GetFilmGenres.
func (mr *MockFilmUseCaseMockRecorder) GetFilmGenres(filmID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmGenres", reflect.TypeOf((*MockFilmUseCase)(nil).GetFilmGenres), filmID)
}

// GetFilms mocks base method.
func (m *MockFilmUseCase) GetFilms(genre, country, producer string) ([]*entity.Film, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilms", genre, country, producer)
	ret0, _ := ret[0].([]*entity.Film)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilms indicates an expected call of GetFilms.
func (mr *MockFilmUseCaseMockRecorder) GetFilms(genre, country, producer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilms", reflect.TypeOf((*MockFilmUseCase)(nil).GetFilms), genre, country, producer)
}

// GetFilmsByActor mocks base method.
func (m *MockFilmUseCase) GetFilmsByActor(ID uint64) ([]*entity.Film, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmsByActor", ID)
	ret0, _ := ret[0].([]*entity.Film)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmsByActor indicates an expected call of GetFilmsByActor.
func (mr *MockFilmUseCaseMockRecorder) GetFilmsByActor(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmsByActor", reflect.TypeOf((*MockFilmUseCase)(nil).GetFilmsByActor), ID)
}

// GetSoonFilms mocks base method.
func (m *MockFilmUseCase) GetSoonFilms() ([]*entity.Film, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSoonFilms")
	ret0, _ := ret[0].([]*entity.Film)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSoonFilms indicates an expected call of GetSoonFilms.
func (mr *MockFilmUseCaseMockRecorder) GetSoonFilms() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSoonFilms", reflect.TypeOf((*MockFilmUseCase)(nil).GetSoonFilms))
}
