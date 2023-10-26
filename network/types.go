package network

var (
	RPCServer     *Rpc
	GatewayServer *Gateway
)

func Init() {
	RPCServer = (&Rpc{}).New()
	GatewayServer = (&Gateway{}).New()
}
