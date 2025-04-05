package apiserver

import (
	"encoding/json"
	"github.com/dev-mackan/gowiki/internal/messages"
	"github.com/dev-mackan/gowiki/internal/middleware"
	"github.com/dev-mackan/gowiki/internal/reposervice"
	"log"
	"net/http"
)

type APIServer struct {
	listenAddr string
	repo       *reposervice.RepoService
}

func NewAPIServer(listenAddr string, repo *reposervice.RepoService) *APIServer {
	return &APIServer{
		listenAddr,
		repo,
	}
}

func (s *APIServer) mount() *http.ServeMux {
	router := http.NewServeMux()
	logger := middleware.NewLoggerMiddleware("API-SERVER")
	// POST
	router.Handle("POST /api/v1/bundled/new", logger(makeApiHandlerFunc(s.createBundledPage)))
	// PUT
	router.Handle("PUT /api/v1/pages/{page_id}/update/title", logger(makeApiHandlerFunc(s.updatePageTitle)))
	router.Handle("PUT /api/v1/pages/{page_id}/update/content", logger(makeApiHandlerFunc(s.updatePageContent)))
	// DELETE
	//TODO: Add page id to url
	router.Handle("DELETE /api/v1/bundled/delete", logger(makeApiHandlerFunc(s.deleteBundledPage)))
	// GET
	router.Handle("GET /api/v1/pages", logger(makeApiHandlerFunc(s.getPages)))
	router.Handle("GET /api/v1/pages/{page_title}", logger(makeApiHandlerFunc(s.getPage)))
	router.Handle("GET /api/v1/bundled/{page_title}", logger(makeApiHandlerFunc(s.getBundledPage)))
	router.Handle("GET /api/v1/pages/{page_title}/revisions", logger(makeApiHandlerFunc(s.getPageRevs)))
	router.Handle("GET /api/v1/bundled/{page_title}/revisions/{rev_id}", logger(makeApiHandlerFunc(s.getBundledPageWithRev)))
	router.Handle("GET /api/v1/revisions/{rev_id}/text/raw", logger(makeApiHandlerFunc(s.getRawTextForPageWithRev)))
	return router
}

func (s *APIServer) Run() error {
	mux := s.mount()
	log.Println("GOWIKI-API listening on: ", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, mux)
}

func encodeJSON[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return ReplyErr(err)
	}
	return nil
}

func decodeJSON[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, ReceiveErr(err)
	}
	return v, nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeApiHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Println(err)
			status := errToStatusCode(err)
			encodeJSON(w, r, status, err)
		}
	}
}

func (s *APIServer) getPage(w http.ResponseWriter, r *http.Request) error {
	title := r.PathValue("page_title")
	ctx := r.Context()
	page, err := s.repo.GetPageByTitle(ctx, title)
	if err != nil {
		return parseDbErr(err)
	}
	return encodeJSON(w, r, 200, page)
}

func (s *APIServer) getBundledPage(w http.ResponseWriter, r *http.Request) error {
	title := r.PathValue("page_title")
	ctx := r.Context()
	bundle, err := s.repo.GetBundledPageByTitle(ctx, title)
	if err != nil {
		return parseDbErr(err)
	}
	return encodeJSON(w, r, 200, bundle)
}

func (s *APIServer) createBundledPage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	br, err := decodeJSON[messages.NewBundleRequest](r)
	if err != nil {
		return BadRequestErr(err)
	}
	err = s.repo.CreateBundledPage(ctx, br.PageTitle, br.TextContent)
	if err != nil {
		return parseDbErr(err)
	}
	return encodeJSON(w, r, 201, &messages.Empty{})
}

func (s *APIServer) deleteBundledPage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	dr, err := decodeJSON[messages.DeletePageRequest](r)
	if err != nil {
		return BadRequestErr(err)
	}
	err = s.repo.DeleteBundledPage(ctx, dr.PageId)
	if err != nil {
		return parseDbErr(err)
	}
	return encodeJSON(w, r, 200, &messages.Empty{})
}

func (s *APIServer) updatePageBundled(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	br, err := decodeJSON[messages.UpdateBundleRequest](r)
	if err != nil {
		return BadRequestErr(err)
	}
	if br.TextContent == "" {
	} else {
		err = s.repo.UpdateBundledPage(ctx, br.PageId, br.PageTitle, br.TextContent)
		if err != nil {
			return parseDbErr(err)
		}
	}
	m := messages.Empty{}
	return encodeJSON(w, r, 200, m)
}

func (s *APIServer) updatePageTitle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	rq, err := decodeJSON[messages.UpdatePageTitleRequest](r)
	if err != nil {
		return BadRequestErr(err)
	}
	if rq.PageTitle == "" {
		return BadRequestErr(err)
	}
	err = s.repo.UpdatePageTitle(ctx, rq.PageId, rq.PageTitle)
	if err != nil {
		return parseDbErr(err)
	}
	m := messages.Empty{}
	return encodeJSON(w, r, 200, m)
}

func (s *APIServer) updatePageContent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	rq, err := decodeJSON[messages.UpdatePageContentRequest](r)
	if err != nil {
		return BadRequestErr(err)
	}
	err = s.repo.NewPageRev(ctx, rq.PageId, rq.TextContent)
	if err != nil {
		return parseDbErr(err)
	}
	m := messages.Empty{}
	return encodeJSON(w, r, 200, m)
}

func (s *APIServer) getPages(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	pages, err := s.repo.GetPages(ctx)
	if err != nil {
		return parseDbErr(err)
	}
	return encodeJSON(w, r, 200, pages)
}

func (s *APIServer) getBundledPageWithRev(w http.ResponseWriter, r *http.Request) error {
	title := r.PathValue("page_title")
	revId, err := parseUintParam(r, "rev_id")
	if err != nil {
		return BadRequestErr(err)
	}
	ctx := r.Context()
	bundle, err := s.repo.GetBundledPageWithRev(ctx, title, revId)
	if err != nil {
		return parseDbErr(err)
	}
	return encodeJSON(w, r, 200, bundle)
}

func (s *APIServer) getRawTextForPageWithRev(w http.ResponseWriter, r *http.Request) error {
	revId, err := parseUintParam(r, "rev_id")
	if err != nil {
		return BadRequestErr(err)
	}
	ctx := r.Context()
	text, err := s.repo.GetTextByRevID(ctx, revId)
	if err != nil {
		return parseDbErr(err)
	}
	return encodeJSON(w, r, 200, text)
}

func (s *APIServer) getPageRevs(w http.ResponseWriter, r *http.Request) error {
	title := r.PathValue("page_title")
	ctx := r.Context()
	revs, err := s.repo.GetPageRevs(ctx, title)
	if err != nil {
		return parseDbErr(err)
	}
	return encodeJSON(w, r, 200, revs)
}

func (s *APIServer) testHandler(w http.ResponseWriter, r *http.Request) error {
	return encodeJSON(w, r, 200, "HELLO")
}
