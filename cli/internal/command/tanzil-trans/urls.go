package tanzilTrans

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/zyedidia/generic/mapset"
)

func downloadTranslationPage(cacheDir string) error {
	// Check if Tanzil index file exist
	indexPath := filepath.Join(cacheDir, "index.html")
	if util.FileExist(indexPath) {
		return nil
	}

	// Download file
	ctx := context.Background()
	client := &http.Client{}
	err := dl.Download(ctx, client, indexPath, dl.Request{
		URL:      "https://tanzil.net/trans/",
		FileName: "index.html"})
	if err != nil {
		return fmt.Errorf("failed to download index: %w", err)
	}

	return nil
}

func parseTranslationPage(cacheDir string) ([]dl.Request, error) {
	// Open index file
	indexPath := filepath.Join(cacheDir, "index.html")
	f, err := os.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open index: %w", err)
	}
	defer f.Close()

	// Parse HTML document
	doc, err := dom.FastParse(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}

	// Extract URLs
	var dlRequests []dl.Request
	urlSet := mapset.New[string]()
	baseURL, _ := url.ParseRequestURI("https://tanzil.net/")

	selector := "table.transList a.download"
	for _, link := range dom.QuerySelectorAll(doc, selector) {
		// Fetch the href, and make sure it hasn't been captured
		href := dom.GetAttribute(link, "href")
		href = strings.TrimSpace(href)
		if urlSet.Has(href) {
			continue
		}

		// Normalize href
		urlSet.Put(href)
		baseURL.Path = href
		href = baseURL.String()

		// Create the file name
		fName := filepath.Base(href)
		fName = strings.ReplaceAll(fName, ".", "-")
		fName += ".txt"

		// Save the URL
		dlRequests = append(dlRequests, dl.Request{
			URL:      href + "?type=txt",
			FileName: fName,
		})
	}

	return dlRequests, nil
}
