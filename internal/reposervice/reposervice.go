package reposervice

import (
	"context"
	"database/sql"
	"github.com/dev-mackan/gowiki/internal/repos"
	"github.com/dev-mackan/gowiki/pkg/models"
	"github.com/dev-mackan/gowiki/pkg/utils"
)

// Interface for the database
// Takes care of bundling items when necessary
type RepoService struct {
	repo *repos.Repository
}

func NewRepoService(repo *repos.Repository) *RepoService {
	return &RepoService{
		repo: repo,
	}
}

func (rs *RepoService) GetPageByTitle(ctx context.Context, title string) (*models.Page, error) {
	pageId, err := rs.getPageIdByTitle(ctx, title)
	if err != nil {
		return nil, handleErr(err)
	}
	page, err := rs.getPageById(ctx, pageId)
	if err != nil {
		return nil, handleErr(err)
	}
	return page, nil
}

func (rs *RepoService) GetBundledPageByTitle(ctx context.Context, title string) (*models.PageBundle, error) {
	//TODO: Use bundle db instead so its a full transaction?
	pageId, err := rs.getPageIdByTitle(ctx, title)
	if err != nil {
		return nil, handleErr(err)
	}

	page, err := rs.getPageById(ctx, pageId)
	if err != nil {
		return nil, handleErr(err)
	}

	rev, err := rs.getRevisionById(ctx, page.LatestRev)
	if err != nil {
		return nil, handleErr(err)
	}

	text, err := rs.getTextById(ctx, rev.TextId)
	if err != nil {
		return nil, handleErr(err)
	}

	return rs.buildPageBundle(page, rev, text), nil
}

func (rs *RepoService) GetBundledPageWithRev(ctx context.Context, title string, revId uint) (*models.PageBundle, error) {
	//TODO: Use bundle db instead so its a full transaction?
	pageId, err := rs.getPageIdByTitle(ctx, title)
	if err != nil {
		return nil, handleErr(err)
	}

	page, err := rs.getPageById(ctx, pageId)
	if err != nil {
		return nil, handleErr(err)
	}

	rev, err := rs.getRevisionById(ctx, revId)
	if err != nil {
		return nil, handleErr(err)
	}

	text, err := rs.getTextById(ctx, rev.TextId)
	if err != nil {
		return nil, handleErr(err)
	}

	return rs.buildPageBundle(page, rev, text), nil
}

func (rs *RepoService) DeleteBundledPage(ctx context.Context, pageId uint) error {
	err := rs.repo.Bundled.DeleteBundle(ctx, pageId)
	if err != nil {
		return handleErr(err)
	}
	return nil
}

func (rs *RepoService) CreateBundledPage(ctx context.Context, title string, content string) error {
	title = utils.SanitizeTitle(title)
	err := rs.repo.Bundled.NewPageBundle(ctx, title, content)
	if err != nil {
		return handleErr(err)
	}
	return nil
}

func (rs *RepoService) UpdateBundledPage(ctx context.Context, pageId uint, title string, content string) error {
	title = utils.SanitizeTitle(title)
	err := rs.repo.Bundled.UpdateBundledPage(ctx, pageId, title, content)
	if err != nil {
		return handleErr(err)
	}
	return err
}
func (rs *RepoService) NewPageRev(ctx context.Context, pageId uint, content string) error {
	err := rs.repo.Bundled.UpdateBundledPageContent(ctx, pageId, content)
	if err != nil {
		return handleErr(err)
	}
	return nil
}

func (rs *RepoService) GetPageRevs(ctx context.Context, title string) (*[]*models.Revision, error) {
	pageId, err := rs.getPageIdByTitle(ctx, title)
	if err != nil {
		return nil, handleErr(err)
	}
	revs, err := rs.getRevsByPageID(ctx, pageId)
	if err != nil {
		return nil, handleErr(err)
	}
	return revs, nil
}

func (rs *RepoService) GetPages(ctx context.Context) (*[]*models.Page, error) {
	pages, err := rs.repo.Page.GetAll(ctx)
	if err != nil {
		return nil, handleErr(err)
	}
	return pages, nil
}

func (rs *RepoService) GetTextByRevID(ctx context.Context, revID uint) (*models.Text, error) {
	rev, err := rs.repo.Revision.GetByID(ctx, revID)
	if err != nil {
		return nil, handleErr(err)
	}
	text, err := rs.repo.Text.GetByID(ctx, rev.TextId)
	if err != nil {
		return nil, handleErr(err)
	}
	return text, nil
}

func (rs *RepoService) UpdatePageTitle(ctx context.Context, pageId uint, title string) error {
	title = utils.SanitizeTitle(title)
	err := rs.repo.Page.UpdateTitle(ctx, pageId, title)
	if err != nil {
		return handleErr(err)
	}
	return nil
}

func (rs *RepoService) getPageIdByTitle(ctx context.Context, title string) (uint, error) {
	return rs.repo.Page.GetIDByTitle(ctx, title)
}

func (rs *RepoService) getPageById(ctx context.Context, pageId uint) (*models.Page, error) {
	return rs.repo.Page.GetByID(ctx, pageId)
}

func (rs *RepoService) getRevisionById(ctx context.Context, revId uint) (*models.Revision, error) {
	return rs.repo.Revision.GetByID(ctx, revId)
}

func (rs *RepoService) getTextById(ctx context.Context, textId uint) (*models.Text, error) {
	return rs.repo.Text.GetByID(ctx, textId)
}

func (rs *RepoService) getRevsByPageID(ctx context.Context, revId uint) (*[]*models.Revision, error) {
	return rs.repo.Revision.GetAllByPageID(ctx, revId)
}

func (rs *RepoService) buildPageBundle(page *models.Page, rev *models.Revision, text *models.Text) *models.PageBundle {
	return &models.PageBundle{
		Page:     page,
		Revision: rev,
		Text:     text,
	}
}

func handleErr(err error) error {
	if err == sql.ErrNoRows {
		return &NotFoundError{err: err.Error()}
	}
	return &DatabaseError{err: err.Error()}
}
