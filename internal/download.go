package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

func DownloadBook(book *Book, downloadsDir string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	var lastErr error

	for mirrorName, downloadPage := range book.Mirrors {
		color.New(color.FgCyan, color.Bold).Printf("Trying %s...\n", mirrorName)
		realURL, err := getRealFileURL(downloadPage)
		if err != nil || realURL == "" {
			lastErr = err
			continue
		}
		for attempt := 1; attempt <= 3; attempt++ {
			color.New(color.FgHiMagenta).Printf("Download attempt %d...\n", attempt)
			resp, err := client.Get(realURL)
			if err != nil {
				lastErr = err
				time.Sleep(2 * time.Second)
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				lastErr = errors.New("bad status: " + resp.Status)
				continue
			}
			filename := fmt.Sprintf("%s - %s.%s", book.Title, book.Author, book.Extension)
			filename = sanitizeFilename(filename)
			filepath := filepath.Join(downloadsDir, filename)
			out, err := os.Create(filepath)
			if err != nil {
				lastErr = err
				continue
			}
			defer out.Close()
			total := resp.ContentLength
			bar := progressbar.NewOptions64(
				total,
				progressbar.OptionSetDescription("Downloading"),
				progressbar.OptionSetWidth(30),
				progressbar.OptionShowBytes(true),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "#",
					SaucerHead:    ">",
					SaucerPadding: "-",
					BarStart:      "[",
					BarEnd:        "]",
				}),
				progressbar.OptionShowCount(),
				progressbar.OptionShowIts(),
			)
			_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
			if err != nil {
				lastErr = err
				continue
			}
			if bar.State().CurrentBytes < 50*1024 {
				color.New(color.FgYellow, color.Bold).Printf("Warning: Downloaded file is very small and may not be a valid book.\n")
			}
			color.New(color.FgGreen, color.Bold).Printf("Downloaded: %s\n", filepath)
			return nil
		}
	}
	return lastErr
}

// getRealFileURL parses the download page for the real file link
func getRealFileURL(downloadPage string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(downloadPage)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	var fileURL string
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if exists && (strings.HasSuffix(href, ".pdf") || strings.HasSuffix(href, ".epub") ||
			strings.HasSuffix(href, ".mobi") || strings.HasSuffix(href, ".djvu") ||
			strings.HasSuffix(href, ".chm") || strings.HasSuffix(href, ".zip") || strings.HasSuffix(href, ".rar")) {
			if strings.HasPrefix(href, "http") {
				fileURL = href
			} else {
				fileURL = "https://libgen.rs" + href
			}
			return false // break
		}
		return true // continue
	})
	if fileURL == "" {
		return "", errors.New("no real file link found")
	}
	return fileURL, nil
}

func sanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune("\\/:*?\"<>|", r) {
			return -1
		}
		return r
	}, name)
}
