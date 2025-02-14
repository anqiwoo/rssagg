package database

import (
	"context"
)

func (q *Queries) CustomDeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteFeedFollow, arg.ID, arg.UserID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
