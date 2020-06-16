package handlers

import (
	"testing"

	"github.com/autonomic-ai/notary/tuf/data"
	"github.com/autonomic-ai/notary/tuf/signed"
	"github.com/autonomic-ai/notary/tuf/testutils"
	"github.com/stretchr/testify/require"
)

func mustCopyKeys(t *testing.T, from signed.CryptoService, roles ...data.RoleName) signed.CryptoService {
	cs, err := testutils.CopyKeys(from, roles...)
	require.NoError(t, err)
	return cs
}
