package db

import (
	"context"
	"database/sql"
	"fmt"
)


// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	SubmitAttemptTx(ctx context.Context, arg CreateAttemptsParams) (SubmitAttemptTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUsersParams) (CreateUserTxResult, error)
	DeleteDictationTx(ctx context.Context, dictationID int64) error
}

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// SubmitAttemptTxResult contains the result of the SubmitAttemptTx operation
type SubmitAttemptTxResult struct {
	Attempt            Attempt
	PerformanceSummary PerformanceSummary
}

// SubmitAttemptTx performs the necessary steps to submit an attempt and update performance summary
func (store *SQLStore) SubmitAttemptTx(ctx context.Context, arg CreateAttemptsParams) (SubmitAttemptTxResult, error) {
	var result SubmitAttemptTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1. Create the Attempt
		result.Attempt, err = q.CreateAttempts(ctx, arg)
		if err != nil {
			return err
		}

		// 2. Get Performance Summary
		summary, err := q.GetPerformanceSummaryByUserAndDictation(ctx, GetPerformanceSummaryByUserAndDictationParams{
			UserID:      arg.UserID,
			DictationID: arg.DictationID,
		})

		// 3. Update or Create Summary
		if err == sql.ErrNoRows {
			// Create new summary
			result.PerformanceSummary, err = q.CreatePerformanceSummary(ctx, CreatePerformanceSummaryParams{
				UserID:          arg.UserID,
				DictationID:     arg.DictationID,
				TotalAttempts:   sql.NullInt32{Int32: 1, Valid: true},
				BestAccuracy:    arg.Accuracy,
				AverageAccuracy: arg.Accuracy,
				AverageTime:     arg.TimeSpent,
				LastAttemptAt:   sql.NullTime{Time: result.Attempt.CreatedAt.Time, Valid: true},
			})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// Update existing summary
			newTotalAttempts := summary.TotalAttempts.Int32 + 1
			newAvgAccuracy := ((summary.AverageAccuracy.Float64 * float64(summary.TotalAttempts.Int32)) + arg.Accuracy.Float64) / float64(newTotalAttempts)
			
			// Handle Average Time (if present)
			currentAvgTime := summary.AverageTime.Float64
			newTime := arg.TimeSpent.Float64
			newAvgTime := ((currentAvgTime * float64(summary.TotalAttempts.Int32)) + newTime) / float64(newTotalAttempts)

			// Determine Best Accuracy
			newBestAccuracy := summary.BestAccuracy.Float64
			if arg.Accuracy.Float64 > newBestAccuracy {
				newBestAccuracy = arg.Accuracy.Float64
			}

			result.PerformanceSummary, err = q.UpdatePerformanceSummary(ctx, UpdatePerformanceSummaryParams{
				ID:              summary.ID,
				TotalAttempts:   sql.NullInt32{Int32: newTotalAttempts, Valid: true},
				BestAccuracy:    sql.NullFloat64{Float64: newBestAccuracy, Valid: true},
				AverageAccuracy: sql.NullFloat64{Float64: newAvgAccuracy, Valid: true},
				AverageTime:     sql.NullFloat64{Float64: newAvgTime, Valid: true},
				LastAttemptAt:   sql.NullTime{Time: result.Attempt.CreatedAt.Time, Valid: true},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

// CreateUserTxResult contains the result of the CreateUserTx operation
type CreateUserTxResult struct {
	User    User
	Setting Setting
}

// CreateUserTx performs the necessary steps to create a user and default settings
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUsersParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1. Create User
		result.User, err = q.CreateUsers(ctx, arg)
		if err != nil {
			return err
		}

		// 2. Create Default Settings
		result.Setting, err = q.CreateSetting(ctx, CreateSettingParams{
			UserID:                 sql.NullInt64{Int64: result.User.ID, Valid: true},
			DefaultVoice:           sql.NullString{String: "en-US-Neural2-F", Valid: true}, 
			DefaultSpeed:           sql.NullFloat64{Float64: 1.0, Valid: true},
			HighlightColorGrammar:  sql.NullString{String: "#FFA500", Valid: true}, // Orange
			HighlightColorSpelling: sql.NullString{String: "#FF0000", Valid: true}, // Red
			HighlightColorCase:     sql.NullString{String: "#FFFF00", Valid: true}, // Yellow
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// DeleteDictationTx deletes a dictation and all associated data (cascading delete)
func (store *SQLStore) DeleteDictationTx(ctx context.Context, dictationID int64) error {
	return store.execTx(ctx, func(q *Queries) error {
		// 1. Delete Performance Summaries for this Dictation
		// Note: Using direct execution because generated query DeletePerformanceSummaryByDictation is missing
		_, err := store.db.ExecContext(ctx, "DELETE FROM performance_summary WHERE dictation_id = $1", dictationID)
		if err != nil {
			return err
		}

		// 2. Delete Attempts for this Dictation
		// Note: Using direct execution because generated DeleteAttemptsByDictation filters by user_id too
		// We want to delete ALL attempts for this dictation regardless of user
		_, err = store.db.ExecContext(ctx, "DELETE FROM attempts WHERE dictation_id = $1", dictationID)
		if err != nil {
			return err
		}

		// 3. Delete the Dictation itself
		// Note: Generated DeleteDictations deletes by title, we want ID for safety
		_, err = store.db.ExecContext(ctx, "DELETE FROM dictations WHERE id = $1", dictationID)
		if err != nil {
			return err
		}

		return nil
	})
}




