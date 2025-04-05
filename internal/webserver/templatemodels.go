package webserver

import (
	"github.com/dev-mackan/gowiki/pkg/models"
)

type PageTmplModel struct {
	bundle      *models.PageBundle
	htmlContent string
}

func NewPageTmplModel(b *models.PageBundle, htmlContent string) *PageTmplModel {
	return &PageTmplModel{
		b,
		htmlContent,
	}
}

type PageRevsTmplModel struct {
	Page      *models.Page
	Revisions *[]models.Revision
}

func NewPageRevsTmplModel(p *models.Page, revs *[]models.Revision) *PageRevsTmplModel {
	return &PageRevsTmplModel{
		p,
		revs,
	}

}
