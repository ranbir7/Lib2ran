import os
import requests
from libgen_api import LibgenSearch
import inquirer
import argparse
from rich.console import Console
from rich.panel import Panel
from rich.prompt import Prompt, Confirm
from rich.table import Table
from rich.progress import Progress, SpinnerColumn, BarColumn, TextColumn, TimeElapsedColumn
from rich.theme import Theme
from rich import box
import sys
import signal
from rich.rule import Rule

# Theme for easy color changes
custom_theme = Theme({
    "banner": "bold magenta",
    "info": "bold cyan",
    "success": "bold green",
    "warning": "bold yellow",
    "error": "bold red",
    "prompt": "bold blue",
    "table.header": "bold magenta",
    "table.row": "white",
})
console = Console(theme=custom_theme)

# Add a theme variable for easy tweaks
ACCENT_COLOR = "bright_magenta"

ABOUT = """
[bold cyan]lib2ran[/bold cyan] is a beautiful, modern CLI tool to search and download books from Library Genesis (libgen.rs).
Created by [bold yellow]ChatGPT[/bold yellow] with ❤️ 
"""

VERSION = "1.0.0"

def graceful_exit(signum=None, frame=None):
    console.print(Panel("[bold magenta]Exiting program. Goodbye! 👋[/bold magenta]", expand=False))
    sys.exit(0)
signal.signal(signal.SIGINT, graceful_exit)

def download_book(book, searcher):
    """
    Download a book using the resolved download link, with a progress bar.
    """
    try:
        download_links = searcher.resolve_download_links(book)
        if not download_links:
            console.print(Panel(":x: [error]No download links found for this book.[/error]", style="error", expand=False))
            return False
        for source_name, url in download_links.items():
            try:
                response = requests.get(url, stream=True, timeout=10)
                response.raise_for_status()
                filename = f"{book['Title']} - {book['Author']}.{book['Extension']}"
                filename = "".join(c for c in filename if c.isalnum() or c in (".", "_", "-")).rstrip()
                filepath = os.path.join(os.getcwd(), filename)
                total = int(response.headers.get('content-length', 0))
                with open(filepath, "wb") as f, Progress(
                    SpinnerColumn(),
                    BarColumn(),
                    TextColumn("[progress.percentage]{task.percentage:>3.0f}%"),
                    TimeElapsedColumn(),
                    console=console,
                    transient=True
                ) as progress:
                    task = progress.add_task(f"[cyan]Downloading[/cyan] [bold]{filename}[/bold]", total=total)
                    for chunk in response.iter_content(chunk_size=1024):
                        if chunk:
                            f.write(chunk)
                            progress.update(task, advance=len(chunk))
                console.print(Panel(f":white_check_mark: [success]Downloaded:[/success] {filepath}", style="success", expand=False))
                return True
            except requests.RequestException:
                console.print(Panel(f":warning: [warning]Failed to download from {source_name}. Trying next link...[/warning]", style="warning", expand=False))
                continue
        console.print(Panel(":x: [error]All download attempts failed.[/error]", style="error", expand=False))
        return False
    except Exception as e:
        console.print(Panel(f":x: [error]Error downloading file: {e}[/error]", style="error", expand=False))
        return False

def display_results(results):
    """
    Display search results and prompt user to select a book.
    """
    if not results:
        console.print(Panel(":x: [error]No results found.[/error]", style="error", expand=False))
        return None
    limited_results = results[:min(len(results), 100000)]
    table = Table(title="Search Results", show_lines=True, header_style="table.header", box=box.ROUNDED)
    table.add_column("#", style="cyan", width=4)
    table.add_column("Title", style="bold")
    table.add_column("Author", style="green")
    table.add_column("Year", style="yellow")
    table.add_column("Size", style="blue")
    table.add_column("Ext", style="magenta")
    for i, book in enumerate(limited_results):
        table.add_row(str(i+1), book['Title'], book['Author'], book['Year'], book['Size'], book['Extension'])
    console.print(table)
    while True:
        try:
            choice = Prompt.ask(":books: [prompt]Enter the number of the book to download[/prompt]", default="1")
            selected_index = int(choice) - 1
            if 0 <= selected_index < len(limited_results):
                return limited_results[selected_index]
            else:
                console.print(Panel(":warning: [warning]Invalid selection. Please try again.[/warning]", style="warning", expand=False))
        except (ValueError, KeyboardInterrupt):
            console.print(Panel(":x: [error]Invalid input or cancelled. Please try again.[/error]", style="error", expand=False))

