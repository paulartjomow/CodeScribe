package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Prixix/CodeScribe/internal/snippet"
	"github.com/Prixix/CodeScribe/pkg/clipboard"
	"github.com/Prixix/CodeScribe/pkg/database"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func main() {
	dbPath := "snippets.db" // Adjust the database path as needed

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		fmt.Println("Error initializing the database:", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.InitializeSchema(dbPath); err != nil {
		fmt.Println("Error initializing the database schema:", err)
		os.Exit(1)
	}

	snippetManager := snippet.NewManager(db)

	var rootCmd = &cobra.Command{
		Use:   "CodeScribe",
		Short: "A tool to manage code snippets",
		Long:  "CodeScribe is a command-line tool to help programmers organize and manage their code snippets.",
		Run: func(cmd *cobra.Command, args []string) {
			// Display a help message
			cmd.Help()
		},
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new code snippet",
		Run: func(cmd *cobra.Command, args []string) {
			dbPath := "snippets.db"
			db, err := database.NewDatabase(dbPath)
			if err != nil {
				fmt.Println("Error initializing the database:", err)
				os.Exit(1)
			}
			defer db.Close()

			snippetManager := snippet.NewManager(db)

			app := tview.NewApplication()

			form := tview.NewForm().
				AddInputField("Title", "", 20, nil, nil).
				AddInputField("Description", "", 20, nil, nil).
				AddInputField("Tags", "", 20, nil, nil).
				AddTextArea("Code", "", 0, 10, 0, nil).
				AddDropDown("Language", []string{"", "Go", "Python", "Java", "JavaScript", "C++", "C#", "C", "PHP", "Ruby", "Unspecified"}, 0, nil).
				AddButton("Save", func() {
					app.Stop()

				}).AddButton("Cancel", func() {
				app.Stop()
			})

			form.SetBorder(true).SetTitle("Create Snippet").SetTitleAlign(tview.AlignLeft)

			if err := app.SetRoot(form, true).Run(); err != nil {
				fmt.Println("Error:", err)
			}

			// Create snippet in snippetmanager
			title := form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			tags := form.GetFormItemByLabel("Tags").(*tview.InputField).GetText()
			code := form.GetFormItemByLabel("Code").(*tview.TextArea).GetText()

			// Check if the title is empty or code is empty
			if title == "" || code == "" {
				fmt.Println("Title and/or code cannot be empty!")
				return
			}

			// Create the snippet
			err1 := snippetManager.CreateSnippet(title, description, tags, code)
			if err1 != nil {
				fmt.Println("Error creating snippet:", err)
			} else {
				fmt.Println("Snippet created successfully!")
			}

		},
	}

	var copyCmd = &cobra.Command{
		Use:   "copy [snippet-id]",
		Short: "Copy a snippet to the clipboard",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			snippetID := args[0]
			// Convert the snippetID to an integer
			snippetIDInt, err := strconv.Atoi(snippetID)
			if err != nil {
				fmt.Println("Invalid snippet ID:", err)
				return
			}

			snippet, err := snippetManager.GetSnippetByID(snippetIDInt)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			if err := clipboard.CopyToClipboard(snippet.Code); err != nil {
				fmt.Println("Error copying to clipboard:", err)
			} else {
				fmt.Println("Snippet copied to clipboard!")
			}
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all saved snippets with IDs",
		Run: func(cmd *cobra.Command, args []string) {
			snippets, err := snippetManager.GetAllSnippets()
			if err != nil {
				fmt.Println("Error fetching snippets:", err)
				return
			}

			if len(snippets) == 0 {
				fmt.Println("No snippets found.")
				return
			}

			app := tview.NewApplication()

			list := tview.NewList().
				ShowSecondaryText(true)

			for _, s := range snippets {
				// Add the snippet and its ID in the brackets at the beginning
				list.AddItem(fmt.Sprintf("[%d] %s", s.ID, s.Title), s.Description, ' ', func() {
					if err := clipboard.CopyToClipboard(s.Code); err != nil {
						fmt.Println("Error copying to clipboard:", err)
					} else {
						fmt.Println("Snippet copied to clipboard!")
					}
					app.Stop()
				})
			}

			list.AddItem("Back", "Return to main menu", 'q', func() {
				app.Stop()
			})

			if err := app.SetRoot(list, true).Run(); err != nil {
				fmt.Println("Error:", err)
			}
		},
	}

	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
