// +build pkcs11

package main

import (
	"github.com/autonomic-ai/notary"
	store "github.com/autonomic-ai/notary/storage"
	"github.com/autonomic-ai/notary/trustmanager"
	"github.com/autonomic-ai/notary/trustmanager/yubikey"
)

func getYubiStore(fileKeyStore trustmanager.KeyStore, ret notary.PassRetriever) (*yubikey.YubiStore, error) {
	return yubikey.NewYubiStore(fileKeyStore, ret)
}

func getImporters(baseDir string, ret notary.PassRetriever) ([]trustmanager.Importer, error) {

	var importers []trustmanager.Importer
	if yubikey.IsAccessible() {
		yubiStore, err := getYubiStore(nil, ret)
		if err == nil {
			importers = append(
				importers,
				yubikey.NewImporter(yubiStore, ret),
			)
		}
	}
	fileStore, err := store.NewPrivateKeyFileStorage(baseDir, notary.KeyExtension)
	if err == nil {
		importers = append(
			importers,
			fileStore,
		)
	} else if len(importers) == 0 {
		return nil, err // couldn't initialize any stores
	}
	return importers, nil
}
