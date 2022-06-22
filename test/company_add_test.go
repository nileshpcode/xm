package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	apiError "xm/error"
)

func TestAddCompany(t *testing.T) {
	testApplication.PrepareEmptyTables()

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
		name                    string
		setInvalidRequestOrigin bool
		payload                 payload
		want                    *companyDTO
		wantHttpStatus          int
		wantErr                 *errData
	}{
		{"+ve:ShouldAddCompany",
			false,
			payload{
				Name:    "ABC Enterprise",
				Code:    "001",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			&companyDTO{
				Name:    "ABC Enterprise",
				Code:    "001",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			http.StatusCreated,
			nil,
		},
		{"-ve:ShouldFailWhenInvalidRequestOrigin",
			true,
			payload{
				Name:    "ABC Enterprise",
				Code:    "001",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			&companyDTO{
				Name:    "ABC Enterprise",
				Code:    "001",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			http.StatusUnauthorized,
			nil,
		},
		{"-ve:ShouldFailWhenEmptyNamePassed",
			false,
			payload{
				Name:    "",
				Code:    "001",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			nil,
			http.StatusBadRequest,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"name": apiError.ErrorCodeRequired}},
		},
		{"-ve:ShouldFailWhenEmptyCodePassed",
			false,
			payload{
				Name:    "ABC Enterprise 1",
				Code:    "",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			nil,
			http.StatusBadRequest,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"code": apiError.ErrorCodeRequired}},
		},
		{"-ve:ShouldFailWhenEmptyCountryNotSpecified",
			false,
			payload{
				Name:    "ABC Enterprise 1",
				Code:    "001",
				Country: "",
				Website: "https://www.abc.com",
				Phone:   "990100000",
			},
			nil,
			http.StatusBadRequest,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"country": apiError.ErrorCodeRequired}},
		},
		{"-ve:ShouldFailWhenInvalidWebsitePassed",
			false,
			payload{
				Name:    "ABC Enterprise 1",
				Code:    "001",
				Country: "India",
				Website: "abc.com",
				Phone:   "990100000",
			},
			nil,
			http.StatusBadRequest,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"website": apiError.ErrorCodeInvalidValue}},
		},
		{"-ve:ShouldFailWhenInvalidPhonePassed",
			false,
			payload{
				Name:    "ABC Enterprise 1",
				Code:    "001",
				Country: "India",
				Website: "https://www.abc.com",
				Phone:   "990100000abc",
			},
			nil,
			http.StatusBadRequest,
			&errData{apiError.ErrorCodeInvalidFields, map[string]string{"phone": apiError.ErrorCodeInvalidValue}},
		},
	}

	for _, tt := range testCases {
		os.Unsetenv("ORIGIN_COUNTRY")

		t.Run(tt.name, func(t *testing.T) {

			if tt.setInvalidRequestOrigin {
				os.Setenv("ORIGIN_COUNTRY", "US")
			}

			response := callAPI(http.MethodPost, "/api/companies", tt.payload)

			checkResponseCode(t, tt.wantHttpStatus, response.Code)

			if tt.wantErr == nil && tt.wantHttpStatus != http.StatusCreated {
				return
			}

			if tt.wantErr != nil {
				wantErr := *tt.wantErr
				for errField, errMsg := range wantErr.err {
					fmt.Println(errField, errMsg, "dsf")
					assertErrorResponse(t, response, wantErr.errKey, errField, errMsg)
				}
				return
			}

			var responseDto companyDTO
			if err := json.Unmarshal(response.Body.Bytes(), &responseDto); err != nil {
				t.Errorf("unable to parse response: %v", err)
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
