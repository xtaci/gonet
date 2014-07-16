package helper

import (
	"log"
	"runtime"
)

func PrintPanicStack() {
	if x := recover(); x != nil {
		log.Printf("%v", x)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				log.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
	}
}
