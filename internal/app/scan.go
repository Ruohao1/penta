package app

import (
	"fmt"
	"os"
	"time"

	"github.com/Ruohao1/penta/internal/engine"
	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/targets"
	"github.com/spf13/cobra"
)

func NewScanCmd() *cobra.Command {
	var opts model.RunOptions
	var req model.Request
	var nmap bool
	cmd := &cobra.Command{
		Use:              "scan",
		Short:            "scan targets",
		SilenceUsage:     true,
		TraverseChildren: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
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
				req.Backend = model.BackendInternal
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			targetList, err := targets.Resolve(args[0], model.TargetTypeHost)
			opts.Targets = targetList

			return err
		},
	}

	cmd.PersistentFlags().BoolVar(&nmap, "nmap", false, "use nmap to scan")
	cmd.PersistentFlags().IntVarP(&opts.Concurrency, "concurrency", "c", 100, "concurrency")
	cmd.PersistentFlags().IntVarP(&opts.MinRate, "min-rate", "m", 0, "min rate")
	cmd.PersistentFlags().IntVarP(&opts.MaxRate, "max-rate", "M", 0, "max rate")
	cmd.PersistentFlags().IntVarP(&opts.MaxRetries, "max-retries", "r", 1, "max retries")
	cmd.PersistentFlags().DurationVarP(&opts.Timeout, "timeout", "t", 800*time.Millisecond, "timeout")

	cmd.AddCommand(newScanHostsCmd(&req, &opts))
	return cmd
}

func newScanHostsCmd(req *model.Request, opts *model.RunOptions) *cobra.Command {
	var probeMethods []string
	cmd := &cobra.Command{
		Use:          "hosts",
		Short:        "scan hosts",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Parent().PreRunE(cmd, args)
			if err != nil {
				return err
			}
			if len(probeMethods) != 0 {
				for _, method := range probeMethods {
					switch method {
					case "tcp":
						opts.TCP = true
					case "icmp":
						opts.ICMP = true
					case "arp":
						opts.ARP = true
					default:
						return fmt.Errorf("unknown probe method %q", method)
					}
				}
			}
			req.Mode = model.ModeHosts
			targetList, err := targets.Resolve(args[0], model.TargetTypeHost)
			if err != nil {
				return err
			}
			opts.Targets = targetList
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			evCh := engine.New(opts).Run(cmd.Context(), *req)
			count := 0

			for ev := range evCh {

				if ev.Finding == nil || ev.Finding.Host == nil {
					continue
				}

				if ev.Finding.Host.State != model.HostStateUp {
					continue
				}

				count++
				fmt.Println(ev.Finding.Host.Addr, ev.Finding)
			}
			fmt.Println(count)
			return nil
		},
	}

	cmd.PersistentFlags().StringSliceVarP(&probeMethods, "methods", "P", []string{"arp", "icmp", "tcp"}, "probeMethods")

	return cmd
}
