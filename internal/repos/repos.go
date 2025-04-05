package repos

import (
	"context"
	"database/sql"

	"github.com/dev-mackan/gowiki/internal/repos/sqliterepo"
	"github.com/dev-mackan/gowiki/pkg/models"
)

type Repository struct {
	Bundled interface {
		NewPageBundle(context.Context, string, string) error
		UpdateBundledPage(context.Context, uint, string, string) error
		UpdateBundledPageContent(context.Context, uint, string) error
		DeleteBundle(context.Context, uint) error
	}
	Page interface {
		GetIDByTitle(context.Context, string) (uint, error)
		GetByID(context.Context, uint) (*models.Page, error)
		Create(context.Context, *models.Page) error
		Update(context.Context, *models.Page) error
		UpdateTitle(context.Context, uint, string) error
		GetAll(context.Context) (*[]*models.Page, error)
	}
	Revision interface {
		GetByID(context.Context, uint) (*models.Revision, error)
		Create(context.Context, *models.Revision) error
		Update(context.Context, *models.Revision) error
		GetAllByPageID(context.Context, uint) (*[]*models.Revision, error)
	}
	Text interface {
		GetByID(context.Context, uint) (*models.Text, error)
		Create(context.Context, *models.Text) error
		Update(context.Context, *models.Text) error
	}
}

func NewSqlRepository(db *sql.DB) *Repository {
	return &Repository{
		Bundled:  sqliterepo.NewSqliteBundledRepository(db),
		Page:     sqliterepo.NewSqlitePageRepository(db),
		Revision: sqliterepo.NewSqliteRevisionRepository(db),
		Text:     sqliterepo.NewSqliteTextRepository(db),
	}
}
