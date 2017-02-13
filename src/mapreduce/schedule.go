package mapreduce

import (
	"fmt"
	"sync"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	var wg sync.WaitGroup
	for i := 0; i < ntasks; i++ {
		wg.Add(1)
		var rpc_call func(file string, phase jobPhase, taskNumber int, other int)
		rpc_call = func(file string, phase jobPhase, taskNumber int, other int) { //开始发起rpc请求
			// 构造传参
			doTaskArgs := DoTaskArgs{JobName: mr.jobName, File: file, Phase: phase, TaskNumber: taskNumber, NumOtherPhase: other}
			for worker := range (mr.registerChannel) {
				ok := call(worker, "Worker.DoTask", doTaskArgs, new(struct{}))
				if ok {
					go func() {
						mr.registerChannel <- worker
					}()
					wg.Done()
					return
				} else {
					go rpc_call(file, phase, taskNumber, other)
					return
				}
			} // 拿到一个空闲的worker
		}
		go rpc_call(mr.files[i], phase, i, nios)
	}
	wg.Wait()

	fmt.Printf("Schedule: %v phase done\n", phase)
}
