package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wenwenxiong/host-prometheus/cmd/app/options"
	"github.com/wenwenxiong/host-prometheus/pkg/apiserver"
	"github.com/wenwenxiong/host-prometheus/pkg/utils/signals"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"
	"runtime"
)

var (
	version     bool
	endpoint   string
	listenPort  string
	logger      = logrus.New()
	// Version shows the host-prometheus binary version.
	Version string
	// GitSHA shows the  host-prometheus binary code commit SHA on git.
	GitSHA string
)

func printVersionInfo() {
	logger.Infof("host-prometheus Version: %s", Version)
	logger.Infof("Git SHA: %s", GitSHA)
	logger.Infof("Go Version: %s", runtime.Version())
	logger.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
}

func NewPrometheusCommand() *cobra.Command {
	s := options.NewServerRunOptions()

	// Load configuration from file
	conf, err := apiserver.TryLoadFromDisk()
	if err == nil {
		s = &options.ServerRunOptions{
			ListenPort:  options.DefaultListenPort,
			Config:                  conf,
		}
	}
	var RootCmd = &cobra.Command{
		Use:   "host-prometheus",
		Short: "host-prometheus serve for provide host monitoring metrics value ",
		Long: `The host-prometheus,mini web server, is serve to provide host monitoring metrics value from specify prometheus, it encapsulation metrics for cpu memory disk network usage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if errs := s.Validate(); len(errs) != 0 {
				return utilerrors.NewAggregate(errs)
			} else {
				return  RunServer(listenPort, endpoint, s, signals.SetupSignalHandler())
			}
		},
		SilenceUsage: true,
	}
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "print version info")
	RootCmd.PersistentFlags().StringVarP(&listenPort, "listenPort", "p", "9111", "listen port, default is 9111")
	RootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:9090", "prometheus endpoint, default is http://localhost:9090")
	fs := RootCmd.Flags()
	namedFlagSets := s.Flags()
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(RootCmd.OutOrStdout())
	RootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
	return RootCmd
}

func RunServer(listenPort string, endpoint string, s *options.ServerRunOptions, stopCh <-chan struct{}) error {
	server, err := options.NewApiServer(listenPort, endpoint, s, stopCh)
	if err != nil {
		return err
	}

	return server.Run()
}