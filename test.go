package main

import (
	"fmt"
	"hello/pkg/Cache"
	"sync"

	"hello/pkg/DrawGraph"
	"hello/pkg/Graph"
	"hello/pkg/Resolver"

	"github.com/miekg/dns"
)

func Total(domain string, sy *sync.WaitGroup) {
	fmt.Println("Hlol")
	var dig Resolver.Dig
	var gg Graph.GraphStruct
	Cache.InitERROR()

	Graph.Init(&gg, domain)
	//rsps, err := dig.Trace(domain) //dig +trace

	_, err := dig.Trace(domain, dns.TypeA, &gg) //dig +trace

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
	//DrawGraph.Visual1(domain, gg)
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
	//var domain string
	// var dig Resolver.Dig
	// var gg Graph.GraphStruct

	// Cache.InitERROR()
	//Draw.DrawGraph()
	//fmt.Println("please input your domain：")
	//fmt.Scanf("%s", &domain)
	var wg sync.WaitGroup

	wg.Add(3)
	var domain = [...]string{"www.baidu.com", "eamis.nankai.edu.cn", "69guy.net"}

	for _, value := range domain {
		fmt.Println("协程", value)
		go Total(value, &wg)
	}
	wg.Wait()

}
