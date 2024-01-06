package types

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

type Application struct {
	AlgodClient       *algod.Client
	AppId             uint64
	ApprovalProgram   []byte
	ClearStateProgram []byte
}

func (application Application) AddStateProof() error {
	return nil
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

func (application Application) VerifyTransaction() bool {
	return false
}
