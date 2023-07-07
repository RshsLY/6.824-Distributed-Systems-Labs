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
	if args.Status == "mapped" {
		delete(mapTask, args.FileName)
		return nil
	}
	if args.Status == "reduced" {
		delete(mapTask, args.FileName)
		return nil
	}
	if args.Status == "toMap" {
		for s := range toMap {
			//已经有进程在map了
			if toMap[s] != 0 {
				continue
			}
			toMap[s] = 1
			reply.Status = "toMap"
			reply.FileName = s
			return nil
		}
	}
	reply.Status = "allInDeal"
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

	return len(reduceTask) == 0
}

//0:not map  1:mapping
var mapTask = make(map[string]int)

//0:toReduce 1:reducing
var reduceTask = make([]int, 10)
var nReduce int = 0

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	// Your code here.
	for i := range files {
		mapTask[files[i]] = 0
	}
	reduceTask = make([]int, nReduce)
	for i := 0; i < nReduce; i++ {
		reduceTask[i] = 0
	}
	nReduce = nReduce
	c.server()
	return &c
}
