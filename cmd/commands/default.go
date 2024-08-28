package commands

import (
	"github.com/hanapedia/metrics-processor/internal/application/core"
	"github.com/hanapedia/metrics-processor/internal/application/usecases"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/config"
	"github.com/spf13/cobra"
)

// defaultCmd represents the default command
var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Query default metrics",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.NewConfigFromEnv()
		prometheusAdapter := usecases.PrometheusQueryAdapter(config)
		s3Adapter := usecases.NewS3Adapter(config)

		processor := core.NewMetricsProcessor(prometheusAdapter, s3Adapter)
		processor.Process()
	},
}

func init() {
	rootCmd.AddCommand(defaultCmd)
}
