package cmd

import (
	"code.vegaprotocol.io/vega/logging"
	"github.com/vegaprotocol/data-metrics-store/clients/comet"
	"github.com/vegaprotocol/data-metrics-store/config"
	"github.com/vegaprotocol/data-metrics-store/services"
	"github.com/vegaprotocol/data-metrics-store/services/read"
	"github.com/vegaprotocol/data-metrics-store/services/update"
)

type AllServices struct {
	Config        *config.Config
	Log           *logging.Logger
	StoreService  *services.StoreService
	ReadService   *read.ReadService
	UpdateService *update.UpdateService
}

func SetupServices(configFilePath string, forceDebug bool) (svc AllServices, err error) {
	svc.Config, svc.Log, err = config.GetConfigAndLogger(configFilePath, forceDebug)
	if err != nil {
		return
	}

	svc.StoreService, err = services.NewStoreService(&svc.Config.SQLStore, svc.Log)
	if err != nil {
		return
	}

	cometClient := comet.NewCometClient(&svc.Config.CometBFT)
	svc.ReadService, err = read.NewReadService(cometClient, svc.StoreService, svc.Log)
	if err != nil {
		return
	}

	svc.UpdateService, err = update.NewUpdateService(svc.ReadService, svc.StoreService, svc.Log)
	if err != nil {
		return
	}
	return
}
