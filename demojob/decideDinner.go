package demojob

import (
	"context"
	"fmt"
	"log"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

var readyClose = make(chan struct{})

func DecideDinner() {
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

	// deploy process
	response, err := client.NewDeployResourceCommand().AddResourceFile("PreparingDinner.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.String())

	// create a new process instance
	variables := make(map[string]interface{})
	variables["meal"] = "Salad"
	//set ProcessId
	request, err := client.NewCreateInstanceCommand().BPMNProcessId("preparing-dinner").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	msg, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(msg.String())

	//runtask - Decide Dinner
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

	variables["meal"] = "Salad"

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		failjob(client, job)
		return
	}

	log.Println("Completed job", job.GetKey())
	//log.Println("Processing order:", variables["meal"])

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

// func roleToString(role pb.Partition_PartitionBrokerRole) string {
// 	switch role {
// 	case pb.Partition_LEADER:
// 		return "Leader"
// 	case pb.Partition_FOLLOWER:
// 		return "Follower"
// 	default:
// 		return "Unknown"
// 	}

// }
