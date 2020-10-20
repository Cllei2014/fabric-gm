/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"

	"github.com/tw-bc-group/fabric-gm/core/chaincode/shim"
	"github.com/tw-bc-group/fabric-gm/integration/chaincode/simple"
)

func main() {
	err := shim.Start(&simple.SimpleChaincode{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exiting Simple chaincode: %s", err)
		os.Exit(2)
	}
}
