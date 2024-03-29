// just a simple program for retrieving the omm votes
package main

import (
	"fmt"
	"github.com/eyeonicon/go-icon-sdk/networks"
	"github.com/icon-project/goloop/client"
	"github.com/paulrouge/get-omm-votes-go/internal/calls"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("run program with address of validator as only argument.")
		return
	}

	validator := args[1]

	if len(validator) != 42 {
		fmt.Println("invalid address")
		return
	}
	c := client.NewClientV3(networks.Mainnet().URL)

	calls.ExportOMMVoters(c, validator)
}
