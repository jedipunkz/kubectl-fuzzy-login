package main

import (
	"log"

	"github.com/jedipunkz/kubectl-login-pod/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute root command: %v", err)
	}
}
