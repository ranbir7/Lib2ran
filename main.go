package main

import (
	"fmt"
	"lib2ran/internal"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	internal.ShowWelcome()
	for {
		query := internal.GetUserQuery()
		if query == "" {
			internal.ShowError("Search query cannot be empty.")
			continue
		}
		var books []internal.Book
		internal.ShowSpinner("Searching for books...", func() {
			books = internal.SearchLibgen(query)
		})
		if len(books) == 0 {
			internal.ShowError("No results found. Try another search.")
			continue
		}
		internal.ShowResultsTable(books)
		for {
			book := internal.PromptSelectResult(books)
			if book == nil {
				internal.ShowError("No book selected.")
				continue
			}
			internal.ShowSpinner("Downloading...", func() {
				downloadsDir, _ := os.UserHomeDir()
				downloadsDir = filepath.Join(downloadsDir, "Downloads")
				err := internal.DownloadBook(book, downloadsDir)
				if err != nil {
					internal.ShowError("Download failed: " + err.Error())
				} else {
					internal.ShowSuccess("Download complete! Saved to your Downloads folder.")
				}
			})
			internal.ShowInfo("Would you like to download another book from these results? (y/n)")
			var again string
			fmt.Print("Your choice: ")
			fmt.Scanln(&again)
			if strings.ToLower(strings.TrimSpace(again)) != "y" {
				break
			}
		}
		internal.ShowInfo("Would you like to start a new search? (y/n)")
		var newsearch string
		fmt.Print("Your choice: ")
		fmt.Scanln(&newsearch)
		if strings.ToLower(strings.TrimSpace(newsearch)) != "y" {
			break
		}
	}
	internal.ShowGoodbye()
}
