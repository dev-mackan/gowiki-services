package sqliterepo

import (
	"context"
	"database/sql"

	"github.com/dev-mackan/gowiki/pkg/models"
)

type SqlitePageRepository struct {
	db *sql.DB
}

func NewSqlitePageRepository(db *sql.DB) *SqlitePageRepository {
	return &SqlitePageRepository{
		db,
	}
}

func (r *SqlitePageRepository) GetIDByTitle(ctx context.Context, title string) (uint, error) {
	query := `SELECT page_id FROM Page WHERE UPPER(title)=UPPER(?)`
	var id uint
	err := r.db.QueryRowContext(ctx, query, title).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlitePageRepository) Create(ctx context.Context, page *models.Page) error {
	//query := `INSERT INTO Page (title) VALUES (?) RETURNING page_id, created_at`
	return nil
}
func (r *SqlitePageRepository) Update(ctx context.Context, page *models.Page) error {
	//query := `UPDATE Page SET latest_rev = ?, title = ? WHERE page_id = ?`
	return nil
}
func (r *SqlitePageRepository) UpdateTitle(ctx context.Context, pageId uint, title string) error {
	query := `UPDATE Page SET title = ? WHERE page_id = ?`
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, query, title, pageId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlitePageRepository) GetAll(ctx context.Context) (*[]*models.Page, error) {
	query := `SELECT page_id, title, latest_rev, created_at FROM Page`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	pages := make([]*models.Page, 0)
	for rows.Next() {
		var page models.Page
		err = rows.Scan(&page.PageId, &page.Title, &page.LatestRev, &page.CreatedAt)
		if err != nil {
			continue
		}
		pages = append(pages, &page)
	}
	return &pages, nil
}

func (r *SqlitePageRepository) GetByID(ctx context.Context, pageId uint) (*models.Page, error) {
	query := `SELECT page_id, title, latest_rev, created_at FROM Page WHERE page_id=?`
	var page models.Page
	err := r.db.QueryRowContext(ctx, query, pageId).Scan(&page.PageId, &page.Title, &page.LatestRev, &page.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &page, nil
}
