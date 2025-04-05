package sqliterepo

import (
	"context"
	"database/sql"
)

type SqliteBundledRepository struct {
	db *sql.DB
}

func NewSqliteBundledRepository(db *sql.DB) *SqliteBundledRepository {
	return &SqliteBundledRepository{
		db,
	}
}

func (s *SqliteBundledRepository) NewPageBundle(ctx context.Context, title string, content string) error {
	textQuery := `INSERT INTO Text (content) VALUES (?) RETURNING text_id`
	pageQuery := `INSERT INTO Page (title, latest_rev) VALUES (?,0) RETURNING page_id`
	revQuery := `INSERT INTO Revision (text_id,page_id) VALUES (?,?) RETURNING rev_id`
	pageUpdQuery := `UPDATE Page SET latest_rev=? WHERE page_id=?`

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var textId uint
	err = tx.QueryRowContext(ctx, textQuery, content).Scan(&textId)
	if err != nil {
		return err
	}

	var pageId uint
	err = tx.QueryRowContext(ctx, pageQuery, title).Scan(&pageId)
	if err != nil {
		return err
	}

	var revId uint
	err = tx.QueryRowContext(ctx, revQuery, textId, pageId).Scan(&revId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, pageUpdQuery, revId, pageId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteBundledRepository) UpdateBundledPage(ctx context.Context, pageId uint, title string, content string) error {
	textQuery := `INSERT INTO Text (content) VALUES (?) RETURNING text_id`
	revQuery := `INSERT INTO Revision (text_id,page_id) VALUES (?,?) RETURNING rev_id`
	pageUpdQuery := `UPDATE Page SET latest_rev=?, title=? WHERE page_id=?`

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var textId uint
	err = tx.QueryRowContext(ctx, textQuery, content).Scan(&textId)
	if err != nil {
		return err
	}

	var revId uint
	err = tx.QueryRowContext(ctx, revQuery, textId, pageId).Scan(&revId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, pageUpdQuery, revId, title, pageId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteBundledRepository) UpdateBundledPageContent(ctx context.Context, pageId uint, content string) error {
	textQuery := `INSERT INTO Text (content) VALUES (?) RETURNING text_id`
	revQuery := `INSERT INTO Revision (text_id,page_id) VALUES (?,?) RETURNING rev_id`
	pageUpdQuery := `UPDATE Page SET latest_rev=?  WHERE page_id=?`

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var textId uint
	err = tx.QueryRowContext(ctx, textQuery, content).Scan(&textId)
	if err != nil {
		return err
	}

	var revId uint
	err = tx.QueryRowContext(ctx, revQuery, textId, pageId).Scan(&revId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, pageUpdQuery, revId, pageId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteBundledRepository) DeleteBundle(ctx context.Context, pageId uint) error {
	pageQuery := `DELETE FROM Page WHERE page_id = ?`
	revQuery := `DELETE FROM Revision WHERE page_id = ?`
	textQuery := `DELETE FROM Text WHERE text_id IN (SELECT text_id FROM Revision WHERE page_id=?)`
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, textQuery, pageId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, revQuery, pageId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, pageQuery, pageId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
