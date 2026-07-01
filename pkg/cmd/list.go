package cmd

import (
	"fmt"
	"os"

	"github.com/rayzalzero/outline-cli/pkg/api"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accessible collections",
	Long:  `List all collections you have access to with their IDs and names.`,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// Get base URL from environment or default
	baseURL := os.Getenv("OUTLINE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://outline-rbi.jatismobile.com"
	}

	// Get token from environment (priority: OUTLINE_API_KEY > OUTLINE_TOKEN)
	token := os.Getenv("OUTLINE_API_KEY")
	if token == "" {
		token = os.Getenv("OUTLINE_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("authentication required: set OUTLINE_API_KEY or OUTLINE_TOKEN environment variable")
	}

	// Initialize API client
	client := api.NewClient(baseURL, token)

	// List collections
	fmt.Println("Fetching collections...")
	collections, err := client.ListCollections()
	if err != nil {
		return fmt.Errorf("list collections: %w", err)
	}

	if len(collections) == 0 {
		fmt.Println("No collections found.")
		return nil
	}

	// Print collections
	fmt.Printf("\nFound %d collection(s):\n\n", len(collections))
	for i, coll := range collections {
		fmt.Printf("[%d] %s\n", i+1, coll.Name)
		fmt.Printf("    ID:  %s\n", coll.ID)
		if coll.Description != "" {
			fmt.Printf("    Desc: %s\n", coll.Description)
		}
		fmt.Printf("    URL: %s\n", coll.URL)
		fmt.Println()
	}

	fmt.Println("Usage: outline clone <collection-id> <directory>")
	return nil
}
