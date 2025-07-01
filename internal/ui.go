package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func ShowWelcome() {
	mag := color.New(color.FgHiMagenta, color.Bold)
	cyan := color.New(color.FgHiCyan, color.Bold)
	blue := color.New(color.FgHiBlue, color.Bold)
	fmt.Println()
	mag.Println("╔" + strings.Repeat("═", 58) + "╗")
	mag.Print("║")
	fmt.Print(strings.Repeat(" ", 6))
	cyan.Print("🌟 ")
	blue.Print("Lib2ran")
	cyan.Print(" – The Ultimate LibGen CLI ")
	mag.Print("🌟")
	fmt.Print(strings.Repeat(" ", 6))
	mag.Println("	   ║")
	mag.Println("╚" + strings.Repeat("═", 58) + "╝")
	fmt.Println()
	// Animated spinner for a second
	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
	s.Suffix = "  Initializing premium experience..."
	s.Color("magenta")
	s.Start()
	time.Sleep(1 * time.Second)
	s.Stop()
	fmt.Println()
}

func ShowInfo(msg string) {
	blue := color.New(color.FgHiBlue, color.Bold)
	fmt.Println()
	blue.Printf("┃ ℹ %s\n", msg)
	fmt.Println()
}

func ShowError(msg string) {
	red := color.New(color.FgHiRed, color.Bold)
	fmt.Println()
	red.Println("╭" + strings.Repeat("━", 58) + "╮")
	red.Printf("│ ✖ %s%s│\n", msg, strings.Repeat(" ", 56-len(msg)))
	red.Println("╰" + strings.Repeat("━", 58) + "╯")
	fmt.Println()
}

func ShowSuccess(msg string) {
	green := color.New(color.FgHiGreen, color.Bold)
	fmt.Println()
	green.Println("╭" + strings.Repeat("━", 58) + "╮")
	green.Printf("│ ✔ %s%s│\n", msg, strings.Repeat(" ", 56-len(msg)))
	green.Println("╰" + strings.Repeat("━", 58) + "╯")
	fmt.Println()
}

func ShowResultsTable(books []Book) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Title", "Author", "Year", "Size", "Ext"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiMagentaColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
	)
	table.SetRowLine(true)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowSeparator("━")
	for i, b := range books {
		rowColor := tablewriter.Colors{tablewriter.FgWhiteColor}
		if i%2 == 1 {
			rowColor = tablewriter.Colors{tablewriter.FgHiBlackColor}
		}
		table.Rich([]string{
			fmt.Sprintf("%d", i+1), b.Title, b.Author, b.Year, b.Size, b.Extension,
		}, []tablewriter.Colors{rowColor, rowColor, rowColor, rowColor, rowColor, rowColor})
	}
	table.Render()
}

func ShowSpinner(msg string, fn func()) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + msg
	s.Color("magenta")
	s.Start()
	fn()
	s.Stop()
}

func GetUserQuery() string {
	cyan := color.New(color.FgHiCyan, color.Bold)
	promptLine := "[🔎 ENTER BOOK TITLE TO SEARCH]"
	padding := (60 - len(promptLine)) / 2
	fmt.Println()
	cyan.Println(strings.Repeat("━", 60))
	fmt.Printf("%s%s%s\n", strings.Repeat(" ", padding), promptLine, strings.Repeat(" ", padding))
	cyan.Println(strings.Repeat("━", 60))
	fmt.Println()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Book Title: ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func PromptSelectResult(books []Book) *Book {
	cyan := color.New(color.FgHiCyan, color.Bold)
	promptLine := "[⬇️ ENTER THE # NUMBER OF THE BOOK TO DOWNLOAD]"
	padding := (60 - len(promptLine)) / 2
	fmt.Println()
	cyan.Println(strings.Repeat("━", 60))
	fmt.Printf("%s%s%s\n", strings.Repeat(" ", padding), promptLine, strings.Repeat(" ", padding))
	cyan.Println(strings.Repeat("━", 60))
	fmt.Println()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("# Number: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		idx, err := strconv.Atoi(input)
		if err == nil && idx > 0 && idx <= len(books) {
			return &books[idx-1]
		}
		ShowError("Invalid selection. Please enter a valid # number from the table above.")
	}
}

func ShowGoodbye() {
	mag := color.New(color.FgHiMagenta, color.Bold)
	msg := "✨ Thank you for using Lib2ran! Have a premium day! ✨"
	width := 50
	fmt.Println()
	mag.Println("╔" + strings.Repeat("═", width) + "╗")
	mag.Printf("║%s║\n", centerText(msg, width))
	mag.Println("╚" + strings.Repeat("═", width) + "╝")
	fmt.Println()
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text // No padding if text is too long
	}
	pad := (width - len(text)) / 2
	return strings.Repeat(" ", pad) + text + strings.Repeat(" ", width-len(text)-pad)
}
