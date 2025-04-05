package webserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dev-mackan/gowiki/internal/messages"
	"github.com/dev-mackan/gowiki/internal/middleware"
	"github.com/dev-mackan/gowiki/pkg/models"
	"github.com/dev-mackan/gowiki/pkg/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type WebServer struct {
	listenAddr string
	apiAddr    string
	html       *Templates
}

func NewWebServer(config *WebServerConfig) *WebServer {
	return &WebServer{
		config.listenAddr,
		config.apiAddr,
		newTemplate(config.templatePaths),
	}
}

func (s *WebServer) mount() *http.ServeMux {
	fs := http.FileServer(http.Dir("web/templates/css"))
	logger := middleware.NewLoggerMiddleware("WEB-SERVER")
	router := http.NewServeMux()
	router.Handle("GET /static/css/", http.StripPrefix("/static/css/", fs))
	//router.Handle("GET /", logger(s.makeApiHandlerFunc(s.testHandler)))
	router.Handle("GET /", logger(s.makeApiHandlerFunc(s.indexHandler)))
	router.Handle("GET /pages", logger(s.makeApiHandlerFunc(s.indexHandler)))
	router.Handle("GET /pages/new", logger(s.makeApiHandlerFunc(s.newPageGETHandler)))
	router.Handle("POST /pages/new", logger(s.makeApiHandlerFunc(s.newPagePOSTHandler)))
	router.Handle("GET /pages/{page_title}", logger(s.makeApiHandlerFunc(s.pageHandler)))
	router.Handle("GET /pages/{page_title}/edit", logger(s.makeApiHandlerFunc(s.editPageGETHandler)))
	router.Handle("POST /pages/{page_title}/edit", logger(s.makeApiHandlerFunc(s.editPagePOSTHandler)))
	router.Handle("GET /pages/{page_title}/delete", logger(s.makeApiHandlerFunc(s.deletePageGETHandler)))
	router.Handle("POST /pages/{page_title}/delete", logger(s.makeApiHandlerFunc(s.deletePagePOSTHandler)))
	router.Handle("DELETE /pages/{page_title}/delete", logger(s.makeApiHandlerFunc(s.testHandler)))
	router.Handle("GET /pages/{page_title}/revisions", logger(s.makeApiHandlerFunc(s.revisionsHandler)))
	router.Handle("GET /pages/{page_title}/revisions/{rev_id}", logger(s.makeApiHandlerFunc(s.pageWithRevHandler)))
	router.Handle("GET /pages/{page_title}/revisions/{rev_id}/raw.md", logger(s.makeApiHandlerFunc(s.rawTextHandler)))
	return router
}

func (s *WebServer) Run() error {
	mux := s.mount()
	log.Println("GOWIKI-WEB listening on: ", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, mux)
}

func (s *WebServer) testHandler(w http.ResponseWriter, r *http.Request) error {
	return s.html.Render(w, "index", 200, nil)
}

func (s *WebServer) indexHandler(w http.ResponseWriter, r *http.Request) error {
	url := fmt.Sprintf("%s/pages", s.apiAddr)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	pagesBytes, err := readRespBytes(resp)
	if err != nil {
		return err
	}
	var pages []models.Page
	err = json.Unmarshal(pagesBytes, &pages)
	if err != nil {
		return err
	}
	return s.html.Render(w, "index", 200, &pages)
}

func (s *WebServer) pageHandler(w http.ResponseWriter, r *http.Request) error {
	pageTitle := r.PathValue("page_title")
	url := fmt.Sprintf("%s/bundled/%s", s.apiAddr, pageTitle)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	bundleBytes, err := readRespBytes(resp)
	if err != nil {
		log.Println(err)
		return err
	}
	var bundle models.PageBundle
	err = json.Unmarshal(bundleBytes, &bundle)
	if err != nil {
		log.Println(err)
		return err
	}
	mdBuf, err := MarkdownToHTML([]byte(bundle.Text.Content))
	if err != nil {
		return err
	}
	bundle.Text.Content = mdBuf.String()
	return s.html.Render(w, "page", 200, bundle)
}

func (s *WebServer) rawTextHandler(w http.ResponseWriter, r *http.Request) error {
	revId := r.PathValue("rev_id")
	url := fmt.Sprintf("%s/revisions/%s/text/raw", s.apiAddr, revId)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	textBytes, err := readRespBytes(resp)
	if err != nil {
		log.Println("read", err)
		return err
	}
	var text models.Text
	err = json.Unmarshal(textBytes, &text)
	if err != nil {
		log.Println(err)
		return err
	}
	// NOTE: No need to use templates for raw text
	_, err = w.Write([]byte(text.Content))
	return err
}

