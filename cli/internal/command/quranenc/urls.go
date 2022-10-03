package quranenc

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/zyedidia/generic/mapset"
)

func downloadIndexPage(cacheDir string) error {
	// Check if quranec index file exist
	indexPath := filepath.Join(cacheDir, "index.html")
	if util.FileExist(indexPath) {
		return nil
	}

	// Download file
	ctx := context.Background()
	client := &http.Client{}
	err := dl.Download(ctx, client, indexPath, dl.Request{
		URL:      "https://quranenc.com/en/home",
		FileName: "index.html"})
	if err != nil {
		return fmt.Errorf("failed to download index: %w", err)
	}

	return nil
}

func parseIndexPage(cacheDir string) ([]dl.Request, error) {
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

	selector := ".modal-body a[href*='/download/xml']"
	for _, link := range dom.QuerySelectorAll(doc, selector) {
		// Fetch the href, and make sure it hasn't been captured
		href := dom.GetAttribute(link, "href")
		href = strings.TrimSpace(href)
		if urlSet.Has(href) {
			continue
		}

		// Create the file name
		fName := filepath.Base(href)
		fName = strings.TrimSuffix(fName, filepath.Ext(fName))
		fName += ".xml"

		// Save the URL
		urlSet.Put(href)
		dlRequests = append(dlRequests, dl.Request{
			URL:      href,
			FileName: fName,
		})
	}

	return dlRequests, nil
}
