package webserver

import (
	"bytes"
	"github.com/yuin/goldmark"
)

func MarkdownToHTML(md []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := goldmark.Convert([]byte(md), &buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
