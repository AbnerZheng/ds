package mapreduce

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

// Shutdown is an RPC method that shuts down the Master's RPC server.
func (mr *Master) Shutdown(_, _ *struct{}) error {
	debug("Shutdown: registration server\n")
	close(mr.shutdown)
	mr.l.Close() // causes the Accept to fail
	return nil
}

// startRPCServer starts the Master's RPC server. It continues accepting RPC
// calls (Register in particular) for as long as the worker is alive.
// startRPCServer 启动一个Master的RPC服务器. 它持续接受RPC调用直到worker死去
func (mr *Master) startRPCServer() {
	rpcs := rpc.NewServer()
	rpcs.Register(mr)
	os.Remove(mr.address) // only needed for "unix"
	l, e := net.Listen("unix", mr.address)
	if e != nil {
		log.Fatal("RegstrationServer", mr.address, " error: ", e)
	}
	mr.l = l

	// now that we are listening on the master address, can fork off
	// accepting connections to another thread.
	// 既然我们监听在master地址， 可以fork 一个新的线程来接受连接
	go func() {
	loop:
		for {
			select {
			case <-mr.shutdown:
				break loop
			default:
			}
			conn, err := mr.l.Accept()
			if err == nil {
				go func() {
					rpcs.ServeConn(conn)
					conn.Close()
				}()
			} else {
				debug("RegistrationServer: accept error", err)
				break
			}
		}
		debug("RegistrationServer: done\n")
	}()
}

// stopRPCServer stops the master RPC server.
// This must be done through an RPC to avoid race conditions between the RPC
// server thread and the current thread.
// 需要使用RPC来避免竞争条件
// stopRPCServer 停止master RPC server
func (mr *Master) stopRPCServer() {
	var reply ShutdownReply
	ok := call(mr.address, "Master.Shutdown", new(struct{}), &reply)
	if ok == false {
		fmt.Printf("Cleanup: RPC %s error\n", mr.address)
	}
	debug("cleanupRegistration: done\n")
}
