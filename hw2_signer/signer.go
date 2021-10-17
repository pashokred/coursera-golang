package main

import (
	"fmt"
	"log"
	"strconv"
)

// сюда писать код

//func startWorker(worker job, wg* sync.WaitGroup, in, out chan interface{})  {
//	wg.Add(1)
//	defer wg.Done()
//	go worker(in, out)
//}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{}, MaxInputDataLen)
	out := make(chan interface{}, MaxInputDataLen)

	for _, job := range jobs {
		go job(in, out)
	}

LOOP:
	for {
		select {
		case val := <-out:
			fmt.Printf("get value %d in out\n", val.(uint32))
			in <- val
		default:
			break LOOP
		}
	}
	//select {
	//case data = <-out:
	//	val, ok := data.(uint32)
	//	if !ok {
	//		log.Fatal("cannot convert data to string")
	//	}
	//	fmt.Println("out val", val)
	//case in <- data:
	//	val, ok := data.(uint32)
	//	if !ok {
	//		log.Fatal("cannot convert data to string")
	//	}
	//	fmt.Println("put val to in:", val)
	//}

}

func SingleHash(in, out chan interface{}) {
	for val := range in {
		data, ok := val.(string)
		if !ok {
			log.Fatal("cant convert result data to string in SingleHash")
		}
		md5 := DataSignerMd5(data)
		crc32FromMd5 := DataSignerCrc32(md5)
		crc32 := DataSignerCrc32(data)
		result := crc32 + "~" + crc32FromMd5

		fmt.Printf(data, "SingleHash data", data)
		fmt.Printf(data, "SingleHash md5(data)", md5)
		fmt.Printf(data, "SingleHash crc32(md5(data))", crc32FromMd5)
		fmt.Printf(data, "SingleHash crc32(data)", crc32)
		fmt.Printf(data, "SingleHash result", result)
		out <- result
	}
}

func MultiHash(in, out chan interface{}) {
	for val := range in {
		data, ok := val.(string)
		if !ok {
			log.Fatal("cant convert result data to string in MultiHash")
		}
		var result string
		for i := 0; i < 6; i++ {
			th := DataSignerCrc32(strconv.Itoa(i) + data)
			result += th
			fmt.Printf(data, "MultiHash: crc32(th+step1))", i, th)
		}
		fmt.Printf(data, "MultiHash result:", result)
		out <- result
	}
}

func CombineResults(in, out chan interface{}) {

}