func (s *WebServer) pageWithRevHandler(w http.ResponseWriter, r *http.Request) error {
	pageTitle := r.PathValue("page_title")
	revId := r.PathValue("rev_id")
	url := fmt.Sprintf("%s/bundled/%s/revisions/%s", s.apiAddr, pageTitle, revId)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	bundleBytes, err := readRespBytes(resp)
	if err != nil {
		log.Println(err)
		return err
	}
	var bundle models.PageBundle
	err = json.Unmarshal(bundleBytes, &bundle)
	if err != nil {
		log.Println(err)
		return err
	}
	mdBuf, err := MarkdownToHTML([]byte(bundle.Text.Content))
	if err != nil {
		return err
	}
	bundle.Text.Content = mdBuf.String()
	return s.html.Render(w, "page", 200, bundle)
}

func (s *WebServer) revisionsHandler(w http.ResponseWriter, r *http.Request) error {
	pageTitle := r.PathValue("page_title")
	pageurl := fmt.Sprintf("%s/pages/%s", s.apiAddr, pageTitle)
	revsurl := fmt.Sprintf("%s/pages/%s/revisions", s.apiAddr, pageTitle)
	resp, err := http.Get(pageurl)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err := readRespBytes(resp)
	if err != nil {
		log.Println(err)
		return err
	}
	var page models.Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return err
	}
	resp, err = http.Get(revsurl)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err = readRespBytes(resp)
	if err != nil {
		log.Println(err)
		return err
	}
	var revs []models.Revision
	err = json.Unmarshal(body, &revs)
	if err != nil {
		log.Println(err)
		return err
	}
	pageRevs := NewPageRevsTmplModel(&page, &revs)
	return s.html.Render(w, "revisions", 200, &pageRevs)
}

func (s *WebServer) editPageGETHandler(w http.ResponseWriter, r *http.Request) error {
	pageTitle := r.PathValue("page_title")
	pageurl := fmt.Sprintf("%s/pages/%s", s.apiAddr, pageTitle)
	//TODO: Should not use a default client here. A client should be part of the server config
	resp, err := http.DefaultClient.Get(pageurl)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err := readRespBytes(resp)
	if err != nil {
		log.Println(err)
		return err
	}
	var page models.Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Println(err)
		return err
	}
	return s.html.Render(w, "edit", 200, &page)
}

func (s *WebServer) editPagePOSTHandler(w http.ResponseWriter, r *http.Request) error {
	log.Println(r.Form)
	pageIdStr := r.FormValue("page_id")
	pageId, err := utils.ParseUintFromStr(pageIdStr)
	if err != nil {
		log.Println(err)
		return err
	}
	formAction := r.FormValue("form_action")
	if formAction == "editContent" {
		err = s.editPageContentHelper(w, r, pageId)
	} else if formAction == "editName" {
		err = s.editPageTitleHelper(w, r, pageId)
	} else {
		err = errors.New(string(http.StatusNotFound))
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *WebServer) editPageContentHelper(w http.ResponseWriter, r *http.Request, pageId uint) error {
	fileBytes, err := parseContentFileFromForm(r)
	if err != nil {
		log.Println(err)
		return err
	}
	reqStruct := messages.UpdatePageContentRequest{
		PageId:      pageId,
		TextContent: string(fileBytes),
	}
	reqBytes, err := json.Marshal(&reqStruct)
	if err != nil {
		log.Println(err)
		return err
	}
	reader := bytes.NewBuffer(reqBytes)
	url := fmt.Sprintf("%s/pages/%d/update/content", s.apiAddr, reqStruct.PageId)
	err = sendReq(http.MethodPut, url, reader)
	if err != nil {
		return err
	}
	pageTitle := r.PathValue("page_title")
	redirectUrl := fmt.Sprintf("/pages/%s", utils.SanitizeTitle(pageTitle))
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	return nil
}

func (s *WebServer) editPageTitleHelper(w http.ResponseWriter, r *http.Request, pageId uint) error {
	newTitle := r.FormValue("page_title")
	reqStruct := messages.UpdatePageTitleRequest{PageId: pageId, PageTitle: newTitle}
	reqBytes, err := json.Marshal(&reqStruct)
	if err != nil {
		log.Println(err)
		return err
	}
	reader := bytes.NewBuffer(reqBytes)
	url := fmt.Sprintf("%s/pages/%d/update/title", s.apiAddr, reqStruct.PageId)
	err = sendReq(http.MethodPut, url, reader)
	if err != nil {
		return err
	}
	redirectUrl := fmt.Sprintf("/pages/%s", utils.SanitizeTitle(newTitle))
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	return nil
}

func (s *WebServer) newPageGETHandler(w http.ResponseWriter, r *http.Request) error {
	return s.html.Render(w, "new", 200, nil)
}

func (s *WebServer) newPagePOSTHandler(w http.ResponseWriter, r *http.Request) error {
	//TODO: Validate form values
	fileBytes, err := parseContentFileFromForm(r)
	if err != nil {
		log.Println(err)
		return err
	}
	br := messages.NewBundleRequest{
		PageTitle:   r.FormValue("page_title"),
		TextContent: string(fileBytes),
	}
	brJson, err := json.Marshal(&br)
	if err != nil {
		log.Println(err)
		return err
	}
	if br.PageTitle == "" {
		return fmt.Errorf("Provide a title")
	}
	reader := bytes.NewBuffer(brJson)
	url := fmt.Sprintf("%s/bundled/new", s.apiAddr)
	//TODO: Should not use a default client here. A client should be part of the server config
	resp, err := http.DefaultClient.Post(url, "application/json", reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		//TODO: Add a unique error here
		log.Println(resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}
	redirectUrl := fmt.Sprintf("/pages/%s", utils.SanitizeTitle(br.PageTitle))
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	return nil
}

func (s *WebServer) deletePageGETHandler(w http.ResponseWriter, r *http.Request) error {
	pageTitle := r.PathValue("page_title")
	pageurl := fmt.Sprintf("%s/pages/%s", s.apiAddr, pageTitle)
	//TODO: Should not use a default client here. A client should be part of the server config
	resp, err := http.DefaultClient.Get(pageurl)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err := readRespBytes(resp)
	if err != nil {
		log.Println(err)
		return err
	}
	var page models.Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Println(err)
		return err
	}
	return s.html.Render(w, "delete", 200, &page)
}

func (s *WebServer) deletePagePOSTHandler(w http.ResponseWriter, r *http.Request) error {
	pageIdStr := r.FormValue("page_id")
	pageId, err := strconv.ParseUint(pageIdStr, 10, 64)
	if err != nil {
		log.Println(err)
		return err
	}
	dr := messages.DeletePageRequest{PageId: uint(pageId)}
	reqBytes, err := json.Marshal(&dr)
	if err != nil {
		log.Println(err)
		return err
	}
	reader := bytes.NewBuffer(reqBytes)
	url := fmt.Sprintf("%s/bundled/delete", s.apiAddr)
	apiReq, err := http.NewRequest(http.MethodDelete, url, reader)
	if err != nil {
		return err
	}
	//TODO: Should not use a default client here. A client should be part of the server config
	resp, err := http.DefaultClient.Do(apiReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//TODO: Add a unique error here
		log.Println(resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}

	redirectUrl := fmt.Sprintf("/pages")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	return nil
}

type WebError struct {
	Error string `json:"error"`
}

type webFunc func(http.ResponseWriter, *http.Request) error

func (s *WebServer) makeApiHandlerFunc(f webFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			s.html.Render(w, "error", 500, err.Error())
		}
	}
}

