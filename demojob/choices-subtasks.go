package demojob

import (
	"context"
	"fmt"
	"log"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func SubTask() {
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

	// deploy process - 輸入.bpmn
	response, err := client.NewDeployResourceCommand().AddResourceFile("subtask.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.String())

	// create a new process instance
	variables := make(map[string]interface{})
	variables["item"] = "Task Assignment"
	//set ProcessId
	request, err := client.NewCreateInstanceCommand().BPMNProcessId("choices-subtasks").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	msg, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(msg.String())

	//runtask - 選擇事項
	jobWorker := client.NewJobWorker().JobType("io.camunda.zeebe:userTask").Handler(handleJobSubtasks).Open()

	<-readyClose
	jobWorker.Close()
	jobWorker.AwaitClose()

}

func handleJobSubtasks(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		failjob(client, job)
		return
	}

	variables["name"] = "pin"
	variables["time"] = "2023-03-25 12:57:03"

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
