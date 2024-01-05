package utils

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
)

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
