package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "baseupgrader",
		Short: "Find the files to update during base app upgrades",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseRepoUrl, _ := cmd.Flags().GetString("baseRepoUrl")
			targetRepoUrl, _ := cmd.Flags().GetString("targetRepoUrl")
			oldCommit, _ := cmd.Flags().GetString("oldCommit")
			newCommit, _ := cmd.Flags().GetString("newCommit")

			return nil
		},
	}
}

func main() {
	var (
		ghAuthToken   string
		baseRepoUrl   string
		targetRepoUrl string
		oldCommit     string
		newCommit     string
	)

	rootCmd := NewCmd()

	rootCmd.Flags().StringVar(&ghAuthToken, "ghAuthToken", "", "GitHub authentication token")
	if err := rootCmd.MarkFlagRequired("ghAuthToken"); err != nil {
		log.Panic(err)
	}
	rootCmd.Flags().StringVar(&baseRepoUrl, "baseRepoUrl", "", "Base repository URL")
	if err := rootCmd.MarkFlagRequired("baseRepoUrl"); err != nil {
		log.Panic(err)
	}
	rootCmd.Flags().StringVar(&targetRepoUrl, "targetRepoUrl", "", "Target repository URL")
	if err := rootCmd.MarkFlagRequired("targetRepoUrl"); err != nil {
		log.Panic(err)
	}
	rootCmd.Flags().StringVar(&oldCommit, "oldCommit", "", "Old commit hash")
	if err := rootCmd.MarkFlagRequired("oldCommit"); err != nil {
		log.Panic(err)
	}
	rootCmd.Flags().StringVar(&newCommit, "newCommit", "", "New commit hash")
	if err := rootCmd.MarkFlagRequired("newCommit"); err != nil {
		log.Panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}
