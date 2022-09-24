package main

import (
	config "github.com/content-services/content-sources-backend/pkg/config"
	"github.com/content-services/content-sources-backend/pkg/db"
	"github.com/content-services/content-sources-backend/pkg/event"
	"github.com/content-services/content-sources-backend/pkg/event/handler"
)

func main() {
	cfg := config.Get()
	db.Connect()
	handler := handler.NewIntrospectHandler(db.DB)
	event.Start(cfg, handler)
}