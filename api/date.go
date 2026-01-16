package api

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

func NewBaseDate(t time.Time) openapi_types.Date {
	return openapi_types.Date{Time: t}
}
