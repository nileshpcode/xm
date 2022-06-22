package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"xm/app"
	"xm/client"
	"xm/controller"
	apiError "xm/error"
	"xm/model"
	"xm/repository"
)

var testApplication *app.TestApp

func TestMain(m *testing.M) {
	routeProvider := func(app2 *app.App) []app.RouteSpecifier {
		ipLocationClient := client.NewIpLocationClient("https://ipapi.co")
		companyRepository := repository.NewRepository()

		return []app.RouteSpecifier{
			controller.NewCompanyController(app2, ipLocationClient, companyRepository),
		}
	}
	testApplication = app.NewTestApp("XM", routeProvider, initializeDB)
	testApplication.Initialize()
	code := m.Run()
	testApplication.Stop()
	os.Exit(code)
}

func initializeDB(db *gorm.DB) {
	db.Migrator().DropTable(&model.Company{})
	db.Migrator().AutoMigrate(&model.Company{})
}

// callAPI invokes http API
func callAPI(httpMethod string, apiURL string, req interface{}) *httptest.ResponseRecorder {

	var payload io.Reader
	if req != nil {
		reqJSON, _ := json.Marshal(req)
		payload = bytes.NewBuffer(reqJSON)
	}

	httpReq, _ := http.NewRequest(httpMethod, apiURL, payload)

	rr := httptest.NewRecorder()
	testApplication.Application.Router.ServeHTTP(rr, httpReq)
	return rr
}

// checkResponseCode checks if the http response is as expected
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// assertErrorResponse checks if the http response contains expected errorKey, errorField and errorMessage
func assertErrorResponse(t *testing.T, response *httptest.ResponseRecorder, expectedErrorKey string, expectedErrorField string, expectedError string) {
	var errData map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &errData); err != nil {
		t.Errorf("Unable to parse response: %v", err)
	}
	if errData["errorKey"] != expectedErrorKey {
		t.Errorf("Expected errorKey [%v], Got [%v]!", expectedErrorKey, errData["errorKey"])
	}
	errors := errData["errors"].(map[string]interface{})
	if fmt.Sprintf("%v", errors[expectedErrorField]) != expectedError {
		t.Errorf("Expected error [%v], Got [%v]!", expectedError, errors[expectedErrorField])
	}
}

// data transfer object of company
type companyDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}

func addCompanyToDB(t *testing.T, name, code, country, website, phone string) *model.Company {
	company, err := model.NewCompany(name, code, country, website, phone)
	if err != nil {
		t.Errorf("unable to add company [%v]!", err)
	}

	err = testApplication.Application.DB.Create(company).Error
	if err != nil {
		t.Errorf("unable to insert company into DB [%v]!", err)
	}

	return company
}

func getCompanyToDB(t *testing.T, id string) (bool, *model.Company) {
	company := &model.Company{}
	err := testApplication.Application.DB.First(company, "id = ?", id).Error
	if err != nil {
		if apiError.NewDatabaseError(err).IsRecordNotFoundError() {
			return false, nil
		}
		t.Errorf("unable to get company from DB [%v]!", err)
	}

	return true, company
}
