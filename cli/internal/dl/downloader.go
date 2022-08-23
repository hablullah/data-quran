package dl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Request struct {
	URL      string
	FileName string
}

func BatchDownload(ctx context.Context, dstDir string, requests []Request) error {
	// Prepare semaphore and error group
	nWorker := int64(runtime.GOMAXPROCS(0))
	sem := semaphore.NewWeighted(nWorker)
	g, ctx := errgroup.WithContext(ctx)

	// Prepare http client
	client := &http.Client{}

	// Download each request
	for _, req := range requests {
		req := req

		// Acquire semaphore
		if err := sem.Acquire(ctx, 1); err != nil {
			return fmt.Errorf("acquire semaphore failed: %w", err)
		}

		g.Go(func() error {
			defer sem.Release(1)

			// Download the url
			logrus.Printf("downloading %s", req.URL)
			resp, err := client.Get(req.URL)
			if err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
			defer resp.Body.Close()

			// Save to destination
			dstPath := filepath.Join(dstDir, req.FileName)
			dst, err := os.Create(dstPath)
			if err != nil {
				return fmt.Errorf("create dst failed: %w", err)
			}
			defer dst.Close()

			_, err = io.Copy(dst, resp.Body)
			if err != nil {
				return fmt.Errorf("write dst failed: %w", err)
			}

			return nil
		})
	}

	// Wait until all process finished
	return g.Wait()
}
