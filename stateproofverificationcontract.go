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
	stateproofverificationcontracttypes "stateproofverificationcontract/types"
)

func CreateApplication(algodClient *algod.Client, signer crypto.Account) (*stateproofverificationcontracttypes.Application, error) {
	var (
		localInts   uint64 = 0
		localBytes  uint64 = 0
		globalInts  uint64 = 0
		globalBytes uint64 = 1 // 1 for the admin account
	)

	approvalProgramBinary, err := utils.CompileTealProgram(algodClient, ".build/approval_program.teal")
	if err != nil {
		return nil, fmt.Errorf("failed to compile approval program: %+v", err)
	}
	clearStateProgramBinary, err := utils.CompileTealProgram(algodClient, ".build/clear_state_program.teal")
	if err != nil {
		return nil, fmt.Errorf("failed to compile clear state program: %+v", err)
	}

	suggestedParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting suggested transaction params: %+v", err)
	}
	txn, err := transaction.MakeApplicationCreateTx(
		false,
		approvalProgramBinary,
		clearStateProgramBinary,
		types.StateSchema{NumUint: globalInts, NumByteSlice: globalBytes},
		types.StateSchema{NumUint: localInts, NumByteSlice: localBytes},
		[][]byte{[]byte(signer.Address.String())},
		nil,
		nil,
		nil,
		suggestedParams,
		signer.Address,
		nil,
		types.Digest{},
		[32]byte{},
		types.ZeroAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create app creation txn: %+v", err)
	}

	// create the application on-chain
	_, confirmedTxn, err := utils.SignAndSendTransaction(txn, signer, algodClient)
	if err != nil {
		return nil, err
	}
	appId := confirmedTxn.ApplicationIndex

	log.Printf("created app with id: %d", appId)

	application := &stateproofverificationcontracttypes.Application{
		AlgodClient:       algodClient,
		AppId:             appId,
		ApprovalProgram:   approvalProgramBinary,
		ClearStateProgram: clearStateProgramBinary,
		Signer:            signer,
	}

	// fund the minimum app account balance
	err = application.FundAppAccount(100000, signer, suggestedParams, []byte(fmt.Sprintf("payment for minimum balance for application '%d' account", appId)))

	return application, nil
}

func InitializeApplication(algodClient *algod.Client, appId uint64, signer crypto.Account) (*stateproofverificationcontracttypes.Application, error) {
	approvalProgramBinary, err := utils.CompileTealProgram(algodClient, ".build/approval_program.teal")
	if err != nil {
		return nil, fmt.Errorf("failed to compile approval program: %+v", err)
	}
	clearStateProgramBinary, err := utils.CompileTealProgram(algodClient, ".build/clear_state_program.teal")
	if err != nil {
		return nil, fmt.Errorf("failed to compile clear state program: %+v", err)
	}

	return &stateproofverificationcontracttypes.Application{
		AlgodClient:       algodClient,
		AppId:             appId,
		ApprovalProgram:   approvalProgramBinary,
		ClearStateProgram: clearStateProgramBinary,
		Signer:            signer,
	}, nil
}
