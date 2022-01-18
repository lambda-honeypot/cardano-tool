package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
)

var _ = Describe("cardano-tool gen-addr", func() {
	Describe("gen-addr --help", func() {
		It("display the gen-addr help text", func() {
			output, err := exec.Command(cardanoToolBinary, "gen-addr", "--help").CombinedOutput()
			Expect(err).ToNot(HaveOccurred(), string(output))
			Expect(string(output)).To(ContainSubstring("Usage:"))
			Expect(string(output)).To(ContainSubstring("cardano-tool gen-addr [mainnet or testnet number] [flags]"))
			Expect(string(output)).To(ContainSubstring("Flags:"))
			Expect(string(output)).To(ContainSubstring("-h, --help"))
			Expect(string(output)).To(ContainSubstring("help for gen-addr"))
			Expect(string(output)).To(ContainSubstring("Output directory to write to"))
			Expect(string(output)).To(ContainSubstring("-o, --output-dir string"))
		})
	})
})
