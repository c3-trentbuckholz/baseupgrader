package main

import (
	"log"
	"net/http"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
	"github.com/c3-trentbuckholz/baseupgrader/pkg/report"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"golang.org/x/sync/errgroup"
)

type OutputType enumflag.Flag

const (
	Basic OutputType = iota
	Json
	Html
)

var OutputTypeIds = map[OutputType][]string{
	Basic: {"basic", "b"},
	Json:  {"json", "j"},
	Html:  {"html", "h"},
}

func main() {
	var (
		ghAuthToken    string
		baseRepoUrl    string
		targetRepoUrl  string
		targetRepoPath string
		oldCommit      string
		newCommit      string
		outputType     OutputType
	)

	rootCmd := &cobra.Command{
		Use:   "baseupgrader",
		Short: "Find the files to update during base app upgrades",
		RunE: func(cmd *cobra.Command, args []string) error {
			gitClientBase := git.NewGitClient(baseRepoUrl, ghAuthToken, http.DefaultClient)
			gitClientTarget := git.NewGitClient(targetRepoUrl, ghAuthToken, http.DefaultClient)

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

			var reporter report.Reporter
			switch outputType {
			case Basic:
				reporter = report.NewBasicReport(targetFiles, baseDiff)
			case Json:
				reporter = report.NewJsonReport(targetFiles, baseDiff)
			case Html:
				log.Panic("HTML report not yet implemented")
			default:
				log.Panic("Unknown output type")
			}

			reportContents, err := reporter.Create()
			if err != nil {
				log.Panic(err)
			}
			err = reporter.Write(reportContents)
			if err != nil {
				log.Panic(err)
			}

			return nil
		},
	}

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
	rootCmd.Flags().Var(
		enumflag.New(&outputType, "outputType", OutputTypeIds, enumflag.EnumCaseSensitive),
		"outputType",
		"Set the output type (available options: basic, json, html)",
	)

	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}
