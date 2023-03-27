package demojob

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

// 建立本地端
func NewZeebeClient() (zbc.Client, error) {
	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: "localhost:26500"}
	client, err := zbc.NewClient(&config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// 部署bpmn
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

// 開啟一個流程實例
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

// 使用ProcessInstanceKey尋找jobKey
func FindJobKeysByProcessInstanceKey(processInstanceKey int64) ([]int64, error) {
	client, err := NewZeebeClient()
	if err != nil {
		return nil, err
	}

	var taskJobKeys []int64
	response, err := client.NewActivateJobsCommand().
		JobType("io.camunda.zeebe:userTask").
		MaxJobsToActivate(100).
		WorkerName("my-worker").
		Timeout(time.Minute).
		Send(context.Background())
	if err != nil {
		return nil, err
	}
	for _, job := range response {
		headers, ok := job.GetCustomHeadersAsMap()
		if ok {
			if myHeaders, ok := headers["ProcessInstanceKey"].(string); ok {
				isEqual := myHeaders == processInstanceKey
				if isEqual {
					taskJobKeys = append(taskJobKeys, job.GetKey())
				}
			}
		}
	}

	return taskJobKeys, nil
}

// 執行工作
func CompleteJob(jobKey int64, key []string, value interface{}) {
	client, err := NewZeebeClient()
	if err != nil {
		panic(err)
	}

	variables := make(map[string]interface{})
	// 將字串陣列轉換為單一的字串
	var keyStr string
	for _, k := range key {
		keyStr += k
	}
	variables[keyStr] = value

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
