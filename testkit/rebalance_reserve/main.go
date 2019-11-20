package main

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"log"
)

func main() {
	rebalanceReserve()
}

func rebalanceReserve() {
	client, err := vclient.NewDefaultTestNetClient("dev-velo-cen-node-01.velo.org:8080", "SBIK3BFXGBY25LUGYX42EYDZHHNRXDZJX3PTNB6VWVRDGNBIROSBZMFM")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	result, err := client.RebalanceReserve(context.Background(), vtxnbuild.RebalanceReserve{})
	if err != nil {
		panic(err)
	}
	log.Println(result.HorizonResult.TransactionSuccessToString())
}
