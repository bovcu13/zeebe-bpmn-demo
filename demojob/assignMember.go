package demojob

import (
	"context"
	"fmt"
	"log"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func Assign() {
	//建立本地端
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

	//新增變數
	variables := make(map[string]interface{})
	variables["case"] = "no.5"

	//set ProcessDefinitionKey
	client.NewCreateInstanceCommand().ProcessDefinitionKey(2251799813721383)

	jobWorker := client.NewJobWorker().JobType("io.camunda.zeebe:userTask").Handler(handleJobAssign).Open()

	<-readyClose
	jobWorker.Close()
	jobWorker.AwaitClose()
}

func handleJobAssign(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		failjob(client, job)
		return
	}

	variables["name"] = "pin"
	variables["time"] = "2023-03-25 13:25:12"
	variables["case"] = "no.5"

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		failjob(client, job)
		return
	}

	log.Println("Completed job", job.GetKey())

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Print("Successfully completed job")
}
