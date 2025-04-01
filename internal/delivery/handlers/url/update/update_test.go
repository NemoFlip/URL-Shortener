package update_test

import (
	"RESTProject/internal/delivery/handlers/url/update"
	"RESTProject/internal/delivery/handlers/url/update/mocks"
	"RESTProject/internal/lib/logger/handlers/slogdiscard"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		newURL    string
		respError string
		mockError error
	}{
		{
			name:   "Success",
			alias:  "test_alias",
			newURL: "https://example.com",
		},
		{
			name:      "Empty alias",
			alias:     "",
			newURL:    "https://example.com",
			respError: "field Alias is a required field",
		},
		{
			name:      "Empty URL",
			newURL:    "",
			alias:     "some_alias",
			respError: "field NewURL is a required field",
		},
		{
			name:      "Invalid URL",
			newURL:    "invalid URL",
			alias:     "some_alias",
			respError: "field NewURL is not a valid URL",
		},
		{
			name:      "UpdateURL Error",
			alias:     "test_alias",
			newURL:    "https://example.com",
			respError: "failed to update url",
			mockError: errors.New("unexpected error"),
		},
		{
			name:      "Alias not found",
			alias:     "non_existing_alias",
			newURL:    "https://example.com",
			respError: "url not found",
			mockError: errors.New("url not found"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlUpdaterMock := mocks.NewURLUpdater(t)

			if tc.respError == "" || tc.mockError != nil {
				urlUpdaterMock.On("UpdateURL", tc.alias, tc.newURL).
					Return(tc.mockError).
					Once()
			}

			handler := update.New(slogdiscard.NewDiscardLogger(), urlUpdaterMock)

			input := fmt.Sprintf(`{"alias": "%s", "new_url": "%s"}`, tc.alias, tc.newURL)

			req, err := http.NewRequest(http.MethodPut, "/update", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp update.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			// TODO: add more checks
		})
	}
}
