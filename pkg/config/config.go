package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/Pineapple217/mb/pkg/auth"
	"github.com/joho/godotenv"
)

const (
	defaultLogo    = " • ▌ ▄ ·. ▄▄▄▄· \n ·██ ▐███▪▐█ ▀█▪\n ▐█ ▌▐▌▐█·▐█▀▀█▄\n ██ ██▌▐█▌██▄▪▐█\n ▀▀  █▪▀▀▀·▀▀▀▀ "
	defaultLink    = "https://mb.dev"
	defaultRights  = "mb.dev"
	defaultMessage = "Created without any JS."
	defaultHost    = "http://localhost:3000"

	DataDir   = "./data"
	UploadDir = DataDir + "/uploads"
	BackupDir = DataDir + `/backups`
)

var (
	OutputTimezone  *time.Location
	HomepageLogo    string
	HomepageLink    string
	HomepageRights  string
	HomepageMessage string
	Host            string
)

func Load() {
	slog.Info("Loading configs")
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found")
	}
	password, isSet := os.LookupEnv("MB_AUTH_PASSWORD")
	if !isSet {
		slog.Info("AUHT_PASSWORD is not set. Using random password")
		slog.Info("Random password generated", "password", auth.SecretPassword)
	} else {
		slog.Info("Auth password succesfully set")
		auth.SecretPassword = password

	}
	initTimezone()
	initLogo()
	initLink()
	initRights()
	initMessage()
	initHost()
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

func initHost() {
	host, isSet := os.LookupEnv("MB_HOST")
	if !isSet {
		Host = defaultHost
		return
	}
	Host = host
}
