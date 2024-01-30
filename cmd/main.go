package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	server "latipe-promotion-services/internal"
	"net"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("Init application")
	defer log.Fatalf("[Info] Application has closed")
	numCPU := runtime.NumCPU()
	fmt.Printf("Number of CPU cores: %d\n", numCPU)

	serv, err := server.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	//subscriber
	var wg sync.WaitGroup
	{
		wg.Add(2)
		go serv.CommitPurchaseTransactionSubscriber().ListenProductPurchaseCreate(&wg)
		go serv.RollbackPurchaseTransactionSubscriber().ListenProductPurchaseCreate(&wg)
	}

	//api handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.App().Listen(serv.Config().Server.RestPort); err != nil {
			fmt.Printf("%s", err)
		}
	}()

	//grpc handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Infof("Start grpc server on port: localhost%v", serv.Config().GRPC.Port)
		lis, err := net.Listen("tcp", serv.Config().GRPC.Port)
		if err != nil {
			log.Fatalf("failed to listen: %v\n", err)
		}

		if err := serv.GrpcServ().Serve(lis); err != nil {
			log.Infof("%s", err)
		}
	}()

	wg.Wait()
}
