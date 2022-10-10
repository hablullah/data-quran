package dl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Request struct {
	URL      string
	FileName string
}

type BatchOption struct {
	NWorker int
	Delay   time.Duration
}

func Download(ctx context.Context, client *http.Client, dstPath string, req Request) error {
	// Download the url
	logrus.Printf("downloading %s", req.URL)
	resp, err := client.Get(req.URL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	// Save to destination
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
}

func BatchDownload(ctx context.Context, dstDir string, requests []Request, opts *BatchOption) error {
	// Prepare opts
	var nWorker int
	var delay time.Duration
	if opts != nil {
		nWorker = opts.NWorker
		delay = opts.Delay
	}

	if nWorker <= 0 {
		nWorker = runtime.GOMAXPROCS(0)
	}

	// Prepare semaphore and error group
	sem := semaphore.NewWeighted(int64(nWorker))
	g, ctx := errgroup.WithContext(ctx)

	// Prepare http client
	client := &http.Client{}

	// Download each request
	for _, req := range requests {
		req := req

		// Acquire semaphore
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}

		g.Go(func() error {
			defer sem.Release(1)
			dstPath := filepath.Join(dstDir, req.FileName)
			err := Download(ctx, client, dstPath, req)
			if err != nil {
				return err
			}

			time.Sleep(delay)
			return nil
		})
	}

	// Wait until all process finished
	return g.Wait()
}
