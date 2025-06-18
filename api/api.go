package api

import (
	"github.com/holonet/core/database"
)

var dbHandler *database.DBHandler

func SetDBHandler(handler *database.DBHandler) {
	dbHandler = handler
}