func sendReq(method string, url string, reader io.Reader) error {
	apiReq, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}
	//TODO: Should not use a default client here. A client should be part of the server config
	resp, err := http.DefaultClient.Do(apiReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//TODO: Add a unique error here
		log.Println(resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}
	return nil
}

func readRespBytes(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		//TODO: implement custom error type for this case
		return nil, errors.New(fmt.Sprintf("%d", resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HasMarkdownSuffix(filename string) bool {
	if !strings.HasSuffix(filename, ".md") {
		return false
	}
	return true
}

func IsMarkdownContent(contentType string) bool {
	if contentType != "text/markdown" && contentType != "text/plain" {
		return false
	}
	return true
}

func parseContentFileFromForm(r *http.Request) ([]byte, error) {
	var fileBytes []byte
	if file, handler, err := r.FormFile("content"); err == nil {
		defer file.Close()
		if !HasMarkdownSuffix(handler.Filename) {
			return nil, errors.New("Invalid file type. Only .md files are allowed")
		}
		if !IsMarkdownContent(handler.Header.Get("Content-Type")) {
			return nil, errors.New("Invalid file type. Only Markdown files are allowed")
		}
		fileBytes, err = io.ReadAll(file)
		if err != nil {
			return nil, errors.New("Error reading the file")
		}
	} else if err != http.ErrMissingFile {
		//NOTE: Leave the TextContent field empty
	} else {
		return nil, err
	}
	return fileBytes, nil
}

func parseFormToBundleRequest(r *http.Request) (*messages.BundleRequest, error) {
	var req messages.BundleRequest
	r.ParseMultipartForm(10 << 20) // 10MB
	if file, handler, err := r.FormFile("content"); err == nil {
		defer file.Close()
		if !HasMarkdownSuffix(handler.Filename) {
			return nil, errors.New("Invalid file type. Only .md files are allowed")
		}
		if !IsMarkdownContent(handler.Header.Get("Content-Type")) {
			return nil, errors.New("Invalid file type. Only Markdown files are allowed")
		}
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, errors.New("Error reading the file")
		}
		req.TextContent = string(fileBytes)
	} else if err != http.ErrMissingFile {
		//NOTE: Leave the TextContent field empty
	}
	title := r.FormValue("page_title")
	if title == "" {
		return nil, fmt.Errorf("No parameters provided")
	}
	req.PageTitle = title
	return &req, nil
}
