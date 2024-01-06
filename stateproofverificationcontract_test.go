package stateproofverificationcontract

import (
	"fmt"
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"log"
	"strings"
	"testing"
)

var (
	algodHost         = "http://localhost"
	algodPort         = "4001"
	algodToken        = strings.Repeat("a", 64)
	kmdHost           = "http://localhost"
	kmdPort           = "4002"
	kmdToken          = strings.Repeat("a", 64)
	kmdWalletName     = "unencrypted-default-wallet"
	kmdWalletPassword = ""
)

func createAlgodClient() *algod.Client {
	algodClient, err := algod.MakeClient(
		fmt.Sprintf("%s:%s", algodHost, algodPort),
		algodToken,
	)

	if err != nil {
		log.Fatalf("failed to create algod client: %s", err)
	}

	return algodClient
}

func createKmdClient() kmd.Client {
	kmdClient, err := kmd.MakeClient(
		fmt.Sprintf("%s:%s", kmdHost, kmdPort),
		kmdToken,
	)

	if err != nil {
		log.Fatalf("failed to create kmd client: %s", err)
	}

	return kmdClient
}

func getSandboxAccounts() ([]crypto.Account, error) {
	var accounts []crypto.Account
	var walletId string

	kmdClient := createKmdClient()

	resp, err := kmdClient.ListWallets()
	if err != nil {
		return nil, fmt.Errorf("failed to list wallets: %+v", err)
	}

	for _, wallet := range resp.Wallets {
		if wallet.Name == kmdWalletName {
			walletId = wallet.ID
		}
	}

	if walletId == "" {
		return nil, fmt.Errorf("no wallet named %s", kmdWalletName)
	}

	whResp, err := kmdClient.InitWalletHandle(walletId, kmdWalletPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to init wallet handle: %+v", err)
	}

	addrResp, err := kmdClient.ListKeys(whResp.WalletHandleToken)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %+v", err)
	}

	for _, addr := range addrResp.Addresses {
		expResp, err := kmdClient.ExportKey(whResp.WalletHandleToken, kmdWalletPassword, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to export key: %+v", err)
		}

		acct, err := crypto.AccountFromPrivateKey(expResp.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create account from private key: %+v", err)
		}

		accounts = append(accounts, acct)
	}

	return accounts, nil
}

func TestCreateApplication(t *testing.T) {
	algodClient := createAlgodClient()
	accounts, err := getSandboxAccounts()
	if err != nil {
		log.Fatalf("failed to get sandbox accounts: %s", err)
	}

	creator := accounts[0]

	_, err = CreateApplication(algodClient, creator)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAdminAddress(t *testing.T) {
	algodClient := createAlgodClient()
	accounts, err := getSandboxAccounts()
	if err != nil {
		log.Fatalf("failed to get sandbox accounts: %s", err)
	}

	creator := accounts[0]

	application, err := CreateApplication(algodClient, creator)
	if err != nil {
		log.Fatal(err)
	}

	adminAddress, err := application.GetAdminAddress()
	if err != nil {
		log.Fatal(err)
	}

	if adminAddress != creator.Address {
		t.Errorf("output %s not equal to expected %s", adminAddress.String(), creator.Address.String())
	}
}
