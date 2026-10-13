package data

import "embed"

// This file and the package it creates are intended to distribute the static
// assets of this data directory in the go module, specifically for consumption
// by the ARO wrapper on a temporary basis. The installer itself does not import
// this package and no other codebases should either, as we intend to remove it
// once ARO has moved off the wrapper.

//go:embed all:bootstrap manifests/* coreos/*
var _ embed.FS
