package demojob

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

type info struct {
	JobType                  string `json:"jobType"`
	ProcessId                string `json:"processId"`
	Key                      int64  `json:"key"`
	ProcessInstanceKey       int64  `json:"processInstanceKey"`
	ProcessDefinitionVersion int32  `json:"processDefinitionVersion"`
	ProcessDefinitionKey     int64  `json:"processDefinitionKey"`
	Variables                string `json:"variables"`
	ElementId                string `json:"elementId"`
	ElementInstanceKey       int64  `json:"elementInstanceKey"`
	CustomHeaders            string `json:"customHeaders"`
	Worker                   string `json:"worker"`
	Retries                  int32  `json:"retries"`
	Deadline                 int64  `json:"deadline"`
}

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
		return "", err
	}
	fmt.Println(response.String())

	return fmt.Sprintf("Successfully cancelled process instance with key %d", processInstanceKey), nil
}

// FindUserTask
// @Description: 查詢 userTask 所有待完成任務詳細資訊
func FindUserTask() []info {
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

	var result []info
	for _, header := range jobHeaders {
		info := info{
			JobType:                  header.GetType(),
			ProcessId:                header.GetBpmnProcessId(),
			Key:                      header.GetKey(),
			ProcessInstanceKey:       header.GetProcessInstanceKey(),
			ProcessDefinitionVersion: header.GetProcessDefinitionVersion(),
			ProcessDefinitionKey:     header.GetProcessDefinitionKey(),
			Variables:                header.GetVariables(),
			ElementId:                header.GetElementId(),
			ElementInstanceKey:       header.GetElementInstanceKey(),
			CustomHeaders:            header.GetCustomHeaders(),
			Worker:                   header.GetWorker(),
			Retries:                  header.GetRetries(),
			Deadline:                 header.GetDeadline(),
		}
		result = append(result, info)
	}

	Marshal, _ := json.Marshal(result)
	fmt.Println(string(Marshal))
	return result
}

// GetJobInfoByProcessInstanceKey
// @Description: 查詢 userTask 中指定的 ProcessInstanceKey 所有待完成任務詳細資訊
func GetJobInfoByProcessInstanceKey(processInstanceKey int64) []info {
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

	var result []info
	for _, header := range jobHeaders {
		if processInstanceKey == header.GetProcessInstanceKey() {
			info := info{
				JobType:                  header.GetType(),
				ProcessId:                header.GetBpmnProcessId(),
				Key:                      header.GetKey(),
				ProcessInstanceKey:       header.GetProcessInstanceKey(),
				ProcessDefinitionVersion: header.GetProcessDefinitionVersion(),
				ProcessDefinitionKey:     header.GetProcessDefinitionKey(),
				Variables:                header.GetVariables(),
				ElementId:                header.GetElementId(),
				ElementInstanceKey:       header.GetElementInstanceKey(),
				CustomHeaders:            header.GetCustomHeaders(),
				Worker:                   header.GetWorker(),
				Retries:                  header.GetRetries(),
				Deadline:                 header.GetDeadline(),
			}

			result = append(result, info)
		}
	}
	Marshal, _ := json.Marshal(result)
	fmt.Println(string(Marshal))
	return result
}

// GetJobInfoByProcessDefinitionKey
// @Description: 查詢 userTask 中指定的 ProcessDefinitionKey 所有待完成任務詳細資訊
func GetJobInfoByProcessDefinitionKey(processDefinitionKey int64) []info {
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

	var result []info
	for _, header := range jobHeaders {
		if processDefinitionKey == header.GetProcessDefinitionKey() {
			info := info{
				JobType:                  header.GetType(),
				ProcessId:                header.GetBpmnProcessId(),
				Key:                      header.GetKey(),
				ProcessInstanceKey:       header.GetProcessInstanceKey(),
				ProcessDefinitionVersion: header.GetProcessDefinitionVersion(),
				ProcessDefinitionKey:     header.GetProcessDefinitionKey(),
				Variables:                header.GetVariables(),
				ElementId:                header.GetElementId(),
				ElementInstanceKey:       header.GetElementInstanceKey(),
				CustomHeaders:            header.GetCustomHeaders(),
				Worker:                   header.GetWorker(),
				Retries:                  header.GetRetries(),
				Deadline:                 header.GetDeadline(),
			}

			result = append(result, info)
		}
	}
	Marshal, _ := json.Marshal(result)
	fmt.Println(string(Marshal))
	return result
}

// GetJobInfoByProcessId
// @Description: 查詢 userTask 中指定的 processId 所有待完成任務詳細資訊
func GetJobInfoByProcessId(processId string) []info {
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

	var result []info
	for _, header := range jobHeaders {
		if processId == header.GetBpmnProcessId() {
			info := info{
				JobType:                  header.GetType(),
				ProcessId:                header.GetBpmnProcessId(),
				Key:                      header.GetKey(),
				ProcessInstanceKey:       header.GetProcessInstanceKey(),
				ProcessDefinitionVersion: header.GetProcessDefinitionVersion(),
				ProcessDefinitionKey:     header.GetProcessDefinitionKey(),
				Variables:                header.GetVariables(),
				ElementId:                header.GetElementId(),
				ElementInstanceKey:       header.GetElementInstanceKey(),
				CustomHeaders:            header.GetCustomHeaders(),
				Worker:                   header.GetWorker(),
				Retries:                  header.GetRetries(),
				Deadline:                 header.GetDeadline(),
			}

			result = append(result, info)
		}
	}
	Marshal, _ := json.Marshal(result)
	fmt.Println(string(Marshal))
	return result
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
