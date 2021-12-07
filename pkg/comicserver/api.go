package comicserver

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/klaital/comics/pkg/config"
	"github.com/klaital/comics/pkg/filters"
)

func New(cfg config.Config) *restful.Container {
	c := restful.NewContainer()

	// Set global filters here
	// TODO: add rate limiting
	c.Filter(filters.RequestLogFilter)

	// Comics API
	comicsWS := restful.WebService{}
	comicsWS.Path(cfg.BasePath + "/comics").ApiVersion("1.0.0").Doc("CRUD API for Comics subscriptions")

	c.Add(&comicsWS)

	return c
}
