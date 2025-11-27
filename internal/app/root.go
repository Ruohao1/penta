package app

import (
	"github.com/Ruohao1/penta/internal/config"
	"github.com/Ruohao1/penta/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	human   bool
	verb    int
	rootCmd = NewRootCmd()
)

func Execute() error {
	return rootCmd.Execute()
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "penta",
		Short:        "Ultimate pentest CLI engine",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 1) Map -v count -> log level
			var lvl zerolog.Level
			switch {
			case verb >= 3:
				lvl = zerolog.TraceLevel // -vvv and beyond
			case verb == 2:
				lvl = zerolog.DebugLevel // -vv
			case verb == 1:
				lvl = zerolog.InfoLevel // -v
			default:
				lvl = zerolog.WarnLevel // no -v
			}

			// 2) Load config
			cfg := config.LoadConfig()

			// 3) Build logger with level + human mode
			logger := utils.NewLogger(human, lvl).
				With().
				Str("cmd", cmd.Name()).
				Logger()

			// optional but useful for any other zerolog loggers
			zerolog.SetGlobalLevel(lvl)

			// 4) Attach to context
			ctx := cmd.Context()
			ctx = utils.WithLogger(ctx, logger)
			ctx = utils.WithConfig(ctx, cfg)
			cmd.SetContext(ctx)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().BoolVar(&human, "human", true, "human-friendly log output")
	cmd.PersistentFlags().CountVarP(&verb, "verbose", "v", "increase verbosity (-v, -vv, -vvv)")
	return cmd
}

func init() {
	// no extra verbose flag here, just subcommands
	rootCmd.AddCommand(NewSessionCmd())
	rootCmd.AddCommand(NewScanCmd())
	// rootCmd.AddCommand(NewBruteCmd())
}
