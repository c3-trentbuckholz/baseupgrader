package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "baseupgrader",
		Short: "Find the files to update during base app upgrades",
		RunE: func(cmd *cobra.Command, args []string) error {
			authToken, _ := cmd.Flags().GetString("ghAuthToken")
			baseRepoUrl, _ := cmd.Flags().GetString("baseRepoUrl")
			targetRepoUrl, _ := cmd.Flags().GetString("targetRepoUrl")
			targetRepoPath, _ := cmd.Flags().GetString("targetRepoPath")
			oldCommit, _ := cmd.Flags().GetString("oldCommit")
			newCommit, _ := cmd.Flags().GetString("newCommit")

			gitClientBase := git.NewGitClient(baseRepoUrl, authToken, http.DefaultClient)
			gitClientTarget := git.NewGitClient(targetRepoUrl, authToken, http.DefaultClient)

			var (
				errGroup    errgroup.Group
				baseDiff    git.Diff
				targetFiles []string
			)

			errGroup.Go(func() error {
				var err error
				baseDiff, err = gitClientBase.GetDiff(oldCommit, newCommit)
				return err
			})
			errGroup.Go(func() error {
				var err error
				targetFiles, err = gitClientTarget.GetFiles(targetRepoPath)
				return err
			})
			if err := errGroup.Wait(); err != nil {
				log.Panic(err)
			}
			for _, fileName := range targetFiles {
				if _, ok := baseDiff[fileName]; ok {
					fmt.Println(fileName)
					fmt.Println("--------------------------------------------")
					fmt.Println(baseDiff[fileName])
					fmt.Println("============================================")
				}
			}

			return nil
		},
	}
}

func main() {
	var (
		ghAuthToken    string
		baseRepoUrl    string
		targetRepoUrl  string
		targetRepoPath string
		oldCommit      string
		newCommit      string
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
	rootCmd.Flags().StringVar(&targetRepoPath, "targetRepoPath", "",
		"Path to directory in target repository")
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
