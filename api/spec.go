// Package api embeds the scaffold OpenAPI document.
package api

import _ "embed"

//go:embed openapi.json
var Document []byte
