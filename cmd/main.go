package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Prixix/CodeScribe/internal/snippet"
	"github.com/Prixix/CodeScribe/pkg/clipboard"
	"github.com/Prixix/CodeScribe/pkg/database"
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
			// Get input from the user for snippet details
			var title, description, tags, code string
			fmt.Print("Enter snippet title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter snippet description: ")
			fmt.Scanln(&description)
			fmt.Print("Enter snippet tags: ")
			fmt.Scanln(&tags)
			fmt.Print("Enter code snippet: ")
			fmt.Scanln(&code)

			// Create the snippet in the database
			err := snippetManager.CreateSnippet(title, description, tags, code)
			if err != nil {
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

	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(createCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
