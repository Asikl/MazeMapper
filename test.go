package main

import (
	"fmt"
	"hello/pkg/Cache"

	"hello/pkg/DrawGraph"
	"hello/pkg/Graph"
	"hello/pkg/Resolver"

	"github.com/miekg/dns"
)

func main() {

	var domain string
	var dig Resolver.Dig
	Cache.InitERROR()
	//Draw.DrawGraph()
	fmt.Println("please input your domain：")

	fmt.Scanf("%s", &domain)
	//rsps, err := dig.Trace(domain) //dig +trace
	_, err := dig.Trace(domain, dns.TypeA) //dig +trace

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
	//Graph.Dump()
	Graph.DumpGraphMap()
	fmt.Println("----")
	DrawGraph.Visual()
	//DrawGraph.Test()
	fmt.Println("+++")
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
}
