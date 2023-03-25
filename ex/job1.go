package demojob

import (
	"context"
	"fmt"
	"log"
	"zedgodemo/util"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

var readyClose = make(chan struct{})

func JobWorkOne() {
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

	// // deploy process

	response, err := client.NewDeployProcessCommand().AddResourceFile("subtask.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.String())

	// create a new process instance
	variables := make(map[string]interface{})
	variables["orderId"] = util.RandomInt(10001, 100000)
	variables["isAddAppreved"] = true

	request, err := client.NewCreateInstanceCommand().BPMNProcessId("Process_1mh06z5").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	result, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())

	jobWorker := client.NewJobWorker().JobType("io.camunda.zeebe:userTask").Handler(handleJob).Open()

	<-readyClose
	jobWorker.Close()
	jobWorker.AwaitClose()
}

func handleJob(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		failjob(client, job)
		return
	}

	variables["name"] = "Bod"
	variables["age"] = 30
	variables["time"] = "2021-01-01"

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		failjob(client, job)
		return
	}

	log.Println("Completed job", job.GetKey())
	log.Println("Processing order:", variables["name"])

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Print("Successfully completed job")
	close(readyClose)
}

func failjob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to handle job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(3).Send(ctx)
	if err != nil {
		panic(err)
	}

}
