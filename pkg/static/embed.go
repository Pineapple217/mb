package static

import (
	"embed"
)

var (
	//go:embed bundle/*
	PublicFS embed.FS
)
