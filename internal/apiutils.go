package apiutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Config struct {
	AccessToken string
	WorkbookItemID string
	NoOfIterations int
	InputParams []InputParam
}

type InputParam struct {
	MemCnt int
	RecCnt int
	Curr string
}

type ResponsePayload struct {
	OdataContext  string          `json:"@odata.context"`
	OdataType     string          `json:"@odata.type"`
	OdataID       string          `json:"@odata.id"`
	Address       string          `json:"address"`
	AddressLocal  string          `json:"addressLocal"`
	ColumnCount   int             `json:"columnCount"`
	CellCount     int             `json:"cellCount"`
	ColumnHidden  bool            `json:"columnHidden"`
	RowHidden     bool            `json:"rowHidden"`
	NumberFormat  [][]string      `json:"numberFormat"`
	ColumnIndex   int             `json:"columnIndex"`
	Text          [][]string      `json:"text"`
	Formulas      [][]string      `json:"formulas"`
	FormulasLocal [][]string      `json:"formulasLocal"`
	FormulasR1C1  [][]string      `json:"formulasR1C1"`
	Hidden        bool            `json:"hidden"`
	RowCount      int             `json:"rowCount"`
	RowIndex      int             `json:"rowIndex"`
	ValueTypes    [][]string      `json:"valueTypes"`
	Values        [][]interface{} `json:"values"`
}

type ErrorResponse struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Date      string `json:"date"`
			RequestID string `json:"request-id"`
		} `json:"innerError"`
	} `json:"error"`
}

func RunTests(config Config) {

	// first test if the AccessToken is valid
	_,err := createSession(config.AccessToken, config.WorkbookItemID)

	if err != nil {
		fmt.Printf("Unable to create a session. Access Token expired??\n")
		return		
	}

	var wg sync.WaitGroup

	for i := 0; i < config.NoOfIterations+1; i++ {

		for j := 0; j < len(config.InputParams); j++ {
			wg.Add(1)
			ip := config.InputParams[j]
			go execScenario(config.AccessToken, config.WorkbookItemID, &wg, ip.MemCnt, ip.RecCnt, ip.Curr)
		}

		wg.Wait()

	}

}


func execScenario(accessToken string, workbookItemID string, wg *sync.WaitGroup, memCnt int, recCnt int, curr string) {

	defer wg.Done()

	// Create session
	sessionID, err := createSession(accessToken, workbookItemID)

	if err != nil {
		fmt.Printf("Unable to create a session. Access Token expired??\n")
		return
	}

	// Patch call
	sendInputParams(accessToken, workbookItemID, sessionID, memCnt, recCnt, curr)

	// Get call
	readOutput(accessToken, workbookItemID, sessionID, memCnt, recCnt, curr)
}

func sendInputParams(accessToken string, workbookItemID string, sessionID string, memCnt int, recCnt int, curr string) {

	s := fmt.Sprintf(`{"values" : [["%d"],["%d"],["%s"]] }`, memCnt, recCnt, curr)
	u := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/workbook/worksheets('InputOutput')/range(address='C3:C5')", workbookItemID)

	requestBody := ioutil.NopCloser(strings.NewReader(s))

	patchURL, _ := url.Parse(u)

	req := &http.Request{
		Method: "PATCH",
		URL:    patchURL,
		Header: map[string][]string{
			"Content-Type":        {"application/json; charset=UTF-8"},
			"Authorization":       {"Bearer " + accessToken},
			"workbook-session-id": {sessionID},
		},
		Body: requestBody,
	}

	_, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}
}

func createSession(accessToken string, workbookItemID string) (string, error) {

	u := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/workbook/createSession", workbookItemID)

	requestBody := ioutil.NopCloser(strings.NewReader(`{"persistChanges": false }`))
	sessionURL, _ := url.Parse(u)

	req := &http.Request{
		Method: "POST",
		URL:    sessionURL,
		Header: map[string][]string{
			"Content-Type":  {"application/json; charset=UTF-8"},
			"Authorization": {"Bearer " + accessToken},
		},
		Body: requestBody,
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode == 401 {
		return "", fmt.Errorf("InvalidAuthenticationToken, please use a new Acess Token")
	}

	data, _ := ioutil.ReadAll(res.Body)

	res.Body.Close()

	var parsedResponse map[string]interface{}
	err = json.Unmarshal(data, &parsedResponse)

	sessionID := parsedResponse["id"].(string)

	return sessionID, nil
}

func readOutput(accessToken string, workbookItemID string, sessionID string, memCnt int, recCnt int, curr string) {

	u := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/workbook/worksheets('InputOutput')/range(address='C9:C11')", workbookItemID)
	getURL, _ := url.Parse(u)

	req := &http.Request{
		Method: "GET",
		URL:    getURL,
		Header: map[string][]string{
			"Content-Type":        {"application/json; charset=UTF-8"},
			"Authorization":       {"Bearer " + accessToken},
			"workbook-session-id": {sessionID},
		},
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(res.Body)

	res.Body.Close()

	var parsedResponse ResponsePayload
	err = json.Unmarshal(data, &parsedResponse)

	fmt.Printf("Input - [%d, %d, %s]  Result - %+v\n", memCnt, recCnt, curr, parsedResponse.Values)

}

func GetDefaultInput() []InputParam{

	var result []InputParam

	result = append(result, InputParam{ MemCnt: 230, RecCnt: 79, Curr: "USD"})
	result = append(result, InputParam{ MemCnt: 1230, RecCnt: 79, Curr: "CAD"})
	result = append(result, InputParam{ MemCnt: 12300, RecCnt: 260, Curr: "GBP"})
	result = append(result, InputParam{ MemCnt: 36900, RecCnt: 749, Curr: "EUR"})
	result = append(result, InputParam{ MemCnt: 2500, RecCnt: 441, Curr: "AUD"})

	return result
}