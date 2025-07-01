# lib2ran – Library Genesis Book Downloader

**lib2ran** is a modern Python CLI tool to search and download books from [Library Genesis (libgen.rs)](https://libgen.rs).

- Clean, professional, and attractive terminal interface
- Uses the [rich](https://github.com/Textualize/rich) library for beautiful output (color, tables, progress bars)
- No big ASCII art or distracting banners—just a slim, elegant accent header
- Both **interactive** and **direct command-line** modes for maximum flexibility

---

## ✨ Features

* Search for books by title on LibGen
* Choose from top results and download the desired book
* Modern, readable tables and prompts
* Download progress bar and search spinner
* Works on Windows, Linux, and macOS
* Supports CLI arguments for scripting/automation
* `--about` and `--version` for info
* `--auto` to instantly download the top result

---

## 🧰 Installation

Clone the repository and install it in **editable mode**:

```bash
git clone https://github.com/ranbir7/Lib2ran.git
cd lib2ran
pip install -e .
```

> This will install the `lib2ran` CLI globally on your system.

---

## 🚀 Usage

### Interactive Mode

Run the tool and follow the prompts:

```bash
lib2ran
```

### Command-Line Mode

Download a book directly using a search term:

```bash
lib2ran -b "The Art of Computer Programming"
```

### Additional Options

- Show about info: `lib2ran --about`
- Show version: `lib2ran --version`
- Auto-download top result: `lib2ran -b "Book Title" --auto`

---

## 🛠 Requirements

* Python 3.7+
* Dependencies (automatically installed):
  * `libgen-api`
  * `inquirer`
  * `requests`
  * `rich`

---

## 📁 Project Structure

```
lib2ran/
├── __init__.py     # required for the tool to run!
    __main__.py     # Main CLI logic
setup.py            # Installation script
README.md           # This file that you are viewing!
```

---

## 🧠 Future Plans

* Add page navigation in terminal
* Option to auto-download the top result
* GUI version (Tkinter or Electron)

---

## ⚠️ Legal Disclaimer

This tool is intended for **educational purposes** only. Please ensure that you comply with your country's copyright laws. The author is **not responsible** for misuse of this tool.

---

## 🧑‍💻 Author

**Ranbir**
