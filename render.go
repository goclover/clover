package clover

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// HeaderContentTyp  HTTP Header 中 Content-Type 的 Key
const HeaderContentTyp = "Content-Type"

// HeaderContentLen HTTP Header 中 Content-Length 的 Key
const HeaderContentLen = "Content-Length"

var JSON = func(data interface{}) *JSONRender {
	bf, _ := json.Marshal(data)
	return &JSONRender{
		NopRender: NopRender{
			Status: 0,
			Headers: http.Header{
				HeaderContentTyp: []string{"application/json; charset=utf-8"},
				HeaderContentLen: []string{strconv.Itoa(len(bf))},
			},
		},
		Data: bf,
	}
}

var Text = func(text string) *TextRender {
	bf := []byte(text)
	return &TextRender{
		NopRender: NopRender{
			Status: 0,
			Headers: http.Header{
				HeaderContentTyp: []string{"text/plain; charset=utf-8"},
				HeaderContentLen: []string{strconv.Itoa(len(bf))},
			},
		},
		Text: bf,
	}
}

type Render interface {
	WriteTo(w http.ResponseWriter) error
}

type NopRender struct {
	Status  int
	Headers http.Header
}

func (n *NopRender) WriteTo(w http.ResponseWriter) error {
	copyHeaders(w.Header(), n.Headers)

	if n.Status > 0 {
		w.WriteHeader(n.Status)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	return nil
}

type JSONRender struct {
	NopRender
	Data []byte
}

func (j JSONRender) WriteTo(w http.ResponseWriter) error {
	_ = j.NopRender.WriteTo(w)
	_, errW := w.Write(j.Data)
	return errW
}

type TextRender struct {
	NopRender
	Text []byte
}

func (t TextRender) WriteTo(w http.ResponseWriter) error {
	_ = t.NopRender.WriteTo(w)
	_, errW := w.Write(t.Text)
	return errW
}

func copyHeaders(dst http.Header, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}
