package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thoeni/google-homebase/pkg/apple"
)

func main() {
	username := flag.String("username", "", "Apple account username")
	password := flag.String("password", "", "Apple account password")
	c := apple.NewClient(*username, *password)
	var d apple.Device
	var user string
	if err := apple.FindDevice(c, "iPhone X", &user, &d); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v - %+v", user, d)
}
