package main

import (
	"context"
	"fmt"
	"log"

	"github.com/0xfakespike/aimo-go-sdk/aimo"
)

func main() {
	client := aimo.NewClient(context.Background(), "https://devnet.aimo.network/api/v1/", "your-api-key")
	resp, err := client.Chat("Hello, how are you?")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Choices)
}
