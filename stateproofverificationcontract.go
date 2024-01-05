package stateproofverificationcontract

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"log"
	"stateproofverificationcontract/internal/utils"
)

var Version string

func CreateApplication(algodClient *algod.Client, creator crypto.Account) (uint64, error) {
	var (
		localInts   uint64 = 0
		localBytes  uint64 = 0
		globalInts  uint64 = 0
		globalBytes uint64 = 1 // 1 for the admin account
	)

	approvalProgramBinary, err := utils.CompileTealProgram(algodClient, ".build/approval_program.teal")
	if err != nil {
		return 0, fmt.Errorf("failed to compile approval program: %+v", err)
	}
	clearStateProgramBinary, err := utils.CompileTealProgram(algodClient, ".build/clear_state_program.teal")
	if err != nil {
		return 0, fmt.Errorf("failed to compile clear state program: %+v", err)
	}

	suggestedParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		return 0, fmt.Errorf("error getting suggested transaction params: %+v", err)
	}
	txn, err := transaction.MakeApplicationCreateTx(
		false,
		approvalProgramBinary,
		clearStateProgramBinary,
		types.StateSchema{NumUint: globalInts, NumByteSlice: globalBytes},
		types.StateSchema{NumUint: localInts, NumByteSlice: localBytes},
		[][]byte{[]byte(creator.Address.String())},
		nil,
		nil,
		nil,
		suggestedParams,
		creator.Address,
		nil,
		types.Digest{},
		[32]byte{},
		types.ZeroAddress,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create app creation txn: %+v", err)
	}

	txnid, stxn, err := crypto.SignTransaction(creator.PrivateKey, txn)
	if err != nil {
		return 0, fmt.Errorf("failed to sign transaction: %+v", err)
	}

	_, err = algodClient.SendRawTransaction(stxn).Do(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to send transaction: %+v", err)
	}

	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txnid, 4, context.Background())
	if err != nil {
		return 0, fmt.Errorf("error waiting for confirmation: %+v", err)
	}
	appId := confirmedTxn.ApplicationIndex

	log.Printf("created app with id: %d", appId)

	return appId, nil
}

func GetVersion() string {
	return Version
}
