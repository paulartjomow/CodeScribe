# CodeScribe: Go Snippet Organizer

CodeScribe is a command-line tool built with Go that makes managing and organizing your code snippets a breeze. It allows programmers to store, search, and retrieve code snippets quickly and efficiently, helping streamline your development workflow.

## Features

- **Snippet Storage:** Add, edit, and delete code snippets with titles, descriptions, tags, and code content.
- **Search and Filtering:** Find snippets by title, description, tags, or filter based on programming languages and frameworks.
- **Syntax Highlighting:** Code snippets are displayed with syntax highlighting for improved readability.
- **Clipboard Integration:** Copy selected snippets to the clipboard for easy pasting.
- **Backup and Restore:** Backup your snippet database to prevent data loss and restore when needed.
- **Sharing:** Share your favorite snippets with colleagues and the programming community.
- **Import/Export:** Import and export snippets in various formats, making migration a breeze.
- **Version Control Integration:** Seamlessly integrate with version control systems like Git.

## Installation

1. Clone this repository: `git clone https://github.com/Prixix/CodeScribe.git`
2. Navigate to the project directory: `cd CodeScribe`
3. Build the project: `go build cmd/main.gp`
4. Rename & Run CodeScribe

## Usage

CodeScribe offers an intuitive command-line interface for managing snippets. Here are some example commands:

- To add a new snippet: `./CodeScribe create`
- To copy a snippet to the clipboard: `./CodeScribe copy [id]`
- To show a list of all snippets: `./CodeScribe list`
- For more commands and options, refer to the [User Guide](docs/user-guide.md).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Built with ❤️ by [Prixix](https://github.com/Prixix)
