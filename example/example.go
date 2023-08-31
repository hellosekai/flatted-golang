/**
 * @Description:
 * @FilePath: /flatted-golang/example/example.go
 * @Author: liyibing liyibing@lixiang.com
 * @Date: 2023-08-29 11:35:35
 */
package main

import (
	"fmt"

	"github.com/hellosekai/flatted"
)

type exampleStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	test := exampleStruct{
		ID:   1,
		Name: "test",
	}
	// fmt.Println(test)
	str, _ := flatted.FlattedFromStruct(test)
	// fmt.Println(str)
	ans := exampleStruct{}
	flatted.UnFlattedToStruct(str, &ans)
	// fmt.Println(ans)

	ttt := "{\"startData\":{},\"resultData\":{\"runData\":{}},\"executionData\":{\"contextData\":{},\"nodeExecutionStack\":[{\"node\":{\"parameters\":{\"list\":\"listn8n\",\"options\":{}},\"id\":\"09536973-b387-427e-ac14-51975a25391f\",\"name\":\"Redis Sentinel Trigger Limit\",\"type\":\"CUSTOM.redisSentinelTriggerLimit\",\"typeVersion\":1,\"position\":[860,580],\"credentials\":{\"redisSentinel\":{\"id\":\"4000087\",\"name\":\"Redis Sentinel account 2\"}}},\"data\":{\"main\":[[{\"json\":{\"list\":\"listn8n\",\"value\":\"listn8n01\"}}]]},\"source\":null}],\"waitingExecution\":{},\"waitingExecutionSource\":{}}}"
	tttAns := flatted.Flatted(ttt)
	fmt.Println(tttAns)
}
