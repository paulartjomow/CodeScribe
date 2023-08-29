package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Prixix/CodeScribe/internal/snippet"
	"github.com/Prixix/CodeScribe/pkg/clipboard"
	"github.com/Prixix/CodeScribe/pkg/database"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func main() {
	dbPath := "snippets.db"

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

	languages, err := readLanguagesFromJSON("languages.json")
	if err != nil {
		fmt.Println("Error reading languages:", err)
		fmt.Println("Downloading languages from GitHub...")
		downloadLanguagesFromGitHub("languages.json")
		os.Exit(1)
	}

	languages = append([]string{""}, languages...)

	snippetManager := snippet.NewManager(db)

	var rootCmd = &cobra.Command{
		Use:   "CodeScribe",
		Short: "A tool to manage code snippets",
		Long:  "CodeScribe is a command-line tool to help programmers organize and manage their code snippets.",
		Run: func(cmd *cobra.Command, args []string) {
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
				AddDropDown("Language", languages, 0, nil).
				AddButton("Save", func() {
					app.Stop()

				}).AddButton("Cancel", func() {
				app.Stop()
			})

			form.SetBorder(true).SetTitle("Create Snippet").SetTitleAlign(tview.AlignLeft)

			if err := app.SetRoot(form, true).Run(); err != nil {
				fmt.Println("Error:", err)
			}

			title := form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			tags := form.GetFormItemByLabel("Tags").(*tview.InputField).GetText()
			code := form.GetFormItemByLabel("Code").(*tview.TextArea).GetText()
			_, language := form.GetFormItemByLabel("Language").(*tview.DropDown).GetCurrentOption()

			if title == "" || code == "" {
				fmt.Println("Title and/or code cannot be empty!")
				return
			}

			err1 := snippetManager.CreateSnippet(title, description, tags, code, language)
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

	var editCmd = &cobra.Command{
		Use:   "edit [snippet-id]",
		Short: "Edit a snippet",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			snippetIDInt, err := strconv.Atoi(args[0])

			if err != nil {
				fmt.Println("Invalid snippet ID:", err)
				return
			}

			snippet, err := snippetManager.GetSnippetByID(snippetIDInt)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			app := tview.NewApplication()

			form := tview.NewForm().
				AddInputField("Title", snippet.Title, 20, nil, nil).
				AddInputField("Description", snippet.Description, 20, nil, nil).
				AddInputField("Tags", snippet.Tags, 20, nil, nil).
				AddTextArea("Code", snippet.Code, 0, 10, 0, nil).
				AddDropDown("Language", languages, 0, nil).
				AddButton("Save", func() {
					app.Stop()

				}).AddButton("Cancel", func() {
				app.Stop()
			})

			form.SetBorder(true).SetTitle("Edit Snippet").SetTitleAlign(tview.AlignLeft)

			if err := app.SetRoot(form, true).Run(); err != nil {
				fmt.Println("Error:", err)
			}

			title := form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			tags := form.GetFormItemByLabel("Tags").(*tview.InputField).GetText()
			code := form.GetFormItemByLabel("Code").(*tview.TextArea).GetText()
			_, language := form.GetFormItemByLabel("Language").(*tview.DropDown).GetCurrentOption()

			if title == "" || code == "" {
				fmt.Println("Title and/or code cannot be empty!")
				return
			}

			err1 := snippetManager.UpdateSnippet(snippetIDInt, title, description, tags, code, language)
			if err1 != nil {
				fmt.Println("Error updating snippet:", err)
			} else {
				fmt.Println("Snippet updated successfully!")
			}
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all saved snippets with IDs",
		Run: func(cmd *cobra.Command, args []string) {
			lang, _ := cmd.Flags().GetString("language")

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
				if lang != "" && !strings.EqualFold(lang, s.Language) {
					continue
				}

				language := s.Language
				if language != "" {
					language = fmt.Sprintf("(%s)", language)
				}

				list.AddItem(fmt.Sprintf("[blue][%d] [white]%s", s.ID, s.Title),
					fmt.Sprintf("[grey]%s [yellow]%s", language, s.Description), ' ', func() {
						if err := clipboard.CopyToClipboard(s.Code); err != nil {
							fmt.Println("Error copying to clipboard:", err)
						}
						app.Stop()
					})
			}

			list.AddItem("New", "Create a new snippet", 'n', func() {
				app.Stop()
				createCmd.Run(cmd, args)
			})

			list.AddItem("Back", "Return to main menu", 'q', func() {
				app.Stop()
			})

			if err := app.SetRoot(list, true).Run(); err != nil {
				fmt.Println("Error:", err)
			}
		},
	}
	listCmd.Flags().String("language", "", "Filter snippets by programming language")

	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func downloadLanguagesFromGitHub(s string) {
	url := "https://raw.githubusercontent.com/Prixix/CodeScribe/main/languages.json"
	if err := downloadFile(s, url); err != nil {
		fmt.Println("Error downloading file:", err)
		os.Exit(1)
	}
}

func downloadFile(filename string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

func readLanguagesFromJSON(filename string) ([]string, error) {
	var languages []string

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	return languages, nil
}
