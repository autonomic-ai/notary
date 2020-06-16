// This makes sure that the server is compatible with the TUF httpstore.

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/autonomic-ai/notary"
	"github.com/autonomic-ai/notary/server/storage"
	store "github.com/autonomic-ai/notary/storage"
	"github.com/autonomic-ai/notary/tuf/data"
	"github.com/autonomic-ai/notary/tuf/signed"
	"github.com/autonomic-ai/notary/tuf/testutils"
	"github.com/autonomic-ai/notary/tuf/validation"
)

// Ensures that the httpstore can interpret the errors returned from the server
func TestValidationErrorFormat(t *testing.T) {
	ctx := context.WithValue(
		context.Background(), notary.CtxKeyMetaStore, storage.NewMemStorage())
	ctx = context.WithValue(ctx, notary.CtxKeyKeyAlgo, data.ED25519Key)

	handler := RootHandler(ctx, nil, signed.NewEd25519(), nil, nil, nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := store.NewHTTPStore(
		fmt.Sprintf("%s/v2/docker.com/notary/_trust/tuf/", server.URL),
		"",
		"json",
		"key",
		http.DefaultTransport,
	)
	require.NoError(t, err)

	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	rs, rt, _, _, err := testutils.Serialize(r, tg, sn, ts)
	require.NoError(t, err)

	// No snapshot is passed, and the server doesn't have the snapshot key,
	// so ErrBadHierarchy
	err = client.SetMulti(map[string][]byte{
		data.CanonicalRootRole.String():    rs,
		data.CanonicalTargetsRole.String(): rt,
	})
	require.Error(t, err)
	require.IsType(t, validation.ErrBadHierarchy{}, err)
}
