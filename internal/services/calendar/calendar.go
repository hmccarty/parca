package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	m "github.com/hmccarty/parca/internal/models"
	c "github.com/hmccarty/parca/internal/services/config"
)

const (
	getCalendarDataURL = "https://www.googleapis.com/calendar/v3/calendars/%s?%s"
)

type GoogleCalendarClient struct {
	serviceEmail string
	keyData      []byte
	token        string
	expiration   time.Time
}

func NewGoogleCalendarClient(config *c.Config) m.CalendarClient {
	keyData, err := ioutil.ReadFile(config.GooglePrivatePEMFile)
	if err != nil {
		panic(err)
	}

	token, expiration, err := requestToken(
		config.GoogleServiceEmail, keyData)
	if err != nil {
		fmt.Println(err)
	}

	return &GoogleCalendarClient{
		serviceEmail: config.GoogleServiceEmail,
		keyData:      keyData,
		token:        token,
		expiration:   expiration,
	}
}

func (client *GoogleCalendarClient) GetCalendarData(calendarID string) (m.CalendarData, error) {
	if time.Now().After(client.expiration) {
		token, expiration, err := requestToken(client.serviceEmail, client.keyData)
		if err != nil {
			return m.CalendarData{}, err
		}
		client.token = token
		client.expiration = expiration
	}

	data := url.Values{}
	data.Set("access_token", client.token)

	url := fmt.Sprintf(getCalendarDataURL, calendarID, data.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return m.CalendarData{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return m.CalendarData{}, err
	}

	var calendarDataRaw map[string]interface{}
	if err := json.Unmarshal(body, &calendarDataRaw); err != nil {
		return m.CalendarData{}, err
	}

	var name string
	if calendarDataRaw["summary"] != nil {
		name = calendarDataRaw["summary"].(string)
	}

	var description string
	if calendarDataRaw["description"] != nil {
		description = calendarDataRaw["description"].(string)
	}

	return m.CalendarData{
		Name:        name,
		Description: description,
	}, nil
}

func (*GoogleCalendarClient) GetCalendarEvents(calendarID string) ([]m.CalendarEventData, error) {
	return []m.CalendarEventData{}, nil
}
