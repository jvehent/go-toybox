package main

import (
	"bytes"
	"camlistore.org/pkg/misc/gpgagent"
	"camlistore.org/pkg/misc/pinentry"
	"code.google.com/p/go.crypto/openpgp"
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage:\t%s <FINGERPRINT>\nex: \t%s 02699C276B60F27FC55C329968D63477A392B6FA\n",
			os.Args[0], os.Args[0])
		os.Exit(1)
	}
	keyid := os.Args[1]
	data := "some random string to sign"
	secringFile, err := os.Open(FindHomedir() + "/.gnupg/secring.gpg")
	if err != nil {
		panic(err)
	}
	defer secringFile.Close()
	keyring, err := openpgp.ReadKeyRing(secringFile)
	if err != nil {
		err = fmt.Errorf("Keyring access failed: '%v'", err)
		panic(err)
	}

	// find the entity in the keyring
	var signer *openpgp.Entity
	found := false
	for _, entity := range keyring {
		fingerprint := strings.ToUpper(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
		fmt.Println("reading key with fingerprint", fingerprint)
		if keyid == fingerprint {
			signer = entity
			found = true
			fmt.Println("found a match")
			break
		}
	}
	if !found {
		err = fmt.Errorf("No key found for ID '%s'", keyid)
		panic(err)
	}

	// if private key is encrypted, attempt to decrypt it with the cached passphrase
	// then try with an agent or by asking the user for a passphrase
	if signer.PrivateKey.Encrypted {
		// get private key passphrase
		signer, err = decryptEntity(signer)
		if err != nil {
			panic(err)
		}
	}

	// calculate signature
	out := bytes.NewBuffer(nil)
	message := bytes.NewBufferString(data)
	err = openpgp.ArmoredDetachSign(out, signer, message, nil)
	if err != nil {
		err = fmt.Errorf("Signature failed: '%v'", err)
		panic(err)
	}

	fmt.Printf("%s\n", out)
}

func FindHomedir() string {
	if runtime.GOOS == "darwin" {
		return os.Getenv("HOME")
	} else {
		// find keyring in default location
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		return u.HomeDir
	}
}

// decryptEntity calls gnupg-agent and pinentry to obtain a passphrase and
// decrypt the private key of a given entity (thank you, camlistore folks)
func decryptEntity(s *openpgp.Entity) (ds *openpgp.Entity, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("decryptEntity(): %v", e)
		}
	}()
	ds = s
	// TODO: syscall.Mlock a region and keep pass phrase in it.
	pubk := &ds.PrivateKey.PublicKey
	desc := fmt.Sprintf("Need to unlock GPG key %s to use it for signing.",
		pubk.KeyIdShortString())

	conn, err := gpgagent.NewConn()
	switch err {
	case gpgagent.ErrNoAgent:
		fmt.Fprintf(os.Stderr, "Note: gpg-agent not found; resorting to on-demand password entry.\n")
	case nil:
		defer conn.Close()
		req := &gpgagent.PassphraseRequest{
			CacheKey: "go:pgpencrypt:" + pubk.KeyIdShortString(),
			Prompt:   "Passphrase",
			Desc:     desc,
		}
		for tries := 0; tries < 3; tries++ {
			pass, err := conn.GetPassphrase(req)
			if err == nil {
				err = ds.PrivateKey.Decrypt([]byte(pass))
				if err == nil {
					return ds, err
				}
				req.Error = "Passphrase failed to decrypt: " + err.Error()
				conn.RemoveFromCache(req.CacheKey)
				continue
			}
			if err == gpgagent.ErrCancel {
				panic("failed to decrypt key; action canceled")
			}
		}
	default:
		panic(err)
	}

	pinReq := &pinentry.Request{Desc: desc, Prompt: "Passphrase"}
	for tries := 0; tries < 3; tries++ {
		pass, err := pinReq.GetPIN()
		if err == nil {

			err = ds.PrivateKey.Decrypt([]byte(pass))
			if err == nil {
				return ds, err
			}
			pinReq.Error = "Passphrase failed to decrypt: " + err.Error()
			continue
		}
		if err == pinentry.ErrCancel {
			panic("failed to decrypt key; action canceled")
		}
	}
	return ds, fmt.Errorf("decryptEntity(): failed to decrypt key %q: %v", pubk.KeyIdShortString(), err)
}
