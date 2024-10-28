package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hanapedia/metrics-processor/internal/application/core"
	"github.com/hanapedia/metrics-processor/internal/application/usecases"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/config"
	"github.com/spf13/cobra"
)

// hexagonRequeryCmd represents the hexagon requery command
var hexagonRequeryCmd = &cobra.Command{
	Use:   "hexagon-requery",
	Short: "Requery Hexagon metrics",
	Run: func(cmd *cobra.Command, args []string) {

		config := config.NewConfigFromEnv()
		readS3Adapter := usecases.NewS3Adapter(config)
		// extract endtime for this query using s3 bucket_dir provided
		end, err := readS3Adapter.ParseEndTime()
		if err != nil {
			slog.Error("Error parsing end time.", "error", err)
			os.Exit(1)
		}
		config.EndTime = parseTimestampWithPastCheck(end)

		// replace target dir
		parts := strings.Split(config.S3BucketDir, string(filepath.Separator))

		// add suffix to the parent directory
		if len(parts) > 1 {
			parts[0] = fmt.Sprintf("%s-requery", parts[0])
		}

		config.S3BucketDir = filepath.Join(parts...)

		// recreate s3 adapter with updated config
		writeS3Adapter := usecases.NewS3Adapter(config)

		// create subset query
		prometheusAdapter := usecases.SubsetPrometheusQueryAdapter(config)

		processor := core.NewMetricsProcessor(prometheusAdapter, writeS3Adapter)
		processor.Process()
	},
}

// hexagonRequeryDryCmd represents the hexagon dry command
var hexagonRequeryDryCmd = &cobra.Command{
	Use:   "hexagon-requery-dry",
	Short: "View Queries for requery of Hexagon metrics",
	Run: func(cmd *cobra.Command, args []string) {

		config := config.NewConfigFromEnv()
		prometheusAdapter := usecases.SubsetPrometheusQueryAdapter(config)
		prometheusAdapter.PrintQuery()
	},
}

func init() {
	rootCmd.AddCommand(hexagonRequeryCmd)
	rootCmd.AddCommand(hexagonRequeryDryCmd)
}

// Detects and parses the old timestamp used prior to v2.0.4
// the bug divided the end timestamp with 10e3 instead of 10e2
func parseTimestampWithPastCheck(end float64) time.Time {
    // Parse the timestamp into seconds and nanoseconds
    seconds := int64(end)
    nanoSeconds := int64((end - float64(seconds)) * 1e9)
    parsedTime := time.Unix(seconds, nanoSeconds)

    // Check if the parsed time is way in the past (5 years ago or earlier)
    fiftyYearsAgo := time.Now().AddDate(-5, 0, 0)
    if parsedTime.Before(fiftyYearsAgo) {
        // If it's way in the past, the timestamp might be off by a factor of 10, so correct it
        end *= 10
        seconds = int64(end)
        nanoSeconds = int64((end - float64(seconds)) * 1e9)
        parsedTime = time.Unix(seconds, nanoSeconds)
    }

    return parsedTime
}
