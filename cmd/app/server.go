package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wenwenxiong/host-prometheus/pkg/apiserver"
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
	var RootCmd = &cobra.Command{
		Use:   "host-prometheus",
		Short: "host-prometheus serve for provide host monitoring metrics value ",
		Long: `The host-prometheus,mini web server, is serve to provide host monitoring metrics value from specify prometheus, it encapsulation metrics for cpu memory disk network usage.`,
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				printVersionInfo()
			} else {
				RunServer(listenPort, endpoint)
			}
		},
	}
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "print version info")
	RootCmd.PersistentFlags().StringVarP(&listenPort, "listenPort", "p", "9111", "listen port, default is 9111")
	RootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:9090", "prometheus endpoint, default is http://localhost:9090")
	return RootCmd

}

func RunServer(listenPort string, endpoint string) error {
	apiserver, err := apiserver.NewApiServer(listenPort, endpoint)
	if err != nil {
		return err
	}

	return apiserver.Run()
}