package sqliterepo

import (
	"context"
	"database/sql"

	"github.com/dev-mackan/gowiki/pkg/models"
)

type SqliteRevisionRepository struct {
	db *sql.DB
}

func NewSqliteRevisionRepository(db *sql.DB) *SqliteRevisionRepository {
	return &SqliteRevisionRepository{
		db,
	}
}

func (r *SqliteRevisionRepository) Create(ctx context.Context, rev *models.Revision) error {
	return nil
}

func (r *SqliteRevisionRepository) Update(ctx context.Context, rev *models.Revision) error {
	return nil
}
func (r *SqliteRevisionRepository) GetByID(ctx context.Context, revId uint) (*models.Revision, error) {
	query := `SELECT rev_id, page_id, text_id, created_at FROM Revision WHERE rev_id = ?`
	var rev models.Revision
	err := r.db.QueryRowContext(ctx, query, revId).Scan(&rev.RevId, &rev.PageId, &rev.TextId, &rev.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rev, nil
}
func (r *SqliteRevisionRepository) GetAllByPageID(ctx context.Context, pageId uint) (*[]*models.Revision, error) {
	query := `SELECT rev_id, page_id, text_id, created_at FROM Revision WHERE page_id = ?`
	rows, err := r.db.QueryContext(ctx, query, pageId)
	if err != nil {
		return nil, err
	}
	revs := make([]*models.Revision, 0)
	for rows.Next() {
		var rev models.Revision
		err = rows.Scan(&rev.RevId, &rev.PageId, &rev.TextId, &rev.CreatedAt)
		if err != nil {
			continue
		}
		revs = append(revs, &rev)
	}
	return &revs, nil
}
