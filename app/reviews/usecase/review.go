package reviewusecase

import (
	"context"
	"go.uber.org/zap"
	"kinopoisk/app/dto"
	"kinopoisk/app/entity"
	errorapp "kinopoisk/app/errors"
	filmrepo "kinopoisk/app/films/repo/mysql"
	review "kinopoisk/service_review/proto"
)

type ReviewUseCase interface {
	GetFilmReviews(filmID uint64, logger *zap.SugaredLogger) ([]*entity.Review, error)
	NewReview(newReview *dto.ReviewDTO, filmID uint64, user *entity.User, logger *zap.SugaredLogger) (*entity.Review, error)
	DeleteReview(reviewID, userID uint64, logger *zap.SugaredLogger) (bool, error)
	UpdateReview(reviewToUpdate *dto.ReviewDTO, reviewID uint64, user *entity.User, logger *zap.SugaredLogger) (*entity.Review, error)
}

type ReviewGRPCClient struct {
	grpcClient review.ReviewMakerClient
	filmRepo   filmrepo.FilmRepo
}

type loggerKey int

const MyLoggerKey loggerKey = 3

func NewReviewGRPCClient(grpcClient review.ReviewMakerClient, filmRepo filmrepo.FilmRepo) *ReviewGRPCClient {
	return &ReviewGRPCClient{
		grpcClient: grpcClient,
		filmRepo:   filmRepo,
	}
}

func (r *ReviewGRPCClient) GetFilmReviews(filmID uint64, logger *zap.SugaredLogger) ([]*entity.Review, error) {
	reviews, err := r.grpcClient.GetFilmReviews(context.Background(), &review.FilmID{
		ID: filmID,
	})
	if err != nil {
		logger.Errorf("error in getting film reviews: %s", err)
		return nil, err
	}
	reviewsArr := reviews.GetReviews()
	reviewsApp := make([]*entity.Review, len(reviewsArr))
	for i, currentReview := range reviewsArr {
		newReviewApp := getReviewFromGRPCStruct(currentReview)
		reviewsApp[i] = newReviewApp
	}
	return reviewsApp, nil
}

func (r *ReviewGRPCClient) NewReview(newReview *dto.ReviewDTO, filmID uint64, user *entity.User, logger *zap.SugaredLogger) (*entity.Review, error) {
	film, err := r.filmRepo.GetFilmByIDRepo(filmID)
	if err != nil {
		logger.Errorf("error in getting film: %s", err)
		return nil, err
	}
	if film == nil {
		logger.Errorf("no film with id: %d", filmID)
		return nil, errorapp.ErrorNoFilm
	}
	newReviewGRPC, err := r.grpcClient.NewReview(context.Background(), &review.NewReviewData{
		Review: getGRPCReviewFromDTO(newReview),
		FilmID: &review.FilmID{ID: filmID},
		UserID: &review.UserID{ID: user.ID},
	})
	if err != nil {
		logger.Errorf("error in adding film review: %s", err)
		return nil, err
	}
	if newReviewGRPC.ID == nil {
		return nil, nil
	}
	reviewApp := getReviewFromGRPCStruct(newReviewGRPC)
	reviewApp.Author = user

	return reviewApp, nil
}

func (r *ReviewGRPCClient) DeleteReview(reviewID, userID uint64, logger *zap.SugaredLogger) (bool, error) {
	deletedData, err := r.grpcClient.DeleteReview(context.Background(), &review.DeleteReviewData{
		ReviewID: &review.ReviewID{ID: reviewID},
		UserID:   &review.UserID{ID: userID},
	})
	if err != nil {
		logger.Errorf("error in deleting film review: %s", err)
		return false, err
	}
	return deletedData.IsDeleted, nil
}

func (r *ReviewGRPCClient) UpdateReview(reviewToUpdate *dto.ReviewDTO, reviewID uint64, user *entity.User, logger *zap.SugaredLogger) (*entity.Review, error) {
	grpcReview := getGRPCReviewFromDTO(reviewToUpdate)
	grpcReview.ID = &review.ReviewID{
		ID: reviewID,
	}
	updatedReviewGRPC, err := r.grpcClient.UpdateReview(context.Background(), &review.UpdateReviewData{
		Review: grpcReview,
		UserID: &review.UserID{ID: user.ID},
	})
	if err != nil {
		logger.Errorf("error in updating film reviews: %s", err)
		return nil, err
	}
	if updatedReviewGRPC.ID == nil {
		return nil, nil
	}
	updatedReviewApp := getReviewFromGRPCStruct(updatedReviewGRPC)
	updatedReviewApp.Author = user
	return updatedReviewApp, nil
}

func getReviewFromGRPCStruct(reviewGRPC *review.Review) *entity.Review {
	return &entity.Review{
		ID:      reviewGRPC.ID.ID,
		Mark:    reviewGRPC.Mark,
		Comment: reviewGRPC.Comment,
		Author: &entity.User{
			ID:       reviewGRPC.Author.ID.ID,
			Username: reviewGRPC.Author.Username,
		},
	}
}

func getGRPCReviewFromDTO(reviewDTO *dto.ReviewDTO) *review.Review {
	return &review.Review{
		Mark:    reviewDTO.Mark,
		Comment: reviewDTO.Comment,
	}
}
