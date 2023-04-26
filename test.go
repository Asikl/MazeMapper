package main

import (
	"bufio"
	"fmt"
	"hello/pkg/Cache"
	"hello/pkg/DrawGraph"
	"hello/pkg/Graph"
	"hello/pkg/Resolver"
	"os"
	"sync"
	"time"

	"github.com/panjf2000/ants"

	"github.com/miekg/dns"
)

func Total(domain string, sy *sync.WaitGroup) {
	//fmt.Println("Hlol")
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
	//DrawGraph.Visual1()
	// DrawPicture.DrawGraph()
	// fmt.Println("没有错误查询数量：", cc.GetERROR0())
	// fmt.Println("格式错误：", cc.GetERROR1())
	// fmt.Println("Server failer：", cc.GetERROR2())
	// fmt.Println("Name Error：", cc.GetERROR3())
	// fmt.Println("Not implemented：", cc.GetERROR4())
	// fmt.Println("Refused：", cc.GetERROR5())
	// fmt.Println("Time out：", cc.GetERROR6())
	// fmt.Println("数据包中没有NS记录：", cc.GetERROR7())
	// fmt.Println("图中边的个数：", cc.GetEdge())
	//fmt.Println("图中节点个数：", cc.GetNum())
	sy.Done()
}

func main() {
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
		//ants.Submit(Total(value, &wg))
		//pool.Submit(Total(value, &wg))

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
