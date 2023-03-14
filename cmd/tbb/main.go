package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Checando por erros no codigo e loga o erro, usar ao inves de iferr
func Check(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "tbb",
		Short: "Deli da Persy - PersyCoins",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	err := tbbCmd.Execute()
	Check(err)
}
