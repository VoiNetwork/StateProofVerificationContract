package utils

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"log"
	"os"
)

func CalculateAppMbrForBox(keyBytesSize int, valueBytesSize int) uint64 {
	return uint64(2500 + 400*(keyBytesSize*valueBytesSize))
}

func CompileTealProgram(algodClient *algod.Client, path string) ([]byte, error) {
	teal, err := os.ReadFile(path)
	if err != nil {
		log.Println("failed to read teal program")

		return nil, err
	}

	result, err := algodClient.TealCompile(teal).Do(context.Background())
	if err != nil {
		log.Println("failed to compile program")

		return nil, err
	}

	bin, err := base64.StdEncoding.DecodeString(result.Result)
	if err != nil {
		log.Println("failed to decode compiled program")

		return nil, err
	}

	return bin, nil
}

func ConvertUint64ToByteArray(input uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)

	binary.PutUvarint(buf, input)

	return buf
}

func SignAndSendTransaction(txn types.Transaction, signer crypto.Account, client *algod.Client) (string, *models.PendingTransactionInfoResponse, error) {
	txnid, stxn, err := crypto.SignTransaction(signer.PrivateKey, txn)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign transaction: %+v", err)
	}

	_, err = client.SendRawTransaction(stxn).Do(context.Background())
	if err != nil {
		return "", nil, fmt.Errorf("failed to send transaction: %+v", err)
	}

	confirmedTxn, err := transaction.WaitForConfirmation(client, txnid, 4, context.Background())
	if err != nil {
		return "", nil, fmt.Errorf("error waiting for confirmation: %+v", err)
	}

	return txnid, &confirmedTxn, nil
}
