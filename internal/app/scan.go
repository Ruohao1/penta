package app

import (
	"encoding/json"
	"math"
	"os"
	"time"

	"github.com/Ruohao1/penta/internal/engine"
	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/targets"
	"github.com/spf13/cobra"
)

func NewScanCmd() *cobra.Command {
	var opts engine.RunOptions
	var req model.Request
	var nmap bool
	cmd := &cobra.Command{
		Use:          "scan",
		Short:        "scan targets",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if nmap {
				req.Backend = model.BackendNmap
				argv := os.Args
				pre, post := targets.SplitDashDash(argv)

				targetExpr := targets.ExtractTargets(pre, cmd.Flags()) // or cmd.InheritedFlags()/PersistentFlags as needed
				nmapArgs := post

				// If targets were provided before --, append them to nmap args (unless user already included them)
				if len(targetExpr) > 0 {
					nmapArgs = append(nmapArgs, targetExpr)
				}
				req.ToolArgs = nmapArgs

			} else {
				targetList, err := targets.Resolve(args[0], targets.TargetTypeHost)
				opts.Targets = targetList
				return err

			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(&nmap, "nmap", false, "use nmap to scan")
	cmd.PersistentFlags().IntVarP(&opts.Concurrency, "concurrency", "c", 10, "concurrency")
	cmd.PersistentFlags().IntVarP(&opts.MinRate, "min-rate", "m", 0, "min rate")
	cmd.PersistentFlags().IntVarP(&opts.MaxRate, "max-rate", "M", math.MaxInt, "max rate")
	cmd.PersistentFlags().IntVarP(&opts.MaxRetries, "max-retries", "r", 3, "max retries")
	cmd.PersistentFlags().DurationVarP(&opts.Timeout, "timeout", "t", 3*time.Second, "timeout")

	cmd.AddCommand(newScanHostsCmd(&req, &opts))
	return cmd
}

func newScanHostsCmd(req *model.Request, opts *engine.RunOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "hosts",
		Short:        "scan hosts",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			req.Mode = model.ModeHosts
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			evCh := engine.New(opts).Run(cmd.Context(), *req)

			enc := json.NewEncoder(os.Stdout)
			for ev := range evCh {
				if err := enc.Encode(ev); err != nil {
					return err
				}
			}
			return nil
		},
	}

	return cmd
}
