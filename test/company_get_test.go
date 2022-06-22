package test

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"testing"
	"xm/model"
)

func TestGetCompany(t *testing.T) {
	testApplication.PrepareEmptyTables()

	company := addCompanyToDB(t, "ABC Enterprise", "001", "India", "https://www.abc.com", "990100000")

	tests := []struct {
		name      string
		companyID string
		want      *model.Company
		wantErr   bool
	}{
		{"+ve:ShouldGetCompany",
			company.ID.String(),
			company,
			false,
		},
		{"-ve:ShouldFailWhenCompanyDoesntExist",
			uuid.NewV4().String(),
			company,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := callAPI(http.MethodGet, fmt.Sprintf("/api/companies/%s", tt.companyID), nil)

			if tt.wantErr {
				checkResponseCode(t, http.StatusNotFound, response.Code)
				return
			}

			checkResponseCode(t, http.StatusOK, response.Code)

			var responseDto companyDTO
			if err := json.Unmarshal(response.Body.Bytes(), &responseDto); err != nil {
				t.Errorf("unable to parse response: %v", err)
				return
			}

			if tt.want.ID.String() != responseDto.ID {
				t.Errorf("expected id %v\nGot %v", tt.want.ID, responseDto.ID)
				return
			}

			if tt.want.Name != responseDto.Name {
				t.Errorf("expected name %v\nGot %v", tt.want.Name, responseDto.Name)
				return
			}

			if tt.want.Code != responseDto.Code {
				t.Errorf("expected code %v\nGot %v", tt.want.Code, responseDto.Code)
				return
			}

			if tt.want.Country != responseDto.Country {
				t.Errorf("expected country %v\nGot %v", tt.want.Country, responseDto.Country)
				return
			}

			if tt.want.Website != responseDto.Website {
				t.Errorf("expected website %v\nGot %v", tt.want.Website, responseDto.Website)
				return
			}

			if tt.want.Phone != responseDto.Phone {
				t.Errorf("expected phone %v\nGot %v", tt.want.Phone, responseDto.Phone)
				return
			}
		})
	}
}
