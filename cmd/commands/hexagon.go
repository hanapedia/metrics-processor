package commands

import (
	"github.com/hanapedia/metrics-processor/internal/application/core"
	"github.com/hanapedia/metrics-processor/internal/application/usecases"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/config"
	"github.com/spf13/cobra"
)

// hexagonCmd represents the hexagon command
var hexagonCmd = &cobra.Command{
	Use:   "hexagon",
	Short: "Query Hexagon metrics",
	Run: func(cmd *cobra.Command, args []string) {

		config := config.NewConfigFromEnv()
		prometheusAdapter := usecases.HexagonPrometheusQueryAdapter(config)
		s3Adapter := usecases.NewS3Adapter(config)

		processor := core.NewMetricsProcessor(prometheusAdapter, s3Adapter)
		processor.Process()
	},
}

// hexagonDryCmd represents the hexagon dry command
var hexagonDryCmd = &cobra.Command{
	Use:   "hexagon-dry",
	Short: "View Queries for Hexagon metrics",
	Run: func(cmd *cobra.Command, args []string) {

		config := config.NewConfigFromEnv()
		prometheusAdapter := usecases.HexagonPrometheusQueryAdapter(config)
		prometheusAdapter.PrintQuery()
	},
}

func init() {
	rootCmd.AddCommand(hexagonCmd)
	rootCmd.AddCommand(hexagonDryCmd)
}
