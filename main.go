package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"git.ont.io/ontid/otf/config"
	"git.ont.io/ontid/otf/middleware"
	"git.ont.io/ontid/otf/rest"
	"git.ont.io/ontid/otf/service"
	store "git.ont.io/ontid/otf/store/leveldb"
	"git.ont.io/ontid/otf/utils"
	"git.ont.io/ontid/otf/vdri/ontdid"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/urfave/cli"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "agent otf"
	app.Action = startAgent
	app.Flags = []cli.Flag{
		utils.LogLevelFlag,
		utils.HttpIpFlag,
		utils.HttpPortFlag,
		utils.ChainAddrFlag,
	}
	app.Commands = []cli.Command{
		utils.DidCommand,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}
func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func startAgent(ctx *cli.Context) {
	ontSdk := sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(ctx.GlobalString(utils.GetFlagName(utils.ChainAddrFlag)))
	account, err := utils.OpenAccount(utils.DEFAULT_WALLET_PATH, ontSdk)
	if err != nil {
		panic(err)
	}
	port := ctx.GlobalString(utils.GetFlagName(utils.HttpPortFlag))
	ip := ctx.GlobalString(utils.GetFlagName(utils.HttpIpFlag))
	prov := store.NewProvider(utils.DEFAULT_STORE_DIR)
	db, err := prov.OpenStore(utils.DEFAULT_STORE_DIR)
	if err != nil {
		panic(err)
	}
	cfg := &config.Cfg{
		Port: port,
		Ip:   ip,
	}
	r := rest.InitRouter()
	ontvdri := ontdid.NewOntVDRI()
	msgSvr := service.NewMessageService(ontvdri)
	rest.NewService(account, cfg, db, msgSvr, ontvdri)
	middleware.Log.Infof("start agent svr%s,port:%s", account.Address, cfg.Port)
	startPort := ip + ":" + port
	err = r.Run(startPort)
	if err != nil {
		panic(err)
	}
	signalHandle()
}

func signalHandle() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Println("get a signal: stop the rest gateway process", si.String())
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
