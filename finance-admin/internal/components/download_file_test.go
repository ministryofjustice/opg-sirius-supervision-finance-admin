package components

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	tests := []struct {
		name   string
		vars   DownloadFileVars
		assert func(doc *goquery.Document)
	}{
		{
			name: "button rendered correctly",
			vars: DownloadFileVars{
				ErrorMessage: "",
				Filename:     "test.csv",
				Uid:          "abc123",
				AppVars: AppVars{
					EnvironmentVars: EnvironmentVars{
						Prefix: "/finance-admin",
					},
				},
			},
			assert: func(doc *goquery.Document) {
				if btn := doc.Find("#download-button"); btn.Length() == 0 {
					t.Error("expected download button to be rendered, but it wasn't")
				} else {
					href, _ := btn.Attr("href")
					assert.Equal(t, "/finance-admin/download/callback?uid=abc123", href)
				}
			},
		},
		{
			name: "error message",
			vars: DownloadFileVars{
				ErrorMessage: "The file cannot be downloaded",
				Filename:     "test.csv",
				Uid:          "abc123",
				AppVars: AppVars{
					EnvironmentVars: EnvironmentVars{
						Prefix: "/finance-admin",
					},
				},
			},
			assert: func(doc *goquery.Document) {
				if btn := doc.Find("#download-button"); btn.Length() > 0 {
					t.Error("expected download button not to be rendered, but it was!")
				}
				if msg := doc.Find(`[data-testid="error-message"]`).Text(); msg != "The file cannot be downloaded" {
					t.Errorf("expected error message 'The file cannot be downloaded', got %q", msg)
				}
			},
		},
	}
	for _, test := range tests {
		r, w := io.Pipe()
		go func() {
			_ = DownloadFile(test.vars).Render(context.Background(), w)
			_ = w.Close()
		}()
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			t.Fatalf("failed to read template: %v", err)
		}
		test.assert(doc)
	}
}
