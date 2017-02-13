package mapreduce

import (
	"hash/fnv"
	"log"
	"io/ioutil"
	//"encoding/json"
	"os"
	"fmt"
	"encoding/json"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	// TODO:
	// You will need to write this function.
	// You can find the filename for this map task's input to reduce task number
	// r using reduceName(jobName, mapTaskNumber, r). The ihash function (given
	// below doMap) should be used to decide which file a given key belongs into.
	//
	// The intermediate output of a map task is stored in the file
	// system as multiple files whose name indicates which map task produced
	// them, as well as which reduce task they are for.
	// map任务的中间输出存储在文件中，文件名标识是哪个map任务产生该文件以及应该被哪个reduce任务处理
	// 可以使用json格式作为持续化格式
	// Coming up with a scheme for how to store the key/value pairs on disk can be tricky,
	// especially when taking into account that both keys and values could
	// contain newlines, quotes, and any other character you can think of.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!

	dat, err :=  ioutil.ReadFile(inFile)
	var r []([]KeyValue) = make([]([]KeyValue), nReduce)
	if err != nil{
		log.Fatal(err)
		return
	}
	dat_str := string(dat)
	fmt.Println("开始执行Map任务")
	inter := mapF(inFile, dat_str) // mapF 生成一个[]keyValue
	for _, v := range inter{
		slot := int(ihash(v.Key)) % nReduce
		r[slot] = append(r[slot],v)
	}

	for i := 0; i<nReduce;i++  {
		filename := reduceName(jobName, mapTaskNumber, i)
		outputFile, err := os.OpenFile(filename , os.O_WRONLY|os.O_CREATE, 0666)
		defer outputFile.Close()
		if err != nil {
			log.Fatal("写入失败")
		}
		enc:=json.NewEncoder(outputFile)
		for _, kv := range(r[i]){
			err = enc.Encode(&kv)
		}
	}
}

func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
