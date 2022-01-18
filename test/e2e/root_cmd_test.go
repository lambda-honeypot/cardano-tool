package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"testing"
)

const cardanoToolBinary = "../../build/cardano-tool"

func checkDefaultMessage(output []byte, err error) {
	Expect(err).ToNot(HaveOccurred(), string(output))
	Expect(string(output)).To(ContainSubstring("Performs Cardano commands"))
	Expect(string(output)).To(ContainSubstring("Usage:\n  cardano-tool [command]"))
	Expect(string(output)).To(ContainSubstring("Available Commands:"))
	Expect(string(output)).To(ContainSubstring("Generate a new payment address and output the files"))
	Expect(string(output)).To(ContainSubstring("help        Help about any command"))
	Expect(string(output)).To(ContainSubstring("version     Print the version number of Cardano"))
	Expect(string(output)).To(ContainSubstring("Flags:"))
	Expect(string(output)).To(ContainSubstring("-h, --help   help for cardano-tool"))
	Expect(string(output)).To(ContainSubstring("Use \"cardano-tool [command] --help\" for more information about a command."))
}

func TestCardanoTool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test CLI")
}

var _ = Describe("cardano-tool command line", func() {
	Describe("no flags", func() {
		It("should print default message", func() {
			output, err := exec.Command(cardanoToolBinary).CombinedOutput()
			checkDefaultMessage(output, err)
		})
	})

	Describe("--help flag", func() {
		It("--help message", func() {
			output, err := exec.Command(cardanoToolBinary, "--help").CombinedOutput()
			checkDefaultMessage(output, err)
		})
	})

	Describe("help command", func() {
		It("help command message", func() {
			output, err := exec.Command(cardanoToolBinary, "help").CombinedOutput()
			checkDefaultMessage(output, err)
		})
	})

	Describe("--unknown flag", func() {
		It("display unknown message", func() {
			output, err := exec.Command(cardanoToolBinary, "--unknown").CombinedOutput()
			Expect(err).To(HaveOccurred(), string(output))
			Expect(string(output)).To(ContainSubstring("Error: unknown flag: --unknown"))
		})
	})

	Describe("version command", func() {
		It("displays the version in correct format", func() {
			output, err := exec.Command(cardanoToolBinary, "version").CombinedOutput()
			Expect(err).ToNot(HaveOccurred(), string(output))
			Expect(string(output)).To(ContainSubstring("cardano-tool:v0.1.0"))
		})
	})

	Describe("version --help", func() {
		It("display the version help text", func() {
			output, err := exec.Command(cardanoToolBinary, "version", "--help").CombinedOutput()
			Expect(err).ToNot(HaveOccurred(), string(output))
			Expect(string(output)).To(ContainSubstring("Print the version number of Cardano"))
			Expect(string(output)).To(ContainSubstring("Usage:"))
			Expect(string(output)).To(ContainSubstring("cardano-tool version [flags]"))
			Expect(string(output)).To(ContainSubstring("Flags:"))
			Expect(string(output)).To(ContainSubstring("-h, --help   help for version"))
		})
	})
})
