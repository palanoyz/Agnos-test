package hospital

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, "/patient/search/123", request.URL.Path)
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"first_name_en":"Alice","national_id":"123","gender":"F"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	patient, err := client.Search(context.Background(), "123")

	require.NoError(t, err)
	require.Equal(t, "Alice", patient.FirstNameEN)
	require.Equal(t, "123", patient.NationalID)
}

func TestClientSearchRejectsNonSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	_, err := NewClient(server.URL, server.Client()).Search(context.Background(), "123")
	require.Error(t, err)
}
