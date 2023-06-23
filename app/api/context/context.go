package context

import (
	"go-fake-smtp/app/storage"
)

type RequestHandlerContext struct {
	Store *storage.Storage
}
