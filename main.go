package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/dialogflow/v2"
)

func handleRequest(r dialogflow.GoogleCloudDialogflowV2WebhookRequest) (dialogflow.GoogleCloudDialogflowV2WebhookResponse, error) {

	fmt.Println(r.QueryResult)
	fmt.Println("Language:", r.QueryResult.LanguageCode)
	fmt.Println("Params:", string(r.QueryResult.Parameters))

	var err error

	username := os.Getenv("UNAME")
	password := os.Getenv("PWD")
	err = decryptEnvCredentials(&username, &password)
	if err != nil {
		return HomeFailureResponse(err.Error()), err
	}

	c := NewClient(username, password)

	var d Device
	var user string
	for i := 0; i < 3; i++ {
		fmt.Println("Calling iCloud, iteration", i)
		err = FindDevice(c, "iPhone X", &user, &d)
		if err != nil {
			fmt.Println("Error was:", err)
			return HomeFailureResponse("Something went wrong while retrieving the data"), err
		}

		if d.Location.Outdated {
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}
	}

	return HomeSuccessResponse(user, d, r.QueryResult.LanguageCode), nil
}

func main() {
	lambda.Start(handleRequest)
}
