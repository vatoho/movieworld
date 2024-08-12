package ratingservicerepo

import (
	"database/sql"
	"go.uber.org/zap"
)

type RatingChangerDB interface {
	ChangeRatingAfterDeleteReview(oldMark uint32, reviewID uint64) error
	ChangeRatingAfterUpdateReview(oldMark, newMark uint32, reviewID uint64) error
	ChangeRatingAddReview(newMark uint32, reviewID uint64) error
}

type RatingChangerMySQL struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewRatingChangerMySQL(db *sql.DB, logger *zap.SugaredLogger) *RatingChangerMySQL {
	return &RatingChangerMySQL{
		db:     db,
		logger: logger,
	}
}

func (r *RatingChangerMySQL) ChangeRatingAfterDeleteReview(oldMark uint32, filmID uint64) error {
	_, err := r.db.Exec(
		`UPDATE films 
                SET 
                   sum_mark = sum_mark - ?,
                num_of_marks = num_of_marks - 1,
                rating = CASE 
                              WHEN num_of_marks > 1 THEN (sum_mark - ?) / (num_of_marks - 1)
                            ELSE 0
                          END
                WHERE id = ?`,
		oldMark,
		oldMark,
		filmID,
	)
	if err != nil {
		r.logger.Errorf("error in changing rating after delete review")
	}
	return err
}

func (r *RatingChangerMySQL) ChangeRatingAfterUpdateReview(oldMark, newMark uint32, reviewID uint64) error {
	_, err := r.db.Exec(
		"UPDATE films SET sum_mark = sum_mark + ? - ?, rating = (sum_mark + ? - ?) / num_of_marks WHERE id in (SELECT film_id from reviews WHERE id = ?)",
		newMark,
		oldMark,
		newMark,
		oldMark,
		reviewID,
	)
	if err != nil {
		r.logger.Errorf("error in changing rating after update review: %s", err)
	}
	return err
}

func (r *RatingChangerMySQL) ChangeRatingAddReview(newMark uint32, reviewID uint64) error {
	_, err := r.db.Exec(
		"UPDATE films SET sum_mark = sum_mark + ?, num_of_marks = num_of_marks + 1, rating = (sum_mark + ?) / (num_of_marks + 1) WHERE id in (SELECT film_id from reviews WHERE id = ?)",
		newMark,
		newMark,
		reviewID,
	)
	if err != nil {
		r.logger.Errorf("error in changing rating after add review")
	}
	return err
}
