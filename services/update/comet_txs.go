package update

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vega/logging"
	"github.com/vegaprotocol/vega-monitoring/services/read"
	"github.com/vegaprotocol/vega-monitoring/sqlstore"
	"go.uber.org/zap"
)

func (us *UpdateService) UpdateCometTxsAllNew(ctx context.Context) error {
	return us.UpdateCometTxs(ctx, 0, 0)
}

func (us *UpdateService) UpdateCometTxs(ctx context.Context, fromBlock int64, toBlock int64) error {
	var err error
	serviceStore := us.storeService.NewCometTxs()
	logger := us.log.With(zap.String(UpdaterType, "UpdateCometTxs"))

	// get Last Block
	if toBlock <= 0 {
		logger.Debug("getting network toBlock network height")
		toBlock, err = us.readService.GetNetworkLatestBlockHeight()
		if err != nil {
			return fmt.Errorf("failed to Update Comet Txs, %w", err)
		}
	}

	// get First Block
	if fromBlock <= 0 {
		logger.Debug("getting network fromBlock network height")
		lastProcessedBlock, err := serviceStore.GetLastestBlockInStore(context.Background())
		if err != nil {
			return fmt.Errorf("failed to Update Comet Txs, %w", err)
		}
		if lastProcessedBlock > 0 {
			fromBlock = lastProcessedBlock + 1
		} else {
			fromBlock = toBlock - (BLOCK_NUM_IN_24h * 3)
			if fromBlock <= 0 {
				fromBlock = 1
			}
		}
	}
	if fromBlock > toBlock {
		return fmt.Errorf("cannot update Comet Txs, from block '%d' is greater than to block '%d'", fromBlock, toBlock)
	}

	// Update in batches
	logger.Info(
		"Update Comet Txs in batches",
		zap.Int64("first block", fromBlock),
		zap.Int64("last block", toBlock),
	)

	var totalCount int = 0

	logger.Debugf("getting batch blocks from %d to %d", fromBlock, toBlock)
	for batchFirstBlock := fromBlock; batchFirstBlock <= toBlock; batchFirstBlock += 200 {
		batchLastBlock := batchFirstBlock + 199 // endBlock inclusive
		if batchLastBlock > toBlock {
			batchLastBlock = toBlock
		}
		count, err := UpdateCometTxsRange(batchFirstBlock, batchLastBlock, us.readService, serviceStore, logger)
		if err != nil {
			return err
		}
		totalCount += count
	}

	logger.Info(
		"Finished",
		zap.Int64("processed blocks", toBlock-fromBlock+1),
		zap.Int("total row count sotred in SQLStore", totalCount),
	)

	return nil
}

func UpdateCometTxsRange(
	fromBlock int64,
	toBlock int64,
	readService *read.ReadService,
	serviceStore *sqlstore.CometTxs,
	logger *logging.Logger,
) (int, error) {

	txs, err := readService.GetCometTxs(fromBlock, toBlock)
	if err != nil {
		return -1, err
	}
	logger.Info(
		"fetched data from CometBFT",
		zap.Int64("from-block", fromBlock),
		zap.Int64("to-block", toBlock),
		zap.Int("tx count", len(txs)),
	)

	for _, tx := range txs {
		serviceStore.AddWithoutTime(tx)
	}

	storedData, err := serviceStore.FlushUpsertWithoutTime(context.Background())
	storedCount := len(storedData)
	if err != nil {
		return storedCount, err
	}
	logger.Info(
		"stored data in SQLStore",
		zap.Int64("from-block", fromBlock),
		zap.Int64("to-block", toBlock),
		zap.Int("row count", storedCount),
	)

	return storedCount, nil
}
