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
		delete(toMap, args.FileName)
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
	// Your code here.
	return len(toMap) == 0
}

var toMap = make(map[string]int)

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	// Your code here.
	for i := range files {
		toMap[files[i]] = 0
	}
	c.server()
	return &c
}
