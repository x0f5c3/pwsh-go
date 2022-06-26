package cmd

import (
	"crypto/sha256"
	"fmt"
	"github.com/x0f5c3/pwsh-go/pkg"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/pterm/pcli"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "pwsh-go",
	Short:   "pwsh-go is a tool to update your powershell version automatically",
	Version: "v0.0.4", // <---VERSION---> Updating this version, will also create a new GitHub release.
	RunE: func(cmd *cobra.Command, args []string) error {
		// Your code here
		pterm.Debug.Printf("File extension: %s\n", pkg.FileExt)
		rels, err := pkg.GetReleases()
		if err != nil {
			return err
		}
		parsed, err := rels.Parse()
		if err != nil {
			return err
		}
		rel, err := pkg.AskForVersion(parsed)
		if err != nil {
			return err
		}
		dl, err := rel.Download()
		if err != nil {
			return err
		}
		err = dl.Data.CompareSha()
		if err != nil {
			return err
		}
		err = dl.Data.Save(fmt.Sprintf("./pwsh.%s", pkg.FileExt))
		if err != nil {
			return err
		}
		pterm.Info.Printf("SHA256: %s\n", dl.SHASum)
		pterm.Info.Printf("Version: %s\n", dl.Version)
		b, err := ioutil.ReadFile(fmt.Sprintf("./pwsh.%s", pkg.FileExt))
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		pterm.Info.Printf("Computed SHA256: %x\n", sum)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Fetch user interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pterm.Warning.Println("user interrupt")
		handleUpdate()
		os.Exit(0)
	}()

	// Execute cobra
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.PrintOnError(err)
		handleUpdate()
	}

	handleUpdate()
}

func handleUpdate() {
	err := pcli.CheckForUpdates()
	if err != nil {
		pterm.Error.Printf("Failed to check for updates %s\n", err)
		os.Exit(1)
	}
}

var interactive bool

func init() {
	// Adds global flags for PTerm settings.
	// Fill the empty strings with the shorthand variant (if you like to have one).
	rootCmd.PersistentFlags().BoolVarP(&pterm.PrintDebugMessages, "debug", "", false, "enable debug messages")
	rootCmd.PersistentFlags().BoolVarP(&pterm.RawOutput, "raw", "", false, "print unstyled raw output (set it if output is written to a file)")
	rootCmd.PersistentFlags().BoolVarP(&pcli.DisableUpdateChecking, "disable-update-checks", "", false, "disables update checks")
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "Choose the PowerShell version interactively")

	// Use https://github.com/pterm/pcli to style the output of cobra.
	err := pcli.SetRepo("x0f5c3/pwsh-go")
	if err != nil {
		pterm.Fatal.Printf("Failed to set the repo %s\n", err)
	}
	pcli.SetRootCmd(rootCmd)
	pcli.Setup()
	// Change global PTerm theme
	pterm.ThemeDefault.SectionStyle = *pterm.NewStyle(pterm.FgCyan)
}
