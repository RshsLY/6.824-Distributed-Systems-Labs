package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"

type Coordinator struct {
	// Your definitions here.

}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	if args.status == "mapped" {
		delete(toMap, args.fileName)
		return nil
	}
	if args.status == "reduced" {
		delete(toReduced, args.fileName)
		return nil
	}
	for s := range toMap {
		if toMap[s] != 0 {
			continue
		}
		toMap[s] = 1
		reply.status = "toMap"
		reply.fileName = s
		return nil
	}
	for s := range toReduced {
		reply.status = "toReduced"
		reply.fileName = s
		return nil
	}
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	// Your code here.
	return len(toReduced) == 0
}

var toMap = make(map[string]int)
var toReduced = make(map[string]int)

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	// Your code here.
	for i := range files {
		toMap[files[i]] = 0
		toReduced[files[i]] = 0
	}
	c.server()
	return &c
}
