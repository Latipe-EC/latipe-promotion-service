package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	server "latipe-promotion-services/internal"
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
		if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
			fmt.Printf("%s", err)
		}
	}()

	wg.Wait()
}
