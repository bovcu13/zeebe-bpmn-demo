package demojob

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
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
func DeployBPMNFile(fileName string) (response *pb.DeployResourceResponse) {
	ctx := context.Background()
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	// deploy BPMN file
	response, err = client.NewDeployResourceCommand().AddResourceFile(fileName).Send(ctx)
	if err != nil {
		log.Fatalf("Failed to deploy BPMN file: %s", err)
	}

	// return deployment key
	return response
}

// StartInstance
// @Description: 開啟一個流程實例
func StartInstance(processId string) (*pb.CreateProcessInstanceResponse, error) {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	request := client.NewCreateInstanceCommand().BPMNProcessId(processId).LatestVersion()

	ctx := context.Background()
	response, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(response.String())

	return response, nil
}

// CancelProcessInstance
// @Description: 取消一個流程實例
func CancelProcessInstance(processInstanceKey int64) (string, error) {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	request := client.NewCancelInstanceCommand().ProcessInstanceKey(processInstanceKey)
	ctx := context.Background()
	response, err := request.Send(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to send cancel instance command: %w", err)
	}
	fmt.Printf("Cancelled process instance with key %d. Response: %s\n", processInstanceKey, response.String())
	return fmt.Sprintf("Successfully cancelled process instance with key %d", processInstanceKey), nil
}

// FindUserTask
// @Description: 查詢 userTask 所有待完成任務詳細資訊
func FindUserTask() []entities.Job {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	request := client.NewActivateJobsCommand().
		JobType("io.camunda.zeebe:userTask").
		MaxJobsToActivate(100).
		Timeout(time.Minute).
		WorkerName("my-worker")

	jobHeaders, err := request.Send(context.Background())
	if err != nil {
		panic(err)
	}

	return jobHeaders
}

// GetJobInfoByProcessInstanceKey
// @Description: 查詢 userTask 中指定的 ProcessInstanceKey 所有待完成任務詳細資訊
func GetJobInfoByProcessInstanceKey(processInstanceKey int64) []entities.Job {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	request := client.NewActivateJobsCommand().
		JobType("io.camunda.zeebe:userTask").
		MaxJobsToActivate(100).
		Timeout(time.Minute).
		WorkerName("my-worker")

	jobHeaders, err := request.Send(context.Background())
	if err != nil {
		panic(err)
	}

	return jobHeaders
}

// GetJobInfoByProcessDefinitionKey
// @Description: 查詢 userTask 中指定的 ProcessDefinitionKey 所有待完成任務詳細資訊
func GetJobInfoByProcessDefinitionKey(processDefinitionKey int64) []entities.Job {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	request := client.NewActivateJobsCommand().
		JobType("io.camunda.zeebe:userTask").
		MaxJobsToActivate(100).
		Timeout(time.Minute).
		WorkerName("my-worker")

	jobHeaders, err := request.Send(context.Background())
	if err != nil {
		panic(err)
	}

	return jobHeaders
}

// GetJobInfoByProcessId
// @Description: 查詢 userTask 中指定的 processId 所有待完成任務詳細資訊
func GetJobInfoByProcessId(processId string) []entities.Job {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	request := client.NewActivateJobsCommand().
		JobType("io.camunda.zeebe:userTask").
		MaxJobsToActivate(100).
		Timeout(time.Minute).
		WorkerName("my-worker")

	jobHeaders, err := request.Send(context.Background())
	if err != nil {
		panic(err)
	}

	return jobHeaders
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
