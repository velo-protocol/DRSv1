package main

import (
	"context"
	cenGrpc "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/client"
	"log"
)

func main() {

	client, err := vclient.NewDefaultTestNetClient("localhost:8080", "SC6OED46ETTH56IZ2YLCHTKWQZTIK3VFB54HXRSWSJMPBVSTK2H7YE7V")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	getCollateralHealthCheck, err := client.GetCollateralHealthCheck(context.Background(), &cenGrpc.Empty{})
	if err != nil {
		panic(err)
	}
	log.Println("Asset Code: ", getCollateralHealthCheck.AssetCode)
	log.Println("Asset Issuer: ", getCollateralHealthCheck.AssetIssuer)
	log.Println("Required Amount: ", getCollateralHealthCheck.RequiredAmount)
	log.Println("Pool Amount: ", getCollateralHealthCheck.PoolAmount)
}
