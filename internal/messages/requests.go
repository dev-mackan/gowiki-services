package messages

type BundleRequest struct {
	PageId      uint   `json:"page_id,omitempty"`
	PageTitle   string `json:"page_title"`
	TextContent string `json:"text_content,omitempty"`
}

type UpdateBundleRequest struct {
	PageId      uint   `json:"page_id,omitempty"`
	PageTitle   string `json:"page_title"`
	TextContent string `json:"text_content,omitempty"`
}

type UpdatePageContentRequest struct {
	PageId      uint   `json:"page_id,omitempty"`
	TextContent string `json:"text_content,omitempty"`
}

type UpdatePageTitleRequest struct {
	PageId    uint   `json:"page_id,omitempty"`
	PageTitle string `json:"page_title"`
}

type NewBundleRequest struct {
	PageTitle   string `json:"page_title"`
	TextContent string `json:"text_content,omitempty"`
}

type DeletePageRequest struct {
	PageId uint `json:"page_id"`
}
