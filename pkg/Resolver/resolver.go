package Resolver

import (
	"fmt"
	"hello/pkg/Cache"
	"hello/pkg/Graph"

	//"strings"
	"time"

	"github.com/miekg/dns"
)

// Dig dig
type Dig struct {
	Domain           string
	LocalAddr        string
	RemoteAddr       string
	BackupRemoteAddr string
	DialTimeout      time.Duration
	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	Protocol         string
	Retry            int
	NSgetIP          map[string]string
}

// 存解决NS记录得到的结果，go语言好像没有set，我们用map来模拟set
//var NSgetIP = make(map[string]string, 0)

// TraceResponse  dig +trace 响应
type TraceResponse struct {
	Server   string
	ServerIP string
	Msg      *dns.Msg
	//type1    TypeDNSKEY.Msg
}

// Resolver
func (d *Dig) Resolver(domain string, msgType uint16, distination string, gg *Graph.GraphStruct, cc *Cache.CacheStruct, GetIP *[]string) ([]string, error) {
	//var responses = make([]TraceResponse, 0)
	var servers = make([]string, 0)
	//cacheFIX := make(map[cachekey]cachevalue)
	var NsNotGlueIP = make([]string, 0)
	server := distination
	//死循环
	for {
		if err := d.SetDNS(server); err != nil {
			fmt.Println("IP本身有问题", err)
			nodenum := gg.NodeNum(domain, dns.TypeA, server)
			gg.SetRerverseflag(nodenum, Graph.IPerrorNode)
			gg.AddNode(nodenum, Graph.IPerrorNode)
			gg.Setflag(domain, msgType, distination, 1)
			return nil, nil
		}
		if cc.Has(domain, server, dns.TypeA) {
			//先得到cache里的内容
			value := cc.GetCache(domain, server, dns.TypeA)

			//这是NS不带IP的情况
			if value.Flag {
				return nil, nil
			}

			nodenum := gg.NodeNum(domain, dns.TypeA, server)
			//Cache.AddEdge(1)
			for _, tempserver := range value.IP {
				//Cache.AddEdge(1)
				tempnodenum := gg.NodeNum(domain, dns.TypeA, tempserver)
				gg.AddNode(nodenum, tempnodenum)
			}
			//处理指向节点
			return nil, nil
		} else {
			flag := gg.Getflag(domain, dns.TypeA, server)
			if flag != Graph.NoVisit {
				//fmt.Println("该节点已访问")
				return nil, nil
			}
			msg, err, num := d.GetMsg(msgType, domain) //GetMsg

			switch num {
			case 100: //正常节点
				gg.Setflag(domain, msgType, server, Graph.Common)
			case 61:
				gg.Setflag(domain, msgType, server, Graph.Timeout)
				nodenum := gg.NodeNum(domain, dns.TypeA, server)
				gg.SetRerverseflag(nodenum, Graph.TimeoutNode)
				gg.AddNode(nodenum, Graph.TimeoutNode)
				fmt.Println(domain, server)
			case 62:
				gg.Setflag(domain, msgType, server, Graph.Timeout)
				nodenum := gg.NodeNum(domain, dns.TypeA, server)
				gg.SetRerverseflag(nodenum, Graph.TimeoutNode)
				gg.AddNode(nodenum, Graph.TimeoutNode)
			case 63:
				gg.Setflag(domain, msgType, server, Graph.Timeout)
				nodenum := gg.NodeNum(domain, dns.TypeA, server)
				gg.SetRerverseflag(nodenum, Graph.TimeoutNode)
				gg.AddNode(nodenum, Graph.TimeoutNode)
				fmt.Println("TIMEOUT!!!")
				gg.SetRerverseflag(nodenum, Graph.TimeoutNode)
			case 9:
				gg.Setflag(domain, msgType, server, Graph.IDMisMatch)
				nodenum := gg.NodeNum(domain, dns.TypeA, server)
				gg.AddNode(nodenum, Graph.IDMisMatchNode)
				gg.SetRerverseflag(nodenum, Graph.IDMisMatch)

			}
			//Graph.Setflag(domain, msgType, server, 1) //节点标记为已访问
			if err != nil {
				cc.AddERROR6()
				//fmt.Println("出现错误  ", fmt.Errorf("%v", err))
				fmt.Println("出现错误TIMEOUT!!!!!!!!!!!", server, domain)
				// gg.Setflag(domain, msgType, server, Graph.Timeout)
				// nodenum := gg.NodeNum(domain, dns.TypeA, server)
				// gg.SetRerverseflag(nodenum, Graph.TimeoutNode)
				// gg.AddNode(nodenum, Graph.TimeoutNode)

				//Graph.Setflag(domain, msgType, server, 5) //节点标记为错误节点
				return nil, fmt.Errorf("%s:%v", server, err)
			}
			gg.Setflag(domain, msgType, server, 1) //节点标记为已访问

			switch msg.Rcode {
			case 0:
				//fmt.Println("NOERROR")
				cc.AddERROR0()
			case 1:
				fmt.Println("出现错误  格式错误")
				nodenum := gg.NodeNum(domain, dns.TypeA, server)
				gg.SetRerverseflag(nodenum, Graph.CorruptNode)
				gg.AddNode(nodenum, Graph.CorruptNode)
				//gg.Setflag(domain, msgType, server, 1)
				cc.AddERROR1()
				return nil, nil
			case 2:
				fmt.Println("出现错误  Server Failure")
				nodenum := gg.NodeNum(domain, dns.TypeA, server)
				gg.AddNode(nodenum, Graph.ServerfailureNode)
				gg.SetRerverseflag(nodenum, Graph.ServerfailureNode)
				//gg.Setflag(domain, msgType, server, Graph.Serverfailure)
				cc.AddERROR2()
				return nil, nil
			case 3:
				fmt.Println("出现错误  Name Error")
				nodenum := gg.NodeNum(domain, dns.TypeA, server)
				gg.AddNode(nodenum, Graph.NameErrorNode)
				gg.SetRerverseflag(nodenum, Graph.NameErrorNode)
				//gg.Setflag(domain, msgType, server, Graph.NameError)
				cc.AddERROR3()
				return nil, nil
			case 4:
				fmt.Println("出现错误  不支持查询类型")
				nodenum := gg.NodeNum(domain, dns.TypeA, server)
				gg.SetRerverseflag(nodenum, Graph.NotImplementedNode)
				gg.AddNode(nodenum, Graph.NotImplementedNode)
				//gg.Setflag(domain, msgType, server, Graph.NotImplemented)
				cc.AddERROR4()
				return nil, nil
			case 5:
				fmt.Println("出现错误  Refused")

				nodenum := gg.NodeNum(gg.Domain, dns.TypeA, server)
				gg.SetRerverseflag(nodenum, Graph.RefusedNode)
				gg.AddNode(nodenum, Graph.RefusedNode)
				fmt.Println(domain, server)
				//gg.Setflag(domain, msgType, server, Graph.Refused)
				cc.AddERROR5()
				return nil, nil
			}
			//answernum:=msg.MsgHdr.

			//fmt.Println("我想要的东西", msg.Answer)

			if len(msg.Answer) == 0 {
				//fmt.Println("==================================")
				servers = servers[:0]
				if len(msg.Ns) == 0 {
					fmt.Println("没有NS记录XXXXXXXXXXXXXXXXXXXXXXXXXXXX")
					nodenum := gg.NodeNum(gg.Domain, dns.TypeA, server)
					gg.SetRerverseflag(nodenum, Graph.NoNsrecordNode)
					gg.AddNode(nodenum, Graph.NoNsrecordNode)
					return nil, nil
				}
				for _, value := range msg.Extra { //通过range可以直接得到数组元素
					//把A记录记起来，准备递归查询
					if value.Header().Rrtype == dns.TypeA {
						ns, ok := value.(*dns.A)
						if ok {
							servers = append(servers, ns.A.String())
						}
						//return responses, nil
					}
					//把AAAA记录记起来，准备递归查询
					if value.Header().Rrtype == dns.TypeAAAA {
						ns, ok := value.(*dns.AAAA)
						if ok {
							servers = append(servers, ns.AAAA.String())
						}
						//return responses, nil
					}
				}
				//fmt.Println("MSG======:", msg)
				//处理NS记录不附带IP的情况
				if len(servers) == 0 {
					fmt.Println("Not glue IP!")
					NsNotGlueIP = NsNotGlueIP[:0]
					for _, v := range msg.Ns { //通过range可以直接得到数组元素
						ns, ok := v.(*dns.NS) //ok为bool，判断是否为该类型
						if ok {
							//fmt.Println("if ok :", ok)
							NsNotGlueIP = append(NsNotGlueIP, ns.Ns)
						}
					}
					fmt.Println("NSNOTGLUEIP", NsNotGlueIP, domain)

					//====================把NS不带IP的记录也放到缓存里面,缓存特殊标记一下
					var cachevalue Cache.Cachevalue
					cachevalue.IP = NsNotGlueIP
					cachevalue.Flag = true
					cc.Add(domain, server, dns.TypeA, cachevalue)
					nodenum := gg.NodeNum(domain, dns.TypeA, server)
					for _, tempserver := range NsNotGlueIP {
						tempnodenum := gg.NodeNum(tempserver, dns.TypeA, "")
						gg.AddNode(nodenum, tempnodenum)
						gg.SetRerverseflag(tempnodenum, Graph.NsNotGlueIPNode)
					}
					//nodenum := gg.NodeNum(tempserver, dns.TypeA, "")
					for _, NS := range NsNotGlueIP {

						var GetIP1 = make([]string, 0)
						//GetIP[0] = "127.0.0.1"
						nodenum := gg.NodeNum(NS, dns.TypeA, "")
						for _, value := range root46servers {
							//nodenum, _ := Graph.NodeNum(domain, int(dns.TypeA), server)
							tempnodenum := gg.NodeNum(NS, dns.TypeA, value)
							gg.AddNode(nodenum, tempnodenum)
							d.Resolver(NS, dns.TypeA, value, gg, cc, &GetIP1)
						}
						//fmt.Println(" jijiji      ", domain, GetIP1)
						NoDup := d.RemoveDuplicates(GetIP1)

						for _, value := range NoDup {
							nodenum := gg.NodeNum(NS, dns.TypeA, value)
							nodenumA := gg.NodeNum(domain, dns.TypeA, value)
							gg.AddNode(nodenum, nodenumA)
							d.Resolver(domain, dns.TypeA, value, gg, cc, &GetIP1)
						}

						GetIP1 = make([]string, 0)
					}

					return nil, nil
				} else {
					var tempvalue Cache.Cachevalue
					tempvalue.IP = servers
					// value := Cache.GetCache(domain, server, int(dns.TypeA))
					nodenum := gg.NodeNum(domain, dns.TypeA, server)
					for _, tempserver := range tempvalue.IP {
						tempnodenum := gg.NodeNum(domain, dns.TypeA, tempserver)

						if nodenum == tempnodenum {
							fmt.Println("nodenum==tempnodenum", msg)
							gg.SetRerverseflag(nodenum, Graph.OneCircleNode)
						}
						gg.AddNode(nodenum, tempnodenum)
					}
					cc.Add(domain, server, dns.TypeA, tempvalue)
				}

				for _, value := range servers {
					//fmt.Println("递归查询：", index)
					//递归查询
					d.Resolver(domain, dns.TypeA, value, gg, cc, GetIP)
				}
				return nil, nil
			} else {
				var tempvalue Cache.Cachevalue
				//tempvalue.IP = R
				for _, value := range msg.Answer {
					//处理CNAME
					if value.Header().Rrtype == dns.TypeCNAME {
						ns, ok := value.(*dns.CNAME)
						nodenum := gg.NodeNum(domain, dns.TypeA, server)
						//把这一条放到缓存里

						nodenumCNAME := gg.NodeNum(ns.Target, dns.TypeA, "")
						gg.AddNode(nodenum, nodenumCNAME)
						gg.SetRerverseflag(nodenumCNAME, Graph.LeaveCNAMENode)

						cc.Add(domain, server, dns.TypeA, tempvalue)
						//fmt.Println("CNAME放到cache里了")
						if ok {
							for _, value := range root46servers {
								//CNAME插入图中
								//Graph.Setflag(domain, msgType, value, Graph.LeaveCNAME)
								tempnodenum := gg.NodeNum(ns.Target, msgType, value)
								gg.AddNode(nodenumCNAME, tempnodenum)

								//d.ResolverIP(ns.Target, dns.TypeCNAME, value)
								d.Resolver(ns.Target, msgType, value, gg, cc, GetIP)
							}
						}
						return nil, nil
						//return responses, nil
					}
					//打印结果A记录
					if value.Header().Rrtype == dns.TypeA {
						ns, _ := value.(*dns.A)
						nodenum := gg.NodeNum(domain, dns.TypeA, server)
						//fmt.Println("打印A记录", ns.A, gg.Domain, nodenum)
						//gg.Setflag(domain, msgType, string(ns.A), Graph.LeaveA)
						tempnodenum := gg.NodeNum(domain, msgType, ns.A.String())
						*GetIP = append(*GetIP, ns.A.To16().String())
						//fmt.Println("*GetIP", ns.A.To16().String(), *GetIP)
						fmt.Println("打印A记录", ns.A, domain, tempnodenum)
						gg.AddNode(nodenum, tempnodenum)
						gg.SetRerverseflag(tempnodenum, Graph.LeaveANode)
						//gg.AddNode(tempnodenum, Graph.LeaveANode)
						// return responses, nil
					}

					//打印结果AAAA记录
					if value.Header().Rrtype == dns.TypeAAAA {
						ns, _ := value.(*dns.AAAA)
						fmt.Println("打印AAAA记录", ns.AAAA, domain)
						*GetIP = append(*GetIP, ns.AAAA.To16().String())
						//fmt.Println("*GetIP", ns.AAAA.To16().String(), *GetIP)
						nodenum := gg.NodeNum(domain, dns.TypeA, server)
						//Graph.Setflag(domain, msgType, string(ns.AAAA), 6)
						gg.Setflag(domain, msgType, string(ns.AAAA), Graph.LeaveAAAA)
						tempnodenum := gg.NodeNum(domain, msgType, ns.AAAA.String())
						gg.AddNode(nodenum, tempnodenum)
						gg.SetRerverseflag(tempnodenum, Graph.LeaveAAAANode)
						//gg.AddNode(tempnodenum, Graph.LeaveAAAANode)
						// return responses, nil
					}
				}
				//fmt.Println("------------------------------------Return")
				return nil, nil
			}
		}
	}
}

// Trace  类似于 dig +trace,把所有根都遍历一遍
func (d *Dig) Trace(domain string, Qtype uint16, gg *Graph.GraphStruct, cc *Cache.CacheStruct) ([]TraceResponse, error) {

	//gg.Init()
	//gg.Init()

	GetIP := make([]string, 0)
	var trace = make([]TraceResponse, 0)
	for index, value := range root46servers {
		fmt.Println("ROOT：", index)

		//画图
		num := gg.NodeNum(domain, Qtype, value)
		fmt.Println("ROOTnum：", num, domain, Qtype, value)
		gg.AddNode(0, num)
		d.Resolver(domain, dns.TypeA, value, gg, cc, &GetIP)

	}
	return trace, nil
	//return d.TraceForRecord(domain, dns.TypeA, root46servers[5])
}

func Mapmerge(map1 map[string]string, map2 map[string]string) map[string]string {
	x := map1
	y := map2
	n := make(map[string]string)
	for i, v := range x {
		for j, w := range y {
			if i == j {
				n[i] = w

			} else {
				if _, ok := n[i]; !ok {
					n[i] = v
				}
				if _, ok := n[j]; !ok {
					n[j] = w
				}
			}
		}
	}
	return n
}
