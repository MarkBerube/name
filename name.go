package name

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// SheetConfig is a struct that holds configuration values on how the name spreadsheet is requested.
// Configuration only supports Google's V4 sheet API currently, please see README for appropriate default values.
type SheetConfig struct {
	URL               string // URL to the spreadsheet API
	ID                string // ID of the spreadsheet
	APIKey            string // API key to call the spreadsheet API
	Name              string // Name of the sheet in the spreadsheet that has the name list
	MsgLimit          int    // Maximum length of strings returned by GetRandomNameList
	MsgHeader         string // Header attached to the top of the first message string for GetRandomNameList
	Titles            bool   // Allow titles to be used by GetRandomNameList
	SecondNameAppends bool   // Allow second name appends to be used by GetRandomNameList
}

// NewConfig creates a new config struct for the name generator calls.
// Returns a SheetConfig as a reference that can be used to call GetRandomNameList().
func NewConfig(url string, id string, apiKey string, name string, msgLimit int, msgHeader string) *SheetConfig {
	return &SheetConfig{URL: url,
		ID:                id,
		APIKey:            apiKey,
		Name:              name,
		MsgLimit:          msgLimit,
		MsgHeader:         msgHeader,
		Titles:            true,
		SecondNameAppends: true}
}

// IgnoreTitles sets the name generator to not return names with titles.
// By default titles are turned on.
// Should be turned off if your name list doesn't include an appropriate range on label C.
func (request *SheetConfig) IgnoreTitles() {
	request.Titles = false
}

// IgnoreSecondNameAppends sets the name generator to not return last names with extras words/phrases.
// By default second name appends are turned on.
// Should be turned off if your name list doesn't include an appropriate range on label D.
func (request *SheetConfig) IgnoreSecondNameAppends() {
	request.SecondNameAppends = false
}

// GetRandomNameList returns an array of strings that has randomly generated names.
// totalNumber controls the amount of total names returned in the resultant array.
func GetRandomNameList(totalNumber int, config *SheetConfig, c APIClient) []string {
	var bufferResult, firstNameList, lastNameList, titleList, lastNameAppendList []string
	var buffer string

	firstNameList = c.getNameList("A2:A", config)
	lastNameList = c.getNameList("B2:B", config)

	if config.Titles {
		titleList = c.getNameList("C2:C", config)
	}

	if config.SecondNameAppends {
		lastNameAppendList = c.getNameList("D2:D", config)
	}

	firstUsed := make(map[int]bool)
	lastUsed := make(map[int]bool)
	titleUsed := make(map[int]bool)

	buffer += config.MsgHeader

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < totalNumber; i++ {
		firstNameIndex := getNewRandomValue(len(firstNameList), firstUsed)
		firstName := firstNameList[firstNameIndex]

		lastNameIndex := getNewRandomValue(len(lastNameList), lastUsed)
		lastName := lastNameList[lastNameIndex]

		if config.SecondNameAppends && rand.Intn(10) > 3 {
			lastNameAppendIndex := getNewRandomValueWithRepeats(len(lastNameAppendList))
			lastName += lastNameAppendList[lastNameAppendIndex]
		}

		if config.Titles && rand.Intn(10) > 3 {
			titleNameIndex := getNewRandomValue(len(titleList), titleUsed)
			title := titleList[titleNameIndex]

			buffer += fmt.Sprintf("%s %s the %s", firstName, lastName, title) + "\n"
		} else {
			buffer += fmt.Sprintf("%s %s", firstName, lastName) + "\n"
		}
	}

	bufferResult = append(bufferResult, buffer)
	return bufferResult
}

// getNewRandomValue returns a new random index number for a list that hasn't been used yet.
// Do note that if you request 100 names generated and only have a 50 words in a list this will lock.
func getNewRandomValue(maxValue int, valuesUsed map[int]bool) int {
	for {
		newInt := rand.Intn(maxValue)
		_, ok := valuesUsed[newInt]
		if !ok {
			valuesUsed[newInt] = true
			return newInt
		}
	}
}

// getNewRandomValueWithRepeats returns a new random index number that could possibly be used already.
// This is useful for second append lists that will likely have less than 10 values.
func getNewRandomValueWithRepeats(maxValue int) int {
	return rand.Intn(maxValue)
}

// Sheet is a struct model for JSON decoding Google's V4 sheet API.
type Sheet struct {
	Values [][]string
}

// APIClient is the interface for the API client making name request calls.
type APIClient interface {
	getNameList(listRange string, config *SheetConfig) []string
}

// GoogleSheetAPIClient is the client struct making calls out to the Google Sheet API for name lists.
type GoogleSheetAPIClient struct{}

func (g GoogleSheetAPIClient) getNameList(listRange string, config *SheetConfig) []string {
	url := fmt.Sprintf("%s%svalues/%s!%s/?key=%s", config.URL, config.ID, config.Name, listRange, config.APIKey)
	list := g.concatAppend(g.makeSheetRequest(url).Values)
	return list
}

func (g GoogleSheetAPIClient) makeSheetRequest(url string) Sheet {
	nameClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	res, getErr := nameClient.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	body, _ := ioutil.ReadAll(res.Body)
	feed := Sheet{}

	json.Unmarshal(body, &feed)
	return feed
}

func (g GoogleSheetAPIClient) concatAppend(slices [][]string) []string {
	var tmp []string
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}
