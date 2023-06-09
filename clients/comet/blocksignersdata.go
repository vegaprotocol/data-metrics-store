package comet

import (
	"fmt"
	"strconv"
	"time"
)

type BlockSignersData struct {
	Height          int64
	Time            time.Time
	ProposerAddress string
	SignerAddresses []string
}

func newBlockSignersData(response commitResponse) (blockSignersData BlockSignersData, err error) {
	blockSignersData.Height, err = strconv.ParseInt(response.Result.SignedHeader.Header.Height, 10, 64)
	if err != nil {
		err = fmt.Errorf("failed to parse Height '%s' to int, from: %+v, %w", response.Result.SignedHeader.Header.Height, response, err)
		return
	}

	blockSignersData.Time, err = time.Parse(time.RFC3339, response.Result.SignedHeader.Header.Time)
	if err != nil {
		err = fmt.Errorf("failed to parse Time '%s' to int, from: %+v, %w", response.Result.SignedHeader.Header.Time, response, err)
		return
	}
	blockSignersData.ProposerAddress = response.Result.SignedHeader.Header.ProposerAddress

	for _, signature := range response.Result.SignedHeader.Commit.Signatures {
		if len(signature.ValidatorAddress) > 0 {
			blockSignersData.SignerAddresses = append(blockSignersData.SignerAddresses, signature.ValidatorAddress)
		}
	}

	return
}
