package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootConfig config

var rootCmd = &cobra.Command{
	Use:          "loadtester",
	Short:        "Executa testes de carga em serviços HTTP",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(command *cobra.Command, _ []string) error {
		if rootConfig.Requests <= 0 {
			return fmt.Errorf("--requests deve ser maior que zero")
		}
		if rootConfig.Concurrency <= 0 {
			return fmt.Errorf("--concurrency deve ser maior que zero")
		}

		executionReport := runLoadTest(rootConfig)
		printReport(command.OutOrStdout(), executionReport)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&rootConfig.URL, "url", "", "URL do serviço a ser testado")
	rootCmd.Flags().IntVar(&rootConfig.Requests, "requests", 0, "Número total de requisições")
	rootCmd.Flags().IntVar(&rootConfig.Concurrency, "concurrency", 0, "Número de chamadas simultâneas")

	_ = rootCmd.MarkFlagRequired("url")
	_ = rootCmd.MarkFlagRequired("requests")
	_ = rootCmd.MarkFlagRequired("concurrency")
}
