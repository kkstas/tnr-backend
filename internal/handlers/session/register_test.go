package session_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestRegister(t *testing.T) {
	t.Run("registers new user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		userFC := struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		}{
			Email:     "doe@johndoe.com",
			Password:  "mypassword123",
			FirstName: "John",
			LastName:  "Doe",
		}

		request := httptest.NewRequest("POST", "/register", testutils.ToJSONBuffer(t, userFC))
		response := httptest.NewRecorder()

		serv, cleanup, db := testutils.NewTestApplication(t)
		t.Cleanup(cleanup)
		serv.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		foundUsers, err := repositories.NewUserRepo(db).FindAll(ctx)
		testutils.AssertNoError(t, err)

		testutils.AssertEqual(t, foundUsers[0].Email, userFC.Email)
		testutils.AssertEqual(t, foundUsers[0].FirstName, userFC.FirstName)
		testutils.AssertEqual(t, foundUsers[0].LastName, userFC.LastName)
		testutils.AssertValidDate(t, foundUsers[0].CreatedAt)
	})

	t.Run("should reject invalid request properties", func(t *testing.T) {
		t.Parallel()
		serv, cleanup, _ := testutils.NewTestApplication(t)
		t.Cleanup(cleanup)

		type reqBody struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		}

		tests := []struct {
			key    string
			reason string
			value  reqBody
		}{
			{key: "email", reason: "invalid email", value: reqBody{
				Email:     "doe@johndoe",
				Password:  "mypassword123",
				FirstName: "John",
				LastName:  "Doe",
			}},
			{key: "password", reason: "too short", value: reqBody{
				Email:     "doe@johndoe.com",
				Password:  "aa",
				FirstName: "John",
				LastName:  "Doe",
			}},
			{key: "password", reason: "too long", value: reqBody{
				Email:     "doe@johndoe.com",
				Password:  testutils.RandomString(501),
				FirstName: "John",
				LastName:  "Doe",
			}},
			{key: "firstName", reason: "too short", value: reqBody{
				Email:     "doe@johndoe.com",
				Password:  "mypassword123",
				FirstName: "J",
				LastName:  "Doe",
			}},
			{key: "firstName", reason: "too long", value: reqBody{
				Email:     "doe@johndoe.com",
				Password:  "mypassword123",
				FirstName: testutils.RandomString(51),
				LastName:  "Doe",
			}},
			{key: "lastName", reason: "too short", value: reqBody{
				Email:     "doe@johndoe.com",
				Password:  "mypassword123",
				FirstName: "John",
				LastName:  "x",
			}},
			{key: "lastName", reason: "too long", value: reqBody{
				Email:     "doe@johndoe.com",
				Password:  "mypassword123",
				FirstName: "John",
				LastName:  testutils.RandomString(51),
			}},
		}

		for _, tc := range tests {
			t.Run(tc.key+tc.reason, func(t *testing.T) {
				t.Parallel()
				request := httptest.NewRequest("POST", "/register", testutils.ToJSONBuffer(t, tc.value))
				response := httptest.NewRecorder()
				serv.ServeHTTP(response, request)
				testutils.AssertStatus(t, response.Code, http.StatusBadRequest)

				var m map[string]string
				if err := json.NewDecoder(response.Body).Decode(&m); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}
				testutils.AssertNotEmpty(t, m[tc.key])
			})
		}
	})
}
