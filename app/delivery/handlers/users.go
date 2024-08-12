package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"kinopoisk/app/delivery"
	"kinopoisk/app/dto"
	"kinopoisk/app/entity"
	errorapp "kinopoisk/app/errors"
	"kinopoisk/app/middleware"
	userusecase "kinopoisk/app/users/usecase"
	"log"
	"net/http"
)

type UserHandler struct {
	UserUseCases userusecase.UserUseCase
}

func NewUserHandler(userUseCases userusecase.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCases: userUseCases,
	}
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := middleware.GetLoggerFromContext(ctx)
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		middleware.WriteNoLoggerResponse(w)
	}
	userFromLoginForm, err := checkRequestFormat(logger, w, r)
	if err != nil || userFromLoginForm == nil {
		return
	}
	loggedInUser, err := uh.UserUseCases.Login(userFromLoginForm.Username, userFromLoginForm.Password, logger)

	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in getting user by login and password: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	if loggedInUser == nil {
		delivery.WriteResponse(logger, w, []byte(`{"message": "bad username or password"}`), http.StatusUnauthorized)
		return
	}
	uh.HandleGetToken(w, loggedInUser, logger)

}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := middleware.GetLoggerFromContext(ctx)
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		middleware.WriteNoLoggerResponse(w)
	}
	userFromLoginForm, err := checkRequestFormat(logger, w, r)
	if err != nil || userFromLoginForm == nil {
		return
	}
	newUser, err := uh.UserUseCases.Register(userFromLoginForm.Username, userFromLoginForm.Password, logger)

	if errors.Is(err, errorapp.ErrorUserExists) {
		delivery.WriteResponse(logger, w, []byte(`{"message": "user already exists"}`), http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		errText := fmt.Sprintf(`{"message": "unknown error occured in register: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	uh.HandleGetToken(w, newUser, logger)
}

func (uh *UserHandler) HandleGetToken(w http.ResponseWriter, newUser *entity.User, logger *zap.SugaredLogger) {
	token, err := uh.UserUseCases.CreateSession(newUser, logger)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in session creation: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	resp := dto.AuthResponseDTO{
		Token: token,
	}
	tokenJSON, err := json.Marshal(&resp)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in coding response: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	logger.Infof("new token: %s", token)
	delivery.WriteResponse(logger, w, tokenJSON, http.StatusOK)
}

func checkRequestFormat(logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) (*dto.AuthRequestDTO, error) {
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in reading request body: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusUnauthorized)
		return nil, err
	}
	userFromLoginForm := &dto.AuthRequestDTO{}
	err = json.Unmarshal(rBody, userFromLoginForm)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in decoding user: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusUnauthorized)
		return nil, err
	}
	if validationErrors := userFromLoginForm.Validate(); len(validationErrors) != 0 {
		errorsJSON, err := json.Marshal(validationErrors)
		if err != nil {
			errText := fmt.Sprintf(`{"message": "error in decoding validation errors: %s"}`, err)
			delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
			return nil, err
		}
		logger.Errorf("login form did not pass validation: %s", err)
		delivery.WriteResponse(logger, w, errorsJSON, http.StatusUnauthorized)
		return nil, err
	}
	return userFromLoginForm, nil
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := middleware.GetLoggerFromContext(ctx)
	if err != nil {
		log.Printf("can not get logger from context: %s", err)
		middleware.WriteNoLoggerResponse(w)
	}
	token, ok := ctx.Value(middleware.MyTokenKey).(string)
	if !ok {
		delivery.WriteResponse(logger, w, []byte(`{"message": "can not cast context value to user"}`), http.StatusInternalServerError)
		return
	}
	isDeleted, err := uh.UserUseCases.DeleteSession(token, logger)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in logging out: %s"}`, err)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	if !isDeleted {
		errText := fmt.Sprintf(`{"message": "no session with token: %s"}`, token)
		delivery.WriteResponse(logger, w, []byte(errText), http.StatusNotFound)
		return
	}
	message := `{"result":"success"}`
	delivery.WriteResponse(logger, w, []byte(message), http.StatusOK)
}
