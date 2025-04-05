package sqliterepo

import (
	"context"
	"database/sql"
	"github.com/dev-mackan/gowiki/pkg/models"
)

type SqliteTextRepository struct {
	db *sql.DB
}

func NewSqliteTextRepository(db *sql.DB) *SqliteTextRepository {
	return &SqliteTextRepository{
		db,
	}
}

func (r *SqliteTextRepository) Create(ctx context.Context, text *models.Text) error {
	return nil
}

func (r *SqliteTextRepository) Update(ctx context.Context, text *models.Text) error {
	return nil
}
func (r *SqliteTextRepository) GetByID(ctx context.Context, textId uint) (*models.Text, error) {
	query := `SELECT text_id, content, created_at FROM Text WHERE text_id = ?`
	var text models.Text
	err := r.db.QueryRowContext(ctx, query, textId).Scan(&text.TextId, &text.Content, &text.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &text, nil
}
