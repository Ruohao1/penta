package app

import "github.com/spf13/cobra"

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "penta",
	Short: "Ultimate pentest CLI tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(NewSessionCmd())
	// rootCmd.AddCommand(NewScanCmd())
	// rootCmd.AddCommand(NewBruteCmd())
}
