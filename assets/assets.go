package assets

import "embed"

//go:embed css/* js/*
var Public embed.FS
