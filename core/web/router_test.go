package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenAuthRequired_NoCredentials(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/v2/specs/", web.MediaType, bytes.NewBufferString("{}"))
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestTokenAuthRequired_SessionCredentials(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := app.NewHTTPClient()
	resp, cleanup := client.Post("/v2/specs/", nil)
	defer cleanup()

	assert.Equal(t, 400, resp.StatusCode)
}

func TestTokenAuthRequired_TokenCredentials(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	eia := models.NewExternalInitiatorAuthentication()
	ea, err := models.NewExternalInitiator(eia)
	require.NoError(t, err)
	err = app.GetStore().CreateExternalInitiator(ea)
	require.NoError(t, err)

	request, err := http.NewRequest("POST", ts.URL+"/v2/specs/", bytes.NewBufferString("{}"))
	require.NoError(t, err)
	request.Header.Set("Content-Type", web.MediaType)
	request.Header.Set("X-Chainlink-EA-AccessKey", eia.AccessKey)
	request.Header.Set("X-Chainlink-EA-Secret", eia.Secret)

	client := http.Client{}
	resp, err := client.Do(request)
	require.NoError(t, err)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestTokenAuthRequired_BadTokenCredentials(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	eia := models.NewExternalInitiatorAuthentication()
	ea, err := models.NewExternalInitiator(eia)
	require.NoError(t, err)
	err = app.GetStore().CreateExternalInitiator(ea)
	require.NoError(t, err)

	request, err := http.NewRequest("POST", ts.URL+"/v2/specs/", bytes.NewBufferString("{}"))
	require.NoError(t, err)
	request.Header.Set("Content-Type", web.MediaType)
	request.Header.Set("X-Chainlink-EA-AccessKey", eia.AccessKey)
	request.Header.Set("X-Chainlink-EA-Secret", "every unpleasant commercial color from aquamarine to beige")

	client := http.Client{}
	resp, err := client.Do(request)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
}
