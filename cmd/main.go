// just a simple program for retrieving the omm votes
// based on the python script of github.com/izyak

package main

import (
	"fmt"
	"os"
)

func main(){
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
}

