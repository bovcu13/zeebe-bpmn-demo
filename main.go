package main

import "github.com/zb-user/zb-example/demojob"

func main() {
	//demojob.GetJobInfoByProcessInstanceKey(2251799813721897)
	//demojob.GetJobInfoByProcessDefinitionKey(2251799813721383)
	//demojob.GetJobInfoByProcessId("choices-subtasks")
	//demojob.FindUserTask()
	demojob.CancelProcessInstance(2251799813721897)
	//demojob.StartInstance("choices-subtasks", nil)

	//variables := map[string]interface{}{
	//	"meal": "Chicken",
	//}
	//demojob.CompleteJob(2251799813721870, variables)

}
