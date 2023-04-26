package main

import (
	"bufio"
	"fmt"
	"hello/pkg/Cache"
	"hello/pkg/DrawGraph"
	"hello/pkg/Graph"
	"hello/pkg/Resolver"
	"log"
	"os"
	"sync"
	"time"

	"github.com/panjf2000/ants"

	"github.com/miekg/dns"
)

func Total(domain string, sy *sync.WaitGroup) {
	//fmt.Println("Hlol")
	var rwm sync.RWMutex
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

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	//Cache.Dump()

	//fmt.Println("[                                  ]")
	//gg.Dump()
	//gg.DumpGraphMap()
	fmt.Println("----")
	DrawGraph.Visual(domain, &gg)
	//DrawGraph.Visual1(domain, &gg)
	//DrawGraph.Test()
	fmt.Println("+++")
	//Graph.DumpGraphReverse()
	DrawGraph.Visual1(domain, &gg)
	//DrawPicture.DrawGraph()
	str := fmt.Sprintf("\"%s\" \n", domain)
	//加锁
	rwm.Lock()
	if cc.GetERROR1() != 0 {
		//写入文件

		// 打开文件，如果文件不存在则创建
		file, err := os.OpenFile("Corrupt.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 写入新的内容到文件末尾
		file.WriteString(str)
	}
	if cc.GetERROR1() != 0 {
		//写入文件
		// 打开文件，如果文件不存在则创建
		file, err := os.OpenFile("./Result/Serverfailure.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// 写入新的内容到文件末尾
		file.WriteString(str)
	}
	if cc.GetERROR3() != 0 {
		// 打开文件，如果文件不存在则创建
		file, err := os.OpenFile("./Result/NXDOMAIN.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// 写入新的内容到文件末尾
		file.WriteString(str)
	}
	if cc.GetERROR4() != 0 {
		//写入文件
		file, err := os.OpenFile("./Result/NotImplmented.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// 写入新的内容到文件末尾
		file.WriteString(str)
	}
	if cc.GetERROR5() != 0 {
		//写入文件
		file, err := os.OpenFile("./Result/Refused.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// 写入新的内容到文件末尾
		file.WriteString(str)
	}
	if cc.GetERROR6() != 0 {
		//写入文件

		//写入文件
		file, err := os.OpenFile("./Result/Timeout.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// 写入新的内容到文件末尾
		file.WriteString(str)
		//os.WriteFile("Timeout.txt", []byte(str), 0664)
	}
	if cc.GetERROR7() != 0 {
		//写入文件
		file, err := os.OpenFile("./Result/NoNsRecord.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// 写入新的内容到文件末尾
		file.WriteString(str)
	}
	//解锁
	rwm.Unlock()
	// fmt.Println("没有错误查询数量：", cc.GetERROR0())
	// fmt.Println("Corrept：", cc.GetERROR1())
	// fmt.Println("Server failer：", cc.GetERROR2())
	// fmt.Println("Name Error：", cc.GetERROR3())
	// fmt.Println("Not implemented：", cc.GetERROR4())
	// fmt.Println("Refused：", cc.GetERROR5())
	// fmt.Println("Time out：", cc.GetERROR6())
	// fmt.Println("NoNsRecord：", cc.GetERROR7())
	// fmt.Println("图中边的个数：", cc.GetEdge())
	// fmt.Println("图中节点个数：", cc.GetNum())
	sy.Done()
}

func main() {
	//生成存储结果的文件
	CreateFile()
	start := time.Now() // 记录开始时间
	// 打开文本文件,读取txt
	var domain = make([]string, 0)
	file, err := os.Open("example.txt")
	if err != nil {
		fmt.Println("文件读取错误", err)
	}
	defer file.Close()
	// 使用bufio.NewReader创建缓冲读取器
	scanner := bufio.NewScanner(file)
	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()
		domain = append(domain, line)
	}
	wg := sync.WaitGroup{}

	//申请一个协程池对象
	pool, _ := ants.NewPool(100)
	//关闭协程池
	defer pool.Release()
	//fmt.Println("POOL", pool.Running())
	//var wg sync.WaitGroup
	fmt.Println(len(domain), " ", domain)
	wg.Add(len(domain))
	for _, value := range domain {
		fmt.Println("协程", value)
		pool.Submit(func() {
			// 调用 processData 函数处理数据
			Total(value, &wg)
		})
		fmt.Println("POOLOOOOOOOOOOOOOOOOOOOOOO", pool.Running())
	}
	wg.Wait()

	elapsed := time.Since(start) // 计算经过的时间
	fmt.Println("程序运行时间：", elapsed)
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
