package config

import (
	"fmt"
	"os"
	"time"

	"github.com/Pineapple217/mb/auth"
	"github.com/joho/godotenv"
)

const (
	defaultLogo    = " • ▌ ▄ ·. ▄▄▄▄· \n ·██ ▐███▪▐█ ▀█▪\n ▐█ ▌▐▌▐█·▐█▀▀█▄\n ██ ██▌▐█▌██▄▪▐█\n ▀▀  █▪▀▀▀·▀▀▀▀ "
	defaultLink    = "https://mb.dev"
	defaultRights  = "mb.dev"
	defaultMessage = "Created without any JS."

	DataDir = "./data"
)

var (
	OutputTimezone  *time.Location
	HomepageLogo    string
	HomepageLink    string
	HomepageRights  string
	HomepageMessage string
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}
	password, isSet := os.LookupEnv("MB_AUTH_PASSWORD")
	if !isSet {
		fmt.Println("AUHT_PASSWORD is not set. Using random password")
		fmt.Printf("Auth password: %s\n", auth.SecretPassword)
	} else {
		fmt.Println("Auth password succesfully set")
		auth.SecretPassword = password

	}
	initTimezone()
	initLogo()
	initLink()
	initRights()
	initMessage()
}

func initTimezone() {
	// MAYBE: get timezone from user
	locStr, isSet := os.LookupEnv("MB_TIMEZONE")
	if !isSet {
		OutputTimezone = time.Now().Location()
		return
	}
	envLoc, err := time.LoadLocation(locStr)
	if err != nil {
		OutputTimezone = time.Now().Location()
		return
	}
	OutputTimezone = envLoc

}

func initLogo() {
	logo, isSet := os.LookupEnv("MB_LOGO")
	if !isSet {
		HomepageLogo = defaultLogo
		return
	}
	HomepageLogo = logo
}

func initLink() {
	link, isSet := os.LookupEnv("MB_LINK")
	if !isSet {
		HomepageLink = defaultLink
		return
	}
	HomepageLink = link
}

func initRights() {
	rights, isSet := os.LookupEnv("MB_RIGHTS")
	if !isSet {
		HomepageRights = defaultRights
		return
	}
	HomepageRights = rights
}

func initMessage() {
	message, isSet := os.LookupEnv("MB_MESSAGE")
	if !isSet {
		HomepageMessage = defaultMessage
		return
	}
	HomepageMessage = message
}
