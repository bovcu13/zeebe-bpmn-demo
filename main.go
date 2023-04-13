package main

import (
	"fmt"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/olekukonko/tablewriter"
	"github.com/zb-user/zb-example/demojob"
	"os"
)

func main() {
	variables := map[string]interface{}{
		"item": "Task Assignment",
	}

	jobHeaders := demojob.FindUserTask()
	PrintJobInfoTable(jobHeaders)

	fmt.Print("請輸入要審核任務的Id:")
	var taskId int64
	fmt.Scan(&taskId)
	demojob.CompleteJob(taskId, variables)
}

func PrintJobInfoTable(jobHeaders []entities.Job) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Job Key", "Process Instance Key", "Element ID"})

	for _, job := range jobHeaders {
		table.Append([]string{
			fmt.Sprintf("%d", job.GetKey()),
			fmt.Sprintf("%d", job.GetProcessInstanceKey()),
			job.GetElementId()})
	}
	table.Render()
}
