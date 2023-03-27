package demojob

import (
	"context"
	"log"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func CompleteJob123(jobKey int64) {
	//建立本地端
	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: "localhost:26500"}
	client, err := zbc.NewClient(&config)
	if err != nil {
		panic(err)
	}

	variables := make(map[string]interface{})
	variables["status"] = true

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	log.Println("Completing job", jobKey)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Print("Successfully completed job")
}
