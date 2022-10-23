package cmd

import (
	"context"
	glice "github.com/ribice/glice/v3/pkg"
	"github.com/spf13/cobra"
)

// thankCmd represents the thank command
var thankCmd = &cobra.Command{
	Use:   "thanks",
	Run:   RunThanks,
	Short: "Provide thanks to your dependencies by up-voting them",
	Long:  "Provide thanks to your dependencies by up-voting them (for Github that would be 'starring' them.)",
}

func init() {
	rootCmd.AddCommand(thankCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// thankCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// thankCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//goland:noinspection GoUnusedParameter
func RunThanks(cmd *cobra.Command, args []string) {
	var ra glice.RepositoryAdapter
	var err error

	if !glice.HasGitHubAPIKey() {
		Failf(glice.ExitHostNotYetSupported,
			"\nCannot give thanks without GITHUB_API_KEY environment variable being set; %w",
			glice.ErrNoAPIKey)
	}

	ctx := context.Background()
	deps := ScanDependencies(ctx)

	for _, dep := range deps {
		ra, err = glice.GetRepositoryAdapter(ctx, dep.Repository())
		if err != nil {
			Warnf("\nUnable to get repository adapter for %s; %w",
				dep.Host,
				err)
		}
		err = ra.UpVoteRepository(ctx)
		if err != nil {
			Warnf("\nUnable to upvote repository for '%s'; %w", dep.Import, err)
		}
	}
}
