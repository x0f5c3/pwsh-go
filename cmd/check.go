package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/x0f5c3/pwsh-go/pkg"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if there's an update for your version",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := pkg.GetLocalVersion()
		if err != nil {
			return err
		}
		latest, err := pkg.GetLatest()
		if err != nil {
			return err
		}
		if v.LessThan(*latest.Version) {
			return errors.New("you have an earlier version")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
