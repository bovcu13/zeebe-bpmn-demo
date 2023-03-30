package demojob

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

// NewZeebeClient
// @Description: 建立本地端
func NewZeebeClient() (zbc.Client, error) {
	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: "localhost:26500"}
	client, err := zbc.NewClient(&config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// DeployBPMNFile
// @Description: 部署bpmn
func DeployBPMNFile(fileName string) string {
	ctx := context.Background()
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	// deploy BPMN file
	response, err := client.NewDeployResourceCommand().AddResourceFile(fileName).Send(ctx)
	if err != nil {
		log.Fatalf("Failed to deploy BPMN file: %s", err)
	}

	// return deployment key
	return response.String()
}

// StartInstance
// @Description: 開啟一個流程實例
func StartInstance(processId string, variables map[string]interface{}) (string, error) {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	request := client.NewCreateInstanceCommand().BPMNProcessId(processId).LatestVersion()

	if len(variables) > 0 {
		_, err = request.VariablesFromMap(variables)
		if err != nil {
			return "", err
		}
	}

	ctx := context.Background()
	response, err := request.Send(ctx)
	if err != nil {
		return "", err
	}

	fmt.Println(response.String())

	return response.String(), nil
}

// FindUserTask
// @Description: 查詢 userTask 待完成任務資訊
func FindUserTask() {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	request := client.NewActivateJobsCommand().
		JobType("io.camunda.zeebe:userTask").
		MaxJobsToActivate(10).
		Timeout(time.Minute).
		WorkerName("my-worker")

	jobHeaders, err := request.Send(context.Background())
	if err != nil {
		panic(err)
	}

	for _, header := range jobHeaders {
		log.Printf("ProcessId: %s\n"+
			"key: %d\n"+
			"ProcessInstanceKey: %d\n"+
			"ProcessDefinitionKey: %d\n"+
			"variables: %v\n", header.GetBpmnProcessId(), header.GetKey(), header.GetProcessInstanceKey(), header.GetProcessDefinitionKey(), header.GetVariables())
	}
}

// CompleteJob
// @Description: 執行工作
func CompleteJob(jobKey int64, variables map[string]interface{}) {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

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
