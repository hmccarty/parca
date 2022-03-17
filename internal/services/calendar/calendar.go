package calendar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	m "github.com/hmccarty/parca/internal/models"
	c "github.com/hmccarty/parca/internal/services/config"
)

const (
	getCalendarDataURL       = "https://www.googleapis.com/calendar/v3/calendars/%s?%s"
	getCalendarEventsListURL = "https://www.googleapis.com/calendar/v3/calendars/%s/events?%s"
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

func (client *GoogleCalendarClient) GetCalendarEvents(calendarID string, end time.Time) ([]m.CalendarEventData, error) {
	fmt.Println("Starting events request")
	if time.Now().After(client.expiration) {
		token, expiration, err := requestToken(client.serviceEmail, client.keyData)
		if err != nil {
			return nil, err
		}
		client.token = token
		client.expiration = expiration
	}

	data := url.Values{}
	data.Set("access_token", client.token)
	data.Set("timeMin", time.Now().UTC().Format(time.RFC3339))
	data.Set("timeMax", end.Format(time.RFC3339))

	url := fmt.Sprintf(getCalendarEventsListURL, calendarID, data.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println("Sent request")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var calendarEventsDataRaw map[string]interface{}
	if err := json.Unmarshal(body, &calendarEventsDataRaw); err != nil {
		return nil, err
	}

	if calendarEventsDataRaw["items"] == nil {
		return nil, errors.New("calendar events not found")
	}

	fmt.Println("Starting to read in response data")
	eventItemsRaw := calendarEventsDataRaw["items"].([]interface{})
	var calendarEventsData []m.CalendarEventData
	for _, v := range eventItemsRaw {
		eventDataRaw := v.(map[string]interface{})

		var summary string
		if eventDataRaw["summary"] != nil {
			summary = eventDataRaw["summary"].(string)
		} else {
			continue
		}

		calendarEventsData = append(calendarEventsData,
			m.CalendarEventData{
				Name: summary,
			})
	}
	return calendarEventsData, nil
}
