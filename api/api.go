package api

import (
	"context"
	_ "embed"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/numary/ledger/api/controllers"
	"github.com/numary/ledger/ledger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	controllers.Module,
)

type API struct {
	addr   string
	engine *gin.Engine
}

func NewAPI(
	lc fx.Lifecycle,
	resolver *ledger.Resolver,
	transactionController *controllers.TransactionController,
) *API {
	gin.SetMode(gin.ReleaseMode)

	cc := cors.DefaultConfig()
	cc.AllowAllOrigins = true
	cc.AllowCredentials = true
	cc.AddAllowHeaders("authorization")

	router := NewRoutes(cc, resolver, transactionController) //todo: use fx

	h := &API{
		engine: router,
		addr:   viper.GetString("server.http.bind_address"),
	}

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			go h.Start()
			return nil
		},
	})

	return h
}

func (h *API) Start() {
	h.engine.Run(h.addr)
}