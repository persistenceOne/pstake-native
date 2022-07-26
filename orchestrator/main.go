package main

import (
	orc "github.com/persistenceOne/pstake-native/orchestrator/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{Use: "pstake-orchestrator",
		Short: "pstake-orchestrator is a tool to relay pstake transactions to the Native network",
		Long:  "pstake-orchestrator is a tool to relay pstake transactions to the Native network",
	}
	//TODO: add commands

	rootCmd.AddCommand(orc.InitCommand())
	rootCmd.AddCommand(orc.StartCommand())
	err := rootCmd.Execute()
	if err != nil {
		panic(any(err))
	}
}
