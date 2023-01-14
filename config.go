package muse

import (
	"embed"
	"time"
)

// RanAt is the time when the executable started running.
var RanAt = time.Now().Unix()

//go:embed templates/**/*.tmpl assets/**/*
var FS embed.FS
