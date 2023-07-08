package kemenag

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func downloadListSurah(cacheDir string) error {
	ctx := context.Background()
	dstPath := filepath.Join(cacheDir, "list-surah.json")

	if !util.FileExist(dstPath) {
		req := dl.Request{URL: "https://web-api.qurankemenag.net/quran-surah"}
		err := dl.Download(ctx, http.DefaultClient, dstPath, req)
		if err != nil {
			return fmt.Errorf("failed to download list surah: %w", err)
		}
	}

	return nil
}

func downladAllTafsir(cacheDir string) error {
	for surah := 1; surah <= 114; surah++ {
		err := downloadTafsir(cacheDir, surah)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadTafsir(cacheDir string, surah int) error {
	// Prepare destination path
	dstName := fmt.Sprintf("surah-%03d.json", surah)
	dstPath := filepath.Join(cacheDir, dstName)
	if util.FileExist(dstPath) {
		return nil
	}

	// Prepare http client
	client := &http.Client{}

	// Download each tafsir
	surahData := util.ListSurah[surah]
	tafsirs := make([]Ayah, surahData.NAyah)

	for idx := 1; idx <= surahData.NAyah; idx++ {
		ayah := surahData.Start + idx - 1
		err := func() error {
			logrus.Printf("downloading tafsir for %d:%d", surah, idx)

			// Download page
			url := fmt.Sprintf("https://web-api.qurankemenag.net/quran-tafsir/%d", ayah)
			resp, err := client.Get(url)
			if err != nil {
				return fmt.Errorf("failed to download tafsir for %d:%d, %w", surah, idx, err)
			}
			defer resp.Body.Close()

			// Decode data
			var respData RespDownloadTafsir
			err = json.NewDecoder(resp.Body).Decode(&respData)
			if err != nil {
				return fmt.Errorf("failed to decode tafsir for %d:%d %w", surah, idx, err)
			}

			// Save to slice
			tafsirs[idx-1] = respData.Data
			return nil
		}()
		if err != nil {
			return err
		}

		// // Delay
		// time.Sleep(500 * time.Millisecond)
	}

	// Write tafsirs to file
	bt, err := json.MarshalIndent(tafsirs, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to write tafsir for surah %d: %w", surah, err)
	}

	return os.WriteFile(dstPath, bt, os.ModePerm)
}
