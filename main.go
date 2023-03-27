package main

import (
	"fmt"

	"github.com/zb-user/zb-example/demojob"
)

func main() {
	fmt.Println(demojob.FindJobKeysByProcessInstanceKey("2251799813721794"))
}
