package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/Pineapple217/mb/pkg/auth"
	"github.com/joho/godotenv"
)

const (
	DataDir   = "./data"
	UploadDir = DataDir + "/uploads"
	BackupDir = DataDir + "/backups"
)

var (
	OutputTimezone  *time.Location
	HomepageLogo    = " • ▌ ▄ ·. ▄▄▄▄· \n ·██ ▐███▪▐█ ▀█▪\n ▐█ ▌▐▌▐█·▐█▀▀█▄\n ██ ██▌▐█▌██▄▪▐█\n ▀▀  █▪▀▀▀·▀▀▀▀ "
	HomepageLink    = "https://mb.dev"
	HomepageRights  = "mv.dev"
	HomepageMessage = "Created without any JS."
	Host            = "http://localhost:3000"
	Debug           = false
	NavidromePrefix string
)

func Load() {
	slog.Info("Loading configs")
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}

	setAuthPassword()
	setEnvVar("MB_LOGO", &HomepageLogo)
	setEnvVar("MB_LINK", &HomepageLink)
	setEnvVar("MB_RIGHTS", &HomepageRights)
	setEnvVar("MB_MESSAGE", &HomepageMessage)
	setEnvVar("MB_HOST", &Host)
	setDebug()
	setTimezone()
	setEnvVar("MB_NAVIDROME_PREFIX", &NavidromePrefix)
}

func setAuthPassword() {
	if password, ok := os.LookupEnv("MB_AUTH_PASSWORD"); ok {
		auth.SecretPassword = password
		slog.Info("Auth password successfully set")
	} else {
		slog.Info("AUTH_PASSWORD is not set. Using random password", "password", auth.SecretPassword)
	}
}

func setEnvVar(key string, target *string) {
	if val, ok := os.LookupEnv(key); ok {
		*target = val
	}
}

func setDebug() {
	if val, ok := os.LookupEnv("MB_DEBUG"); ok {
		if debug, err := strconv.ParseBool(val); err == nil {
			Debug = debug
		}
	}
}

func setTimezone() {
	if locStr, ok := os.LookupEnv("MB_TIMEZONE"); ok {
		if loc, err := time.LoadLocation(locStr); err == nil {
			OutputTimezone = loc
			return
		}
	}
	OutputTimezone = time.Now().Location()
}
