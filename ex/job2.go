package testjob

import (
	"context"
	"fmt"
	"log"

	"zedgodemo/util"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func JobWorkTwo() {
	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: "localhost:26500"}
	client, err := zbc.NewClient(&config)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	topology, err := client.NewTopologyCommand().Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(topology.String())

	variables := make(map[string]interface{})
	variables["orderId"] = util.RandomInt(10000, 100000)

	client.NewCreateInstanceCommand().ProcessDefinitionKey(2251799813689018).VariablesFromMap(variables)
	// // deploy process

	jobWorker := client.NewJobWorker().JobType("io.camunda.zeebe:userTask").Handler(handleJobTwo).Open()

	<-readyClose
	jobWorker.Close()
	jobWorker.AwaitClose()
}

func handleJobTwo(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		failjob(client, job)
		return
	}

	variables["time"] = "2023-03-03 12:12:12"
	variables["name"] = "zedgo"
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		failjob(client, job)
		return
	}

	log.Println("Completed job", job.GetKey())
	log.Println("Processing order:", variables["isStart"])
	log.Println("is Start", variables["isStart"])

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Print("Successfully completed job")
	close(readyClose)
}
