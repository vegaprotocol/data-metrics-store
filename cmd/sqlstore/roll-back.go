package sqlstore

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/vega-monitoring/config"
	"github.com/vegaprotocol/vega-monitoring/sqlstore"
)

type RollBackArgs struct {
	*SQLStoreArgs
}

var rollBackArgs RollBackArgs

// rollBackCmd represents the rollBack command
var rollBackCmd = &cobra.Command{
	Use:   "roll-back",
	Short: "Roll back the last migration for Monitoring Tables",
	Long:  `Roll back the last migration for Monitoring Tables`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunRollBack(rollBackArgs); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	SQLStoreCmd.AddCommand(rollBackCmd)
	rollBackArgs.SQLStoreArgs = &sqlstoreArgs
}

func RunRollBack(args RollBackArgs) error {
	cfg, logger, err := config.GetConfigAndLogger(args.ConfigFilePath, args.Debug)
	if err != nil {
		return err
	}

	if err := sqlstore.RevertOneVersion(logger, cfg.SQLStore.GetConnectionConfig()); err != nil {
		return err
	}

	return nil
}
