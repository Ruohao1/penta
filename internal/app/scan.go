package app

import (
	"time"

	"github.com/Ruohao1/penta/internal/scan"
	"github.com/Ruohao1/penta/internal/scan/native"
	"github.com/Ruohao1/penta/internal/scan/nmap"
	"github.com/Ruohao1/penta/internal/utils"
	"github.com/spf13/cobra"
)

func NewScanCmd() *cobra.Command {
	var useNmap bool
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan a target",
	}

	cmd.AddCommand(newScanHostsCmd())
	// cmd.AddCommand(newScanPortsCmd())
	// cmd.AddCommand(newScanWebCmd())
	cmd.Flags().BoolVar(&useNmap, "nmap", false, "use nmap")
	return cmd
}

func newScanHostsCmd() *cobra.Command {
	var opts scan.HostsOptions
	var methodStrings []string

	cmd := &cobra.Command{
		Use:   "hosts",
		Short: "Scan hosts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := utils.LoggerFrom(ctx)

			var err error
			opts.Methods, err = scan.ParseMethods(methodStrings)
			if err != nil {
				return err
			}

			var engine scan.Engine
			switch opts.EngineName {
			case "nmap":
				engine = nmap.NewEngine("nmap", nil)
			default:
				engine = native.NewEngine(1024, 3*time.Second)
			}

			svc := scan.NewService(engine)
			logger = logger.With().Str("component", "scan").Logger()
			ctx = utils.WithLogger(ctx, logger)

			svc.ScanHosts(ctx, args[0], opts)

			return nil
		},
	}

	opts.EngineName = scan.NativeEngine
	cmd.Flags().Var(&opts.EngineName, "engine", "scan engine: native|nmap")
	cmd.Flags().StringSliceVarP(&methodStrings, "methods", "m", methodStrings, "scan methods: arp,icmp,tcp")
	cmd.Flags().DurationVar(&opts.Timeout, "timeout", 500*time.Millisecond, "scan timeout")

	return cmd
}
