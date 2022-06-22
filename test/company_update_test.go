package test

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"testing"
	apiError "xm/error"
)

func TestUpdateCompany(t *testing.T) {
	testApplication.PrepareEmptyTables()
	
	company := addCompanyToDB(t, "ABC Enterprise", "001", "India", "https://www.abc.com", "990100000")

	type errData struct {
		errKey string
		err    map[string]string
	}

	type payload struct {
		Name    string `json:"name"`
		Code    string `json:"code"`
		Country string `json:"country"`
		Website string `json:"website"`
		Phone   string `json:"phone"`
	}

	testCases := []struct {
		name      string
		companyID string
		payload   payload
		wantCode  int
		want      *companyDTO
		wantErr   *errData
	}{
		{"+ve:ShouldUpdateCompany",
			company.ID.String(),
			payload{
				Name:    "ABC Enterprise 001",
				Code:    "002",
				Country: "India",
				Website: "https://www.abc1.com",
				Phone:   "990100000",
			},
			http.StatusOK,
			&companyDTO{
				ID:      company.ID.String(),
				Name:    "ABC Enterprise 001",
				Code:    "002",
				Country: "India",
				Website: "https://www.abc1.com",
				Phone:   "990100000",
			},
			nil,
		},
		{"-ve:ShouldFailWhenNonExistingIDPassed",
			uuid.NewV4().String(),
			payload{
				Name:    "ABC Enterprise 001",
				Code:    "002",
				Country: "India",
				Website: "https://www.abc1.com",
				Phone:   "990100000",
			},
			http.StatusNotFound,
			nil,
			nil,
		},
		{"-ve:ShouldFailWhenEmptyNamePassed",
			company.ID.String(),
			payload{
				Name:    "",
				Code:    "002",
				Country: "India",
				Website: "https://www.abc1.com",
				Phone:   "990100000",
			},
			http.StatusBadRequest,
			nil,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"name": apiError.ErrorCodeRequired}},
		},
		{"-ve:ShouldFailWhenEmptyCodePassed",
			company.ID.String(),
			payload{
				Name:    "ABC Enterprise 1",
				Code:    "",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			http.StatusBadRequest,
			nil,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"code": apiError.ErrorCodeRequired}},
		},
		{"-ve:ShouldFailWhenEmptyCountryNotSpecified",
			company.ID.String(),
			payload{
				Name:    "ABC Enterprise 1",
				Code:    "001",
				Country: "",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			http.StatusBadRequest,
			nil,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"country": apiError.ErrorCodeRequired}},
		},
		{"-ve:ShouldFailWhenInvalidWebsitePassed",
			company.ID.String(),
			payload{
				Name:    "ABC Enterprise 1",
				Code:    "001",
				Country: "India",
				Website: "abc.com",
				Phone:   "990100000",
			},
			http.StatusBadRequest,
			nil,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"website": apiError.ErrorCodeInvalidValue}},
		},
		{"-ve:ShouldFailWhenInvalidPhonePassed",
			company.ID.String(),
			payload{
				Name:    "ABC Enterprise 1",
				Code:    "001",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000abc",
			},
			http.StatusBadRequest,
			nil,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"phone": apiError.ErrorCodeInvalidValue}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			response := callAPI(http.MethodPut, fmt.Sprintf("/api/companies/%s", tt.companyID), tt.payload)

			checkResponseCode(t, tt.wantCode, response.Code)

			if tt.wantErr != nil {
				wantErr := *tt.wantErr
				for errField, errMsg := range wantErr.err {
					assertErrorResponse(t, response, wantErr.errKey, errField, errMsg)
				}
				return
			}

			if tt.wantCode == http.StatusNotFound {
				return
			}

			var responseDto companyDTO
			if err := json.Unmarshal(response.Body.Bytes(), &responseDto); err != nil {
				t.Errorf("unable to parse response: %v", err)
				return
			}

			if tt.want.ID != responseDto.ID {
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
