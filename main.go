package main

import (
	"context"
	"fmt"
	"log"
	"task/service"
	"time"
)

func main() {
	timeRange := 5 * time.Second
	maxRequests := 6
	mfc := service.NewMemFloodControl(timeRange, maxRequests)

	for i := 0; i < 9; i++ {
		ok, err := mfc.Check(context.Background(), 777)
		if err != nil {
			log.Fatalln(err)
		}

		if !ok {
			fmt.Println("Flood control not passed")
		} else {
			fmt.Println("Flood control passed")
		}
	}
}
