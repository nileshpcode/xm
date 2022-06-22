package test

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"testing"
)

func TestDeleteCompany(t *testing.T) {
	testApplication.PrepareEmptyTables()

	company := addCompanyToDB(t, "ABC Enterprise", "001", "India", "https://www.abc.com", "990100000")

	tests := []struct {
		name                    string
		setInvalidRequestOrigin bool
		companyID               string
		wantHttpStatusCode      int
	}{
		{"+ve:ShouldDeleteCompany",
			false,
			company.ID.String(),
			http.StatusOK,
		},
		{"-ve:ShouldFailWhenCompanyDoesntExist",
			false,
			uuid.NewV4().String(),
			http.StatusNotFound,
		},
		{"+ve:ShouldDeleteCompany",
			true,
			company.ID.String(),
			http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		os.Unsetenv("ORIGIN_COUNTRY")

		t.Run(tt.name, func(t *testing.T) {

			if tt.setInvalidRequestOrigin {
				os.Setenv("ORIGIN_COUNTRY", "US")
			}

			response := callAPI(http.MethodDelete, fmt.Sprintf("/api/companies/%s", tt.companyID), nil)

			checkResponseCode(t, tt.wantHttpStatusCode, response.Code)

			if tt.wantHttpStatusCode == http.StatusOK {
				exists, _ := getCompanyToDB(t, tt.companyID)
				if exists {
					t.Errorf("compay record not removed from db")
				}
			}
		})
	}
}
