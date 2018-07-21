package main

import (
	"bytes"
	"fmt"
	"time"

	"google.golang.org/api/dialogflow/v2"
)

func HomeSuccessResponse(user string, d Device, locale string) dialogflow.GoogleCloudDialogflowV2WebhookResponse {

	var resp string
	switch locale {
	case "it":
		resp = it(user, d)
	default:
		resp = en(user, d)
	}

	return dialogflow.GoogleCloudDialogflowV2WebhookResponse{
		FulfillmentText: resp,
	}
}

func HomeFailureResponse(message string) dialogflow.GoogleCloudDialogflowV2WebhookResponse {
	return dialogflow.GoogleCloudDialogflowV2WebhookResponse{
		FulfillmentText: message,
	}
}

func translateMonth(m time.Month) string {
	switch m {
	case time.January:
		return "Gennaio"
	case time.February:
		return "Febbraio"
	case time.March:
		return "Marzo"
	case time.April:
		return "Aprile"
	case time.May:
		return "Maggio"
	case time.June:
		return "Giugno"
	case time.July:
		return "Luglio"
	case time.August:
		return "Agosto"
	case time.September:
		return "Settembre"
	case time.October:
		return "Ottobre"
	case time.November:
		return "Novembre"
	case time.December:
		return "Dicembre"
	default:
		return ""
	}
}

func it(user string, d Device) string {
	var w bytes.Buffer

	w.WriteString(fmt.Sprintf("%s ", user))

	if d.isHome() {
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

func en(user string, d Device) string {
	var w bytes.Buffer

	w.WriteString(fmt.Sprintf("%s ", user))

	if d.isHome() {
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
