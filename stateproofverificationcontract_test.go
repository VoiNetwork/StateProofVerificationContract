package stateproofverificationcontract

import (
	"fmt"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()

	fmt.Println("Version:", version)
}
