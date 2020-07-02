package utils

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/urfave/cli"
)

var DidCommand = cli.Command{
	Name:        "did",
	Usage:       "new did",
	Description: "Did management commands can generate did",
	Action:      newdid,
	Flags: []cli.Flag{
		RPCPortFlag,
		TransactionGasPriceFlag,
		TransactionGasLimitFlag,
		WalletFileFlag,
	},
}

func newdid(ctx *cli.Context) error {
	ontSdk := sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(ctx.GlobalString(GetFlagName(RPCPortFlag)))
	gasPrice := ctx.Uint64(TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(TransactionGasLimitFlag.Name)
	optionFile := checkFileName(ctx)
	acc, err := OpenAccount(optionFile, ontSdk)
	if err != nil {
		return fmt.Errorf("open account err:%s", err)
	}
	did, err := NewDID(ontSdk, acc, gasPrice, gasLimit)
	if err != nil {
		return fmt.Errorf("new did err:%s", err)
	}
	fmt.Printf("did:%v", did)
	return nil
}

func checkFileName(ctx *cli.Context) string {
	if ctx.IsSet(GetFlagName(WalletFileFlag)) {
		return ctx.String(GetFlagName(WalletFileFlag))
	} else {
		//default account file name
		return DEFAULT_WALLET_FILE_NAME
	}
}

func NewDID(ontSdk *sdk.OntologySdk, acc *sdk.Account, gasPrice, gasLimit uint64) (string, error) {
	did, err := sdk.GenerateID()
	if err != nil {
		return "", err
	}
	err = RegisterDid(did, ontSdk, acc, gasPrice, gasLimit)
	if err != nil {
		return "", err
	}
	return did, nil
}

func RegisterDid(did string, ontSdk *sdk.OntologySdk, acc *sdk.Account, gasPrice, gasLimit uint64) error {
	if ontSdk.Native == nil || ontSdk.Native.OntId == nil {
		return fmt.Errorf("ontsdk is nil")
	}
	txHash, err := ontSdk.Native.OntId.RegIDWithPublicKey(gasPrice, gasLimit, acc, did, acc)
	if err != nil {
		return err
	}
	fmt.Printf("did:%v,hash:%v", did, txHash.ToHexString())
	return nil
}
