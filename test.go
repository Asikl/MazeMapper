package main

import (
	"bufio"
	"fmt"
	"hello/pkg/Cache"
	"os"
	"sync"
	"time"

	"hello/pkg/DrawGraph"
	"hello/pkg/Graph"
	"hello/pkg/Resolver"

	"github.com/miekg/dns"
)

func Total(domain string, sy *sync.WaitGroup) {
	fmt.Println("Hlol")
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

	fmt.Println("[                                  ]")
	gg.Dump()
	gg.DumpGraphMap()
	fmt.Println("----")
	DrawGraph.Visual(domain, &gg)
	//DrawGraph.Visual1(domain, &gg)
	//DrawGraph.Test()
	fmt.Println("+++")
	//Graph.DumpGraphReverse()
	//DrawGraph.Visual1()
	// DrawPicture.DrawGraph()
	// fmt.Println("没有错误查询数量：", Cache.GetERROR0())
	// fmt.Println("格式错误：", Cache.GetERROR1())
	// fmt.Println("Server failer：", Cache.GetERROR2())
	// fmt.Println("Name Error：", Cache.GetERROR3())
	// fmt.Println("Not implemented：", Cache.GetERROR4())
	// fmt.Println("Refused：", Cache.GetERROR5())
	// fmt.Println("Time out：", Cache.GetERROR6())
	// fmt.Println("数据包中没有NS记录：", Cache.GetERROR7())
	// fmt.Println("图中边的个数：", Cache.GetEdge())
	// fmt.Println("图中节点个数：", Graph.GetNum())
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
		// 获取扫描到的一行文本
		line := scanner.Text()

		// 打印读取到的一行文本
		domain = append(domain, line)
	}
	var wg sync.WaitGroup
	fmt.Println(len(domain), " ", domain)
	wg.Add(len(domain))
	for _, value := range domain {
		fmt.Println("协程", value)
		go Total(value, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start) // 计算经过的时间
	fmt.Println("程序运行时间：", elapsed)
}
