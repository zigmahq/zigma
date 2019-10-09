package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/zigmahq/zigma/version"
)

const tmpl = `
package version

// Current returns the build code of release
var Current = Version{
	Number: "{{ .number }}",
	Name: "{{ .name }}",
	Signature: decode("{{ .sig }}"),
}

func init() {
	Verifier = decode("{{ .public_key }}")
}
`

// DefaultSignerDir defines the default directory for
// the release signing public and private keys
var DefaultSignerDir = os.ExpandEnv("$HOME/.zsigner")

var cmd = &cobra.Command{
	Use:   "zsigner",
	Short: "Zsigner helps signing and verifying zigma binary release",
}

var sign = &cobra.Command{
	Use:   "sign",
	Short: "Sign version and generate a signature for binary release",
	Run: func(cmd *cobra.Command, args []string) {
		pub, pri := getkey()

		dir, err := os.Getwd()
		if err != nil {
			warn(err)
		}

		vff := path.Join(dir, "version.go")

		ver := &version.Version{
			Number: "0.0.1",
			Name:   "autumn-waterfall",
		}

		sig, err := ver.Sign(pri)
		if err != nil {
			warn(err)
		}

		tmpgen(vff, tmpl, map[string]interface{}{
			"number":     ver.Number,
			"name":       ver.Name,
			"sig":        hex.EncodeToString(sig),
			"public_key": hex.EncodeToString(pub),
		})
	},
}

var verify = &cobra.Command{
	Use:   "verify",
	Short: "Verify a binary release",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func warn(err error) {
	log.Println(err)
	os.Exit(1)
}

func tmpgen(path, tmpl string, params map[string]interface{}) {
	f, err := os.Create(path)
	if err != nil {
		warn(err)
	}
	t, err := template.New("zsigner").Parse(strings.TrimSpace(tmpl))
	if err != nil {
		warn(err)
	}
	if err := t.Execute(f, params); err != nil {
		warn(err)
	}
}

func fexist(paths ...string) bool {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getkey() (ed25519.PublicKey, ed25519.PrivateKey) {
	ppub := path.Join(DefaultSignerDir, "zsigner.pub")
	ppri := path.Join(DefaultSignerDir, "zsigner.pri")

	if ok := fexist(ppub, ppri); !ok {
		return keygen()
	}
	pub, err := ioutil.ReadFile(ppub)
	if err != nil {
		warn(err)
	}
	pri, err := ioutil.ReadFile(ppri)
	if err != nil {
		warn(err)
	}
	return pub, pri
}

func keygen() (ed25519.PublicKey, ed25519.PrivateKey) {
	ppub := path.Join(DefaultSignerDir, "zsigner.pub")
	ppri := path.Join(DefaultSignerDir, "zsigner.pri")

	pub, pri, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		warn(err)
	}

	os.Mkdir(DefaultSignerDir, os.ModePerm)
	ioutil.WriteFile(ppub, pub, 0644)
	ioutil.WriteFile(ppri, pri, 0644)

	return pub, pri
}

func init() {
	cmd.AddCommand(sign, verify)
}

func main() {
	if err := cmd.Execute(); err != nil {
		warn(err)
	}
}
