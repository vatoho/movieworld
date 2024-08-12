package ratingserviceusecase

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	ratingservicerepo "kinopoisk/service_rating/repo/mysql"
)

type ChangeRatingInfo struct {
	ChangeType string
	ReviewID   uint64
	OldMark    uint32
	NewMark    uint32
	FilmID     uint64
}

type RatingChanger interface {
	ChangeRating(tasks <-chan amqp.Delivery)
}

type RatingChangerApp struct {
	logger           *zap.SugaredLogger
	changeRatingRepo ratingservicerepo.RatingChangerDB
}

func NewRatingChangerApp(logger *zap.SugaredLogger, changeRatingRepo ratingservicerepo.RatingChangerDB) *RatingChangerApp {
	return &RatingChangerApp{
		logger:           logger,
		changeRatingRepo: changeRatingRepo,
	}
}

func (r *RatingChangerApp) ChangeRating(tasks <-chan amqp.Delivery) {
	for taskItem := range tasks {
		r.logger.Infof("incoming task %+v\n", taskItem)
		task := &ChangeRatingInfo{}
		err := json.Unmarshal(taskItem.Body, task)
		if err != nil {
			r.logger.Errorf("cant unpack json: %s", err)
			err = taskItem.Ack(false)
			if err != nil {
				r.logger.Errorf("error with ack: %s", err)
			}
			continue
		}

		switch task.ChangeType {
		case "Add":
			err = r.changeRatingRepo.ChangeRatingAddReview(task.NewMark, task.ReviewID)
		case "Update":
			err = r.changeRatingRepo.ChangeRatingAfterUpdateReview(task.OldMark, task.NewMark, task.ReviewID)
		case "Delete":
			err = r.changeRatingRepo.ChangeRatingAfterDeleteReview(task.OldMark, task.FilmID)
		default:
			r.logger.Errorf("unknown task: %s", task.ChangeType)
		}
		if err != nil {
			err = taskItem.Nack(false, true) // Возвращаем сообщение в очередь в случае ошибки
			if err != nil {
				r.logger.Errorf("error with nack: %s", err)
			}
			continue
		}
		err = taskItem.Ack(false)
		if err != nil {
			r.logger.Errorf("error with ack: %s", err)

		}
	}
}
