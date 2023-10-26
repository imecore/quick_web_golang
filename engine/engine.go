package engine

import (
	"quick_web_golang/model"
	"quick_web_golang/network"
	"quick_web_golang/provider"
)

func Init() {
	provider.Init()
	network.Init()
}

func Start() {
	provider.Database.Start()
	provider.Cache.Start()
	provider.SessionManager.Start()

	model.Repos = model.NewRepo()

	go network.RPCServer.Start()
	go network.GatewayServer.Start()

}

func Stop() {
	provider.Database.Close()
	provider.Cache.Close()
	provider.SessionManager.Close()
	network.GatewayServer.Close()
	network.RPCServer.Close()
}
