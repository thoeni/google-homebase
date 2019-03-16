package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thoeni/google-homebase/internal"
	"github.com/thoeni/google-homebase/pkg/apple"
	"google.golang.org/api/dialogflow/v2"
)

func handleRequest(r dialogflow.GoogleCloudDialogflowV2beta1WebhookRequest) (dialogflow.GoogleCloudDialogflowV2beta1WebhookResponse, error) {

	fmt.Printf("Request:\n%+v\n", r)
	fmt.Println("Query result:", r.QueryResult)

	locale := "en-gb"

	if r.QueryResult != nil {
		fmt.Println("Language:", r.QueryResult.LanguageCode)
		locale = r.QueryResult.LanguageCode
		fmt.Println("Params:", string(r.QueryResult.Parameters))
	}

	var err error

	enCreds := os.Getenv("CREDS")
	deCreds, err := internal.DecryptEnvCredentials(enCreds)
	if err != nil {
		return internal.HomeFailureResponse(err.Error()), err
	}

	c := apple.NewClient(deCreds["username"], deCreds["password"])

	var d apple.Device
	var user string
	for i := 0; i < 3; i++ {
		fmt.Println("Calling iCloud, iteration", i)
		err = apple.FindDevice(c, "iPhone X", &user, &d)
		if err != nil {
			fmt.Println("Error was:", err)
			return internal.HomeFailureResponse("Something went wrong while retrieving the data"), err
		}

		if d.Location.Outdated {
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}
	}

	return internal.HomeSuccessResponse(user, d, locale), nil
}

func main() {
	lambda.Start(handleRequest)
}