def main():
    """
    Main function to run the LibGen book downloader.
    Supports both interactive and command-line modes.
    """
    parser = argparse.ArgumentParser(description="LibGen Book Downloader by Ran7")
    parser.add_argument("-b", "--book", metavar="BOOK_TITLE", type=str, nargs='?', help="Download a book by title directly")
    parser.add_argument("--about", action="store_true", help="Show about information and exit")
    parser.add_argument("--version", action="store_true", help="Show version and exit")
    parser.add_argument("--auto", action="store_true", help="Automatically download the top search result")
    args = parser.parse_args()
    # Slim, elegant header with accent line
    console.print("\n")
    console.print(Rule("lib2ran – Library Genesis Book Downloader", style=ACCENT_COLOR))
    console.print("\n")
    if args.about:
        console.print(Panel(ABOUT, title="About lib2ran", style="info", expand=False))
        sys.exit(0)
    if args.version:
        console.print(Panel(f"lib2ran version {VERSION}", style="info", expand=False))
        sys.exit(0)
    s = LibgenSearch()
    if args.book:
        query = args.book.strip()
        with Progress(SpinnerColumn(), TextColumn("[cyan]Searching...[/cyan]"), console=console, transient=True) as progress:
            progress.add_task("search", total=None)
            try:
                results = s.search_title(query)
            except Exception as e:
                console.print(Panel(f":x: [error]Error: {e}[/error]", style="error", expand=False))
                return
        if not results:
            console.print(Panel(":x: [error]No results found.[/error]", style="error", expand=False))
            return
        if args.auto:
            selected_book = results[0]
            console.print(Panel(f":rocket: [info]Auto-downloading top result: [bold]{selected_book['Title']}[/bold][/info]", style="info", expand=False))
        else:
            selected_book = display_results(results)
            if not selected_book:
                console.print(Panel(":warning: [warning]No book selected or no results found.[/warning]", style="warning", expand=False))
                return
        console.print(Panel(f":open_book: [success]Downloading [cyan]{selected_book['Title']}[/cyan]...[/success]", style="success", expand=False))
        download_book(selected_book, s)
        return
    # Fallback to interactive mode
    while True:
        try:
            query = Prompt.ask(":mag: [prompt]\nEnter the book name[/prompt] (or type 'quit' to exit)").strip()
        except (KeyboardInterrupt, EOFError):
            graceful_exit()
        if query.lower() == 'quit':
            graceful_exit()
        if not query:
            console.print(Panel(":warning: [warning]Book name cannot be empty.[/warning]", style="warning", expand=False))
            continue
        with Progress(SpinnerColumn(), TextColumn("[cyan]Searching...[/cyan]"), console=console, transient=True) as progress:
            progress.add_task("search", total=None)
            try:
                results = s.search_title(query)
            except Exception as e:
                console.print(Panel(f":x: [error]Error searching LibGen: {e}[/error]", style="error", expand=False))
                continue
        selected_book = display_results(results)
        if not selected_book:
            console.print(Panel(":warning: [warning]No book selected.[/warning]", style="warning", expand=False))
            continue
        console.print(Panel(f":open_book: [success]Downloading [cyan]{selected_book['Title']}[/cyan]...[/success]", style="success", expand=False))
        download_book(selected_book, s)
        if not Confirm.ask(":repeat: [prompt]Do you want to download another book?[/prompt]", default=True):
            graceful_exit()

if __name__ == "__main__":
    main()