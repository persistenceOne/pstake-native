package main

import (
	"github.com/persistenceOne/pStake-native/oracle/command/orc"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{Use: "pstake-oracle",
		Short: "pstake-oracle is a tool to relay pstake transactions to the Native network",
		Long:  "pstake-oracle is a tool to relay pstake transactions to the Native network",
	}

	//TODO: add commands

	rootCmd.AddCommand(orc.InitCommand())

	rootCmd.AddCommand(orc.StartCommand())

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
