package app

import (
	// "github.com/Ruohao1/penta/internal/config"
	// "github.com/Ruohao1/penta/internal/utils"

	"github.com/Ruohao1/penta/internal/config"
	"github.com/Ruohao1/penta/internal/ui"
	"github.com/Ruohao1/penta/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var rootCmd = newRootCmd()

func Execute() error {
	return rootCmd.Execute()
}

func newRootCmd() *cobra.Command {
	var opts globalOptions
	cmd := &cobra.Command{
		Use:          "penta",
		Short:        "Ultimate pentest CLI engine",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 1) Map -v count -> log level
			var lvl zerolog.Level
			switch {
			case opts.verbosity >= 3:
				lvl = zerolog.TraceLevel // -vvv and beyond
			case opts.verbosity == 2:
				lvl = zerolog.DebugLevel // -vv
			case opts.verbosity == 1:
				lvl = zerolog.InfoLevel // -v
			default:
				lvl = zerolog.WarnLevel // no -v
			}

			// 2) Load config
			cfg := config.LoadConfig()

			// 3) Build logger with level + human mode
			logger := utils.NewLogger(opts.human, lvl).
				With().
				Str("cmd", cmd.Name()).
				Logger()
			zerolog.SetGlobalLevel(lvl)

			// 4) Attach to context
			ctx := cmd.Context()
			ctx = utils.WithLogger(ctx, logger)
			ctx = utils.WithConfig(ctx, cfg)
			cmd.SetContext(ctx)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return ui.RunTUI(ctx, ui.TuiOptions{})
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.human, "human", true, "human-friendly log output")
	cmd.PersistentFlags().CountVarP(&opts.verbosity, "verbose", "v", "increase verbosity (-v, -vv, -vvv)")
	cmd.PersistentFlags().BoolVar(&opts.tui, "tui mode", true, "use tui mode")
	return cmd
}

func init() {
	// rootCmd.AddCommand(NewSessionCmd())
	rootCmd.AddCommand(NewScanCmd())
	// rootCmd.AddCommand(NewBruteCmd())
}
