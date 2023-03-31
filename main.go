package main

import "github.com/zb-user/zb-example/demojob"

func main() {
	//demojob.GetJobInfoByProcessInstanceKey(2251799813721865)
	//demojob.FindUserTask()
	//demojob.StartInstance("choices-subtasks", nil)
	variables := map[string]interface{}{
		"meal": "Chicken",
	}
	demojob.CompleteJob(2251799813721870, variables)
}
