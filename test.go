package main

import (
	"bufio"
	"fmt"
	"hello/pkg/Cache"
	"hello/pkg/DrawGraph"

	//"hello/pkg/DrawGraph"
	"hello/pkg/Graph"
	"hello/pkg/Resolver"
	"log"
	"os"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/panjf2000/ants"
)

func Total(domain string) (result map[string]string) {
	//fmt.Println("Hlol")
	// var rwm sync.RWMutex

	var dig Resolver.Dig
	var gg Graph.GraphStruct
	var cc Cache.CacheStruct

	Cache.Init(domain, &cc)
	Graph.Init(&gg, domain)
	//rsps, err := dig.Trace(domain) //dig +trace
	_, err := dig.Trace(domain, dns.TypeA, &gg, &cc) //dig +trace

	//tempmap, err := dig.TraceIP(domain)
	//_, err := dig.RRquery(domain)
	if err != nil {
		fmt.Println(err)
		return
	}
	//os.WriteFile()

	//fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	//fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	//Cache.Dump()

	//fmt.Println("----")

	DrawGraph.GetDot(domain, &gg) //生成图，图中节点为数字

	//DrawGraph.Visual1(domain, &gg)  //生成图，节点为真实数据<domain,qtype,Ip>
	//DrawGraph.Test()
	fmt.Println("+++")
	//gg.DumpGraphReverse()
	//DrawPicture.DrawGraph()
	//str := fmt.Sprintf("\"%s\" \n", domain)

	if cc.GetERROR1() != 0 {
		//写入文件
		gg.Result["./Result/Corrupt.txt"] = gg.Domain + "\n"
	}
	if cc.GetERROR2() != 0 {
		//写入文件
		// 打开文件，如果文件不存在则创建
		gg.Result["./Result/Serverfailure.txt"] = gg.Domain + "\n"
	}
	if cc.GetERROR3() != 0 {
		// 打开文件，如果文件不存在则创建
		gg.Result["./Result/NXDOMAIN.txt"] = gg.Domain + "\n"

	}
	if cc.GetERROR4() != 0 {
		//写入文件
		gg.Result["./Result/NotImplmented.txt"] = gg.Domain + "\n"

	}
	if cc.GetERROR5() != 0 {
		//写入文件
		gg.Result["./Result/Refused.txt"] = gg.Domain + "\n"
	}
	if cc.GetERROR6() != 0 {
		//写入文件

		//写入文件
		gg.Result["./Result/Timeout.txt"] = gg.Domain + "\n"

	}
	if cc.GetERROR7() != 0 {
		//写入文件
		gg.Result["./Result/NoNsRecord.txt"] = gg.Domain + "\n"

	}

	if cc.GetERROR8() != 0 {
		//写入文件
		gg.Result["./Result/Hijack.txt"] = gg.Domain + "\n"
	}
	if cc.GetERROR9() != 0 {
		//写入文件
		gg.Result["./Result/IPerror.txt"] = gg.Domain + "\n"
	}

	fmt.Println("       ========            ", domain, gg.GetNum())

	if gg.GetNum() > 200 {
		gg.Result["./Result/NodeOver200.txt"] = gg.Domain + "\n"
	}

	return gg.Result
}

func main() {
	//生成存储结果的文件
	CreateFile()
	start := time.Now() // 记录开始时间
	// 打开文本文件,读取txt

	wg := sync.WaitGroup{}
	m := sync.Mutex{}
	//var mu sync.Mutex
	//申请一个协程池对象
	pool, _ := ants.NewPool(10000)
	//关闭协程池
	defer pool.Release()
	file, err := os.Open("top1_mpopular.txt")
	if err != nil {
		fmt.Println("文件读取错误", err)
	}
	defer file.Close()
	// 使用bufio.NewReader创建缓冲读取器
	scanner := bufio.NewScanner(file)
	// 逐行读取文件内容

	var num = 1
	for scanner.Scan() {
		line := scanner.Text()
		//domain = append(domain, line)
		num++

		fmt.Println("协程", line, "           num:", num)
		wg.Add(1)
		pool.Submit(func() {
			// 调用 processData 函数处理数据
			defer wg.Done()
			map1 := Total(line)
			m.Lock()
			WriteResult(map1)
			m.Unlock()
		})
		//fmt.Println("POOLOOOOOOOOOOOOOOOOOOOOOO", pool.Running())
	}
	wg.Wait()

	elapsed := time.Since(start) // 计算经过的时间
	fmt.Println("程序运行时间：", elapsed)
}

func WriteResult(str map[string]string) {

	for index, value := range str {
		file, err := os.OpenFile(index, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// 写入新的内容到文件末尾
		file.WriteString(value)

	}
}

func CreateFile() {
	_, Timeouterr := os.Create("./Result/Timeout.txt")
	if Timeouterr != nil {
		log.Fatal(Timeouterr)
		fmt.Println(Timeouterr)
	}
	_, Refusederr := os.Create("./Result/Refused.txt")
	if Refusederr != nil {
		log.Fatal(Refusederr)
		fmt.Println(Refusederr)
	}
	_, NXDOMAINerr := os.Create("./Result/NXDOMAIN.txt")
	if NXDOMAINerr != nil {
		log.Fatal(NXDOMAINerr)
		fmt.Println(NXDOMAINerr)
	}
	_, NoNsRecorderr := os.Create("./Result/NoNsRecord.txt")
	if NoNsRecorderr != nil {
		log.Fatal(NoNsRecorderr)
		fmt.Println(NoNsRecorderr)
	}
	_, Serverfailureerr := os.Create("./Result/Serverfailure.txt")
	if Serverfailureerr != nil {
		log.Fatal(Serverfailureerr)
		fmt.Println(Serverfailureerr)
	}
	_, NotImplmentederr := os.Create("./Result/NotImplmented.txt")
	if NotImplmentederr != nil {
		log.Fatal(NotImplmentederr)
		fmt.Println(NotImplmentederr)
	}
	_, Corrupterr := os.Create("./Result/Corrupt.txt")
	if Corrupterr != nil {
		log.Fatal(Corrupterr)
		fmt.Println(Corrupterr)
	}
}
