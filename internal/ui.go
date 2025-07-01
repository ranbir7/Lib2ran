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
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
)

func ShowWelcome() {
	mag := color.New(color.FgHiMagenta, color.Bold)
	cyan := color.New(color.FgHiCyan, color.Bold)
	mag.Println("╔════════════════════════════════════════════════════════╗")
	mag.Print("║")
	fmt.Print(" ")
	cyan.Print("Lib2ran – The Ultimate LibGen CLI")
	mag.Println("                         		║")
	mag.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func ShowInfo(msg string) {
	color.New(color.FgCyan, color.Bold).Println(msg)
}

func ShowError(msg string) {
	red := color.New(color.FgHiRed, color.Bold)
	red.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	red.Printf("✖ %s\n", msg)
	red.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}

func ShowSuccess(msg string) {
	green := color.New(color.FgHiGreen, color.Bold)
	green.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	green.Printf("✔ %s\n", msg)
	green.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
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
	table.SetRowSeparator("-")
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
	prompt := promptui.Prompt{
		Label: "Enter book title to search",
	}
	result, _ := prompt.Run()
	return result
}

func PromptSelectResult(books []Book) *Book {
	cyan := color.New(color.FgHiCyan, color.Bold)
	fmt.Println()
	cyan.Println(strings.Repeat("━", 60))
	cyan.Printf("%-60s\n", "[ENTER THE # NUMBER OF THE BOOK TO DOWNLOAD]")
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
