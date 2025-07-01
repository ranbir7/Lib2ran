package main

import (
	"bufio"
	"fmt"
	"lib2ran/internal"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	internal.ShowWelcome()
	reader := bufio.NewReader(os.Stdin)
	for {
		query := internal.GetUserQuery()
		results := internal.SearchLibgen(query)
		if len(results) == 0 {
			internal.ShowError("No results found.")
			continue
		}
		internal.ShowResultsTable(results)
		for {
			selected := internal.PromptSelectResult(results)
			if selected == nil {
				internal.ShowError("No book selected.")
				break
			}
			internal.ShowInfo(fmt.Sprintf("Downloading: %s", selected.Title))
			downloadsDir, _ := os.UserHomeDir()
			downloadsDir = filepath.Join(downloadsDir, "Downloads")
			err := internal.DownloadBook(selected, downloadsDir)
			if err != nil {
				internal.ShowError(fmt.Sprintf("Download failed: %v", err))
			} else {
				internal.ShowSuccess(fmt.Sprintf("Downloaded to %s", downloadsDir))
			}
			fmt.Println()
			fmt.Print("Do you want to download another file from the same results? (y/n): ")
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer == "y" || answer == "yes" {
				internal.ShowResultsTable(results)
				continue
			}
			fmt.Print("Do you want to search for new books? (y/n): ")
			newSearch, _ := reader.ReadString('\n')
			newSearch = strings.TrimSpace(strings.ToLower(newSearch))
			if newSearch == "y" || newSearch == "yes" {
				break // break inner loop, start new search
			}
			internal.ShowInfo("Thank you for using Lib2ran! Have a premium day!")
			return
		}
	}
}
