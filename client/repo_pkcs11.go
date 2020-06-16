// +build pkcs11

package client

import (
	"fmt"

	"github.com/autonomic-ai/notary"
	"github.com/autonomic-ai/notary/trustmanager"
	"github.com/autonomic-ai/notary/trustmanager/yubikey"
)

func getKeyStores(baseDir string, retriever notary.PassRetriever) ([]trustmanager.KeyStore, error) {
	fileKeyStore, err := trustmanager.NewKeyFileStore(baseDir, retriever)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key store in directory: %s", baseDir)
	}

	keyStores := []trustmanager.KeyStore{fileKeyStore}
	yubiKeyStore, _ := yubikey.NewYubiStore(fileKeyStore, retriever)
	if yubiKeyStore != nil {
		keyStores = []trustmanager.KeyStore{yubiKeyStore, fileKeyStore}
	}
	return keyStores, nil
}
