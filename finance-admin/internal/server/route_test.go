package server

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"io"
	"net/http"
	"testing"
)

func TestRoute_fullPage(t *testing.T) {
	sut := route{client: mockApiClient{}}
	req, _ := http.NewRequest(http.MethodGet, "", nil)

	r, w := io.Pipe()
	go func() {
		_ = sut.execute(w, req, components.Error(components.ErrorVars{}))
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}
	// Expect the component to be present.
	if doc.Find(`[data-testid="error-template"]`).Length() == 0 {
		t.Error("expected data-testid attribute to be rendered, but it wasn't")
	}
	// Expect the parent to also be present
	if doc.Find(`[data-testid="content-header-template"]`).Length() == 0 {
		t.Error("expected data-testid attribute to not be rendered, but it was")
	}
}

func TestRoute_htmxRequest(t *testing.T) {
	sut := route{client: mockApiClient{}}
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	req.Header.Add("HX-Request", "true")

	r, w := io.Pipe()
	go func() {
		_ = sut.execute(w, req, components.Error(components.ErrorVars{}))
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}
	// Expect the component to be present.
	if doc.Find(`[data-testid="error-template"]`).Length() == 0 {
		t.Error("expected data-testid attribute to be rendered, but it wasn't")
	}
	// Expect the parent to NOT be present
	if doc.Find(`[data-testid="content-header-template"]`).Length() > 0 {
		t.Error("expected data-testid attribute to not be rendered, but it was")
	}
}
