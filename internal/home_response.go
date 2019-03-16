package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/translate"
	"github.com/thoeni/google-homebase/pkg/apple"
	"golang.org/x/text/language"
	"google.golang.org/api/dialogflow/v2"
	"google.golang.org/api/option"
)

// HomeSuccessResponse returns a successful response to Google Home dialogflow
// API by injecting data extracted from the device into the FulFillmentText
// attribute
func HomeSuccessResponse(user string, d apple.Device, locale string) dialogflow.GoogleCloudDialogflowV2beta1WebhookResponse {
	var resp string

	switch locale {
	case "it":
		resp = it(user, d)
	case "en", "en-gb", "en-us":
		resp = en(user, d)
	default:
		resp = any(locale[:2], user, d)
	}

	return dialogflow.GoogleCloudDialogflowV2beta1WebhookResponse{
		FulfillmentText: resp,
	}
}

// HomeFailureResponse is a helper to create a failure response. It just inject
// the message in input into the FulfillmentText attribute
func HomeFailureResponse(message string) dialogflow.GoogleCloudDialogflowV2beta1WebhookResponse {
	return dialogflow.GoogleCloudDialogflowV2beta1WebhookResponse{
		FulfillmentText: message,
	}
}

func translateMonth(m time.Month) string {

	var itMonths = map[time.Month]string{
		time.January:   "Gennaio",
		time.February:  "Febbraio",
		time.March:     "Marzo",
		time.April:     "Aprile",
		time.May:       "Maggio",
		time.June:      "Giugno",
		time.July:      "Luglio",
		time.August:    "Agosto",
		time.September: "Settembre",
		time.October:   "Ottobre",
		time.November:  "Novembre",
		time.December:  "Dicembre",
	}

	if itMonth, exists := itMonths[m]; exists {
		return itMonth
	}
	return m.String()
}

func it(user string, d apple.Device) string {
	var w bytes.Buffer

	w.WriteString(fmt.Sprintf("%s ", user))

	if IsHome(d) {
		w.WriteString("molto probabilmente si trova a casa.\n")
	} else {
		w.WriteString("non si trova a casa in questo momento, o dista almeno 25 metri dal suo divano.\n")
	}

	italy, _ := time.LoadLocation("Europe/Rome")
	lastSeen := time.Unix(d.Location.Timestamp/1000, 0).In(italy)
	w.WriteString(fmt.Sprintf("Ultimo aggiornamento delle %d e %d del giorno %d %v.\n", lastSeen.Hour(), lastSeen.Minute(), lastSeen.Day(), translateMonth(lastSeen.Month())))

	if d.BatteryLevel > 0 {
		w.WriteString(fmt.Sprintf("Stato batteria del suo cellulare: %.f percento.\n", d.BatteryLevel*100))
	} else {
		w.WriteString("Attualmente non risulta possibile rilevare la carica della sua batteria.\n")
	}

	return w.String()
}

func en(user string, d apple.Device) string {
	var w bytes.Buffer

	w.WriteString(fmt.Sprintf("%s ", user))

	if IsHome(d) {
		w.WriteString("is likely to be at home right now.\n")
	} else {
		w.WriteString("is not at home right now, or he's more than 25 meters far from his sofa.\n")
	}

	uk, _ := time.LoadLocation("Europe/London")
	lastSeen := time.Unix(d.Location.Timestamp/1000, 0).In(uk)
	w.WriteString(fmt.Sprintf("Last updated at %d %d of %v, %d of %v.\n", lastSeen.Hour(), lastSeen.Minute(), lastSeen.Weekday(), lastSeen.Day(), lastSeen.Month()))

	if d.BatteryLevel > 0 {
		w.WriteString(fmt.Sprintf("His phone battery is charged at %.f percent.\n", d.BatteryLevel*100))
	} else {
		w.WriteString("His phone battery status is currently unavailable.\n")
	}

	return w.String()
}

func any(locale string, user string, d apple.Device) string {
	en := en(user, d)
	if r, err := translateText(locale, en); err == nil {
		return r
	}

	return en
}

func translateText(targetLanguage, text string) (string, error) {
	gTranslateKey := os.Getenv("TRANSLATE_KEY")
	if gTranslateKey == "" {
		return text, nil
	}

	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		fmt.Println("Error when parsing language:", err)
		return "", err
	}

	client, err := translate.NewClient(ctx, option.WithAPIKey(gTranslateKey))
	if err != nil {
		fmt.Println("Error when creating translate client:", err)
		return "", err
	}
	defer func() {
		_ = client.Close()
	}()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		fmt.Println("Error when translating:", err)
		return "", err
	}

	return resp[0].Text, nil
}
