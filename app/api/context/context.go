package context

import (
	"zinktray/app/storage"
)

type RequestHandlerContext struct {
	Store *storage.Storage
}
