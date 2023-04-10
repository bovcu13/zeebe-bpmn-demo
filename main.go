package main

import "github.com/zb-user/zb-example/demojob"

func main() {
	variables := map[string]interface{}{
		"status": true,
	}
	//demojob.StartInstance("fill-out-timesheets", variables)

	jobInfo := demojob.GetJobInfoByProcessInstanceKey(2251799813722314)
	if len(jobInfo) == 0 {
		panic("找不到對應的 jobKey")
	}

	jobKey := jobInfo[0].Key // 假設 GetJobInfoByProcessInstanceKey 只返回一个 jobInfo

	demojob.CompleteJob(jobKey, variables)

}
