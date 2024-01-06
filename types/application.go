package types

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"log"
	"math"
	"stateproofverificationcontract/internal/utils"
)

const (
	StateProofBlockIntervalSize uint64 = 256
)

type Application struct {
	AlgodClient       *algod.Client
	AppId             uint64
	ApprovalProgram   []byte
	ClearStateProgram []byte
	Signer            crypto.Account
}

func calculateBlockCommitmentInterval(round uint64) uint64 {
	return uint64(math.Ceil(float64(round) / float64(StateProofBlockIntervalSize)))
}

func (application Application) AddBlockHeaderCommitment(lastAttestationRound uint64, blockHeaderCommitment string) error {
	appAddress := crypto.GetApplicationAddress(application.AppId)
	if appAddress.IsZero() {
		return fmt.Errorf("found to get app address for '%d'", application.AppId)
	}
	blockCommitmentInterval := calculateBlockCommitmentInterval(lastAttestationRound)
	log.Printf("AddBlockHeaderCommitment#blockCommitmentInterval: '%d'", blockCommitmentInterval)
	boxKey := utils.ConvertUint64ToByteArray(blockCommitmentInterval)
	appArgs := [][]byte{
		[]byte("add_block_header_commitment"),
		boxKey,
		[]byte(blockHeaderCommitment),
	}
	boxStorageFee := utils.CalculateAppMbrForBox(len(appArgs[1]), len(appArgs[2]))
	suggestedParams, err := application.AlgodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		return fmt.Errorf("error getting suggested transaction params: %+v", err)
	}

	// fund the account with the box fee
	err = application.FundAppAccount(boxStorageFee, application.Signer, suggestedParams, []byte(fmt.Sprintf("payment for the box fee for application '%d'", application.AppId)))
	if err != nil {
		return err
	}

	txn, err := transaction.MakeApplicationNoOpTxWithBoxes(
		application.AppId,
		appArgs,
		nil,
		nil,
		nil,
		[]types.AppBoxReference{
			{AppID: application.AppId, Name: boxKey},
		},
		suggestedParams,
		application.Signer.Address,
		[]byte(fmt.Sprintf("send block header commitment from last attested round '%d'", lastAttestationRound)),
		types.Digest{},
		[32]byte{},
		types.ZeroAddress,
	)

	txnId, _, err := utils.SignAndSendTransaction(txn, application.Signer, application.AlgodClient)
	if err != nil {
		return err
	}

	log.Printf("added block header commitment in transaction '%s'", txnId)

	return nil
}

func (application Application) FundAppAccount(amount uint64, signer crypto.Account, suggestedParams types.SuggestedParams, note []byte) error {
	appAddress := crypto.GetApplicationAddress(application.AppId)
	if appAddress.IsZero() {
		return fmt.Errorf("found to get app address for '%d'", application.AppId)
	}
	txn, err := transaction.MakePaymentTxn(
		application.Signer.Address.String(),
		appAddress.String(),
		amount,
		note,
		"",
		suggestedParams,
	)
	if err != nil {
		return fmt.Errorf("failed to create app creation txn: %+v", err)
	}

	_, _, err = utils.SignAndSendTransaction(txn, signer, application.AlgodClient)

	return err
}

func (application Application) GetAdminAddress() (types.Address, error) {
	app, err := application.AlgodClient.GetApplicationByID(application.AppId).Do(context.Background())
	if err != nil {
		return types.ZeroAddress, fmt.Errorf("failed to get app '%d': %+v", application.AppId, err)
	}

	globalState := app.Params.GlobalState

	for _, tealValue := range globalState {
		keyBytes, err := base64.StdEncoding.DecodeString(tealValue.Key)
		if err != nil {
			return types.ZeroAddress, fmt.Errorf("failed to decode global state for app '%d': %+v", application.AppId, err)
		}

		if string(keyBytes) == "admin" {
			valueBytes, err := base64.StdEncoding.DecodeString(tealValue.Value.Bytes)
			if err != nil {
				return types.ZeroAddress, fmt.Errorf("found admin address, failed to decode the value: %+v", err)
			}

			address, err := types.DecodeAddress(string(valueBytes))
			if err != nil {
				return types.ZeroAddress, fmt.Errorf("failed to decode address '%s': %+v", string(valueBytes), err)
			}

			return address, nil
		}
	}

	return types.ZeroAddress, fmt.Errorf("admin address doesn't exist in state: %+v", err)
}

func (application Application) GetBlockHeaderCommitmentByRound(round uint64) (string, error) {
	blockCommitmentInterval := calculateBlockCommitmentInterval(round)
	log.Printf("GetBlockHeaderCommitmentByRound#blockCommitmentInterval: '%d'", blockCommitmentInterval)
	box, err := application.AlgodClient.GetApplicationBoxByName(application.AppId, utils.ConvertUint64ToByteArray(blockCommitmentInterval)).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get box for round '%d': %+v", round, err)
	}

	return string(box.Value), nil
}

func (application Application) VerifyTransaction() bool {
	return false
}
