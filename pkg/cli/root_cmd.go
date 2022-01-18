package cli

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "cardano-tool",
		Short: "Performs Cardano commands",
	}
	commandRunner Runner
)

//Runner interface so we can mock it later
type Runner interface {
	ExecuteCommand(args ...string) ([]byte, error)
}

//NewCardanoTool returns the cobra command for the application
func NewCardanoTool(runner Runner) *cobra.Command {
	commandRunner = runner
	return rootCmd
}
