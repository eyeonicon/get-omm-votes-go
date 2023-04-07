// just a simple program for retrieving the omm votes
package main

import (
	"fmt"
	"os"
	"github.com/paulrouge/get-omm-votes-go/internal/calls"
	"github.com/icon-project/goloop/client"
	"github.com/eyeonicon/go-icon-sdk/networks"
)

func main(){
	fmt.Println("running")
	
	args := os.Args
	// if len(args) != 2 {
	// 	fmt.Println("run program with address of validator as only argument.")
	// 	return
	// }
	
	validator := args[1]
	_ = validator
	// if len(validator) != 42 {
	// 	fmt.Println("invalid address")
	// 	return 
	// }

	c := client.NewClientV3(networks.Mainnet().URL)
	
	calls.GetOMMTotalVotes(c ,"hx2e7db537ca3ff73336bee1bab4cf733a94ae769b")



	
}

