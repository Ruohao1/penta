package app

import (
	"fmt"
	"time"

	"github.com/Ruohao1/penta/internal/engine"
	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/sinks"
	"github.com/Ruohao1/penta/internal/sinks/summary"
	"github.com/Ruohao1/penta/internal/stages"
	"github.com/Ruohao1/penta/internal/stages/host_discovery"
	"github.com/spf13/cobra"
)

func NewScanCmd() *cobra.Command {
	var opts model.RunOptions
	cmd := &cobra.Command{
		Use:              "scan",
		Short:            "scan targets",
		SilenceUsage:     true,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// var err error
			// task, err = model.NewScanTask(args[0], portsExpr)
			// if err != nil {
			// 	return err
			// }
			//
			// evCh := engine.New(&opts).Run(cmd.Context(), task)
			// count := 0
			//
			// for ev := range evCh {
			//
			// 	if ev.Finding == nil || ev.Finding.Host == nil {
			// 		continue
			// 	}
			//
			// 	if ev.Finding.Host.State != model.HostStateUp {
			// 		continue
			// 	}
			//
			// 	count++
			// 	fmt.Println(ev.Finding.Host.Addr, ev.Finding)
			// }
			// fmt.Println(count)
			return nil
		},
	}

	cmd.PersistentFlags().IntVarP(&opts.Limits.MaxInFlight, "concurrency", "c", 400, "max concurrent operations (global)")
	cmd.PersistentFlags().IntVarP(&opts.Limits.MaxInFlightPerHost, "concurrency-per-host", "H", 4, "max concurrent operations per host")
	cmd.PersistentFlags().IntVarP(&opts.Limits.MinRate, "min-rate", "m", 0, "minimum rate (ops/s), 0 disables")
	cmd.PersistentFlags().IntVarP(&opts.Limits.MaxRate, "max-rate", "M", 200, "maximum rate (ops/s), 0 disables")
	cmd.PersistentFlags().IntVarP(&opts.Limits.MaxRetries, "max-retries", "r", 0, "max retries per operation")
	cmd.PersistentFlags().DurationVarP(&opts.Timeouts.Overall, "timeout", "t", 5*time.Second, "overall operation timeout (0 disables)")
	cmd.PersistentFlags().DurationVar(&opts.Timeouts.TCP, "timeout-tcp", 1500*time.Millisecond, "TCP connect timeout")

	cmd.AddCommand(newScanHostsCmd(&opts))
	cmd.AddCommand(newScanPortsCmd(&opts))
	return cmd
}

func newScanHostsCmd(opts *model.RunOptions) *cobra.Command {
	var probeMethods []string
	var portsExpr []string
	var task model.Task

	cmd := &cobra.Command{
		Use:          "hosts",
		Short:        "Host discovery",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// err := cmd.Parent().PreRunE(cmd, args)
			// if err != nil {
			// 	return err
			// }
			var err error
			if len(probeMethods) != 0 {
				for _, method := range probeMethods {
					switch method {
					case "tcp":
						opts.ProbeOpts.TCP = true
					case "icmp":
						opts.ProbeOpts.ICMP = true
					case "arp":
						opts.ProbeOpts.ARP = true
					default:
						return fmt.Errorf("unknown probe method %q", method)
					}
				}
			}

			task, err = model.NewHostDiscoveryTask(args[0], portsExpr)

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			human, _ := cmd.InheritedFlags().GetBool("human")
			sink := sinks.NewSink(sinks.SinkOptions{
				JSON: !human,
				Out:  cmd.OutOrStdout(),
				Err:  cmd.ErrOrStderr(),
			})
			sink = summary.NewHostDiscoverySink(sink, cmd.OutOrStdout())

			verbose, _ := cmd.InheritedFlags().GetCount("verbose")
			var filterFunc sinks.FilterFunc
			switch {
			case verbose >= 2:
				filterFunc = func(model.Event) bool { return true }
			default:
				filterFunc = sinks.OnlyOpen
			}
			sink = sinks.NewFilterSink(sink, filterFunc)

			hostDiscoveryStage := host_discovery.New()
			eng := engine.Engine{
				Stages: []stages.Stage{hostDiscoveryStage},
				Pool:   engine.DefaultPool,
				Sink:   sink,
			}
			err := eng.Run(ctx, task, *opts)
			sink.Close()
			return err
		},
	}

	cmd.PersistentFlags().StringSliceVarP(&probeMethods, "methods", "P", []string{"arp", "icmp", "tcp"}, "methods use to probe")
	cmd.PersistentFlags().StringSliceVarP(&portsExpr, "ports", "p", []string{"22", "80", "443"}, "tcp probe ports")

	return cmd
}

func newScanPortsCmd(opts *model.RunOptions) *cobra.Command {
	var portsExpr []string
	var nmap bool
	var task model.Task
	cmd := &cobra.Command{
		Use:          "ports",
		Short:        "scan ports",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Parent().PreRunE(cmd, args)
			if err != nil {
				return err
			}
			task, err = model.NewPortScanTask(args[0], portsExpr)
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = task
			// evCh := engine.New(opts).Run(cmd.Context(), task)
			// count := 0
			//
			// for ev := range evCh {
			//
			// 	if ev.Finding == nil || ev.Finding.Host == nil {
			// 		continue
			// 	}
			//
			// 	if ev.Finding.Host.State != model.HostStateUp {
			// 		continue
			// 	}
			//
			// 	count++
			// 	fmt.Println(ev.Finding.Host.Addr, ev.Finding)
			// }
			// fmt.Println(count)
			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(&nmap, "nmap", false, "use nmap to scan")
	cmd.PersistentFlags().StringSliceVarP(&portsExpr, "ports", "p", []string{"22", "80", "443"}, "ports to probe with tcp")

	return cmd
}
