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
	LocalAddr        string
	RemoteAddr       string
	BackupRemoteAddr string
	DialTimeout      time.Duration
	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	Protocol         string
	Retry            int
}

// 存解决NS记录得到的结果，go语言好像没有set，我们用map来模拟set
var NSgetIP = make(map[string]string, 0)

// TraceResponse  dig +trace 响应
type TraceResponse struct {
	Server   string
	ServerIP string
	Msg      *dns.Msg
	//type1    TypeDNSKEY.Msg
}

var responses = make([]TraceResponse, 0)

// Resolver
func (d *Dig) Resolver(domain string, msgType uint16, distination string) ([]TraceResponse, error) {
	//var responses = make([]TraceResponse, 0)
	var servers = make([]string, 0)
	//cacheFIX := make(map[cachekey]cachevalue)
	var NsNotGlueIP = make([]string, 0)
	server := distination
	//死循环
	for {
		if err := d.SetDNS(server); err != nil {
			fmt.Println("IP本身有问题", err)
			Graph.SetNodetype(domain, msgType, distination, 8)
			return nil, nil
		}
		if Cache.Has(domain, server, dns.TypeA) {
			//先得到cache里的内容
			value := Cache.GetCache(domain, server, dns.TypeA)
			nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
			//Cache.AddEdge(1)
			for _, tempserver := range value.IP {
				//Cache.AddEdge(1)
				tempnodenum, flag := Graph.NodeNum(domain, dns.TypeA, tempserver)
				if flag {
					//fmt.Println("两个节点图里面都有了")
					Graph.AddNode(nodenum, tempnodenum)
					continue
				} else {
					//fmt.Println("只有一个节点图里面有")
					Graph.AddNode(nodenum, tempnodenum)
				}
			}
			//处理指向节点
			return nil, nil
		} else {
			flag := Graph.Getflag(domain, dns.TypeA, server)
			if flag != 0 {
				//fmt.Println("该节点已访问")
				return nil, nil
			}
			msg, err, num := d.GetMsg(msgType, domain) //GetMsg
			switch num {
			case 100: //正常节点
				Graph.SetNodetype(domain, msgType, server, 100)
			case 61:
				Graph.SetNodetype(domain, msgType, server, 61)
			case 62:
				Graph.SetNodetype(domain, msgType, server, 62)
			case 63:
				Graph.SetNodetype(domain, msgType, server, 63)
			case 9:
				Graph.SetNodetype(domain, msgType, server, 9)
			}
			Graph.Setflag(domain, msgType, server, 1) //节点标记为已访问
			if err != nil {
				Cache.AddERROR6()
				fmt.Println("出现错误  ", fmt.Errorf("%s,%v", server, err))
				//Graph.Setflag(domain, msgType, server, 5) //节点标记为错误节点
				return responses, fmt.Errorf("%s:%v", server, err)
			}
			//Graph.Setflag(domain, msgType, server, 1) //节点标记为已访问
			var rsp TraceResponse
			rsp.Server = server
			rsp.ServerIP = server
			rsp.Msg = msg
			responses = append(responses, rsp)

			switch msg.Rcode {
			case 0:
				//fmt.Println("NOERROR")
				Cache.AddERROR0()
			case 1:
				fmt.Println("出现错误  格式错误")

				Graph.SetNodetype(domain, msgType, server, 1)
				Cache.AddERROR1()
				return nil, nil
			case 2:
				fmt.Println("出现错误  Server Failure")
				Graph.SetNodetype(domain, msgType, server, 2)
				Cache.AddERROR2()
				return nil, nil
			case 3:
				fmt.Println("出现错误  Name Error")
				Graph.SetNodetype(domain, msgType, server, 3)
				Cache.AddERROR3()
				return nil, nil
			case 4:
				fmt.Println("出现错误  不支持查询类型")
				Graph.SetNodetype(domain, msgType, server, 4)
				Cache.AddERROR4()
				return nil, nil
			case 5:
				fmt.Println("出现错误  Refused")
				Graph.SetNodetype(domain, msgType, server, 5)
				Cache.AddERROR5()
				return nil, nil
			}
			//answernum:=msg.MsgHdr.

			//fmt.Println("我想要的东西", msg.Answer)

			if len(msg.Answer) == 0 {
				//fmt.Println("==================================")
				servers = servers[:0]
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
					fmt.Println("NSNOTGLUEIP", NsNotGlueIP)

					//====================把NS不带IP的记录也放到缓存里面
					var cachevalue Cache.Cachevalue
					cachevalue.IP = NsNotGlueIP
					Cache.Add(domain, server, dns.TypeA, cachevalue)
					nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
					for _, tempserver := range NsNotGlueIP {
						tempnodenum, flag := Graph.NodeNum(tempserver, dns.TypeA, "")
						if flag {
							Graph.AddNode(nodenum, tempnodenum)
							continue
						} else {
							//fmt.Println("只有一个节点图里面有")
							Graph.AddNode(nodenum, tempnodenum)
						}
					}
					//用于DEBUG
					//fmt.Println("成功插入cache")
					//fmt.Println("打印cache", cacheFIX)

					//完全没有NS记录，无法继续进行查询
					if len(NsNotGlueIP) == 0 {
						Cache.AddERROR7()
						Graph.SetNodetype(domain, msgType, server, 7)
						fmt.Println("出现错误  ", fmt.Sprintf("从%s得到的数据没有NS记录", server))
						return nil, nil
					} else {
						//nodenum, _ := Graph.NodeNum(domain, int(dns.TypeA), server)
						//把所有NS不带IP的情况都给从根部开始重新遍历一遍
						for _, NS := range NsNotGlueIP {
							nodenum, _ := Graph.NodeNum(NS, dns.TypeA, "")
							for _, value := range root46servers {
								//nodenum, _ := Graph.NodeNum(domain, int(dns.TypeA), server)
								tempnodenum, flag := Graph.NodeNum(NS, dns.TypeA, value)
								if flag {
									//fmt.Println("两个节点图里面都有了")
									Graph.AddNode(nodenum, tempnodenum)
									continue
								} else {
									//fmt.Println("只有一个节点图里面有")
									Graph.AddNode(nodenum, tempnodenum)
								}

								d.ResolverIP(NS, dns.TypeA, value)
							}
						}

						//serverstemp, _ := d.TraceIP(NsNotGlueIP[0])
						//把解析得到的IP供之前的查询继续进行下去

						for index, _ := range NSgetIP {
							servers = append(servers, index)
						}
						//servers = append(servers)

						fmt.Println("成功解决Ns不附带IP的情况", NSgetIP)

						if len(servers) == 0 {
							fmt.Println("没有可以继续访问的节点")
							Graph.SetNodetype(domain, msgType, server, 7)
						}

						//domainnode, _ := Graph.NodeNum(domain, dns.TypeA, server)

						for index, value := range NSgetIP {
							tempnodenum, _ := Graph.NodeNum(domain, dns.TypeA, value)
							domainnode, flag := Graph.NodeNum(index, dns.TypeA, value)
							if flag {
								//fmt.Println("两个节点图里面都有了")
								Graph.AddNode(domainnode, tempnodenum)
							} else {
								//fmt.Println("只有一个节点图里面有")
								Graph.AddNode(domainnode, tempnodenum)
							}
							d.Resolver(domain, dns.TypeA, value)
						}
						return nil, nil
						//return nil, nil
						//DEBUG
					}
				} else {
					var tempvalue Cache.Cachevalue
					tempvalue.IP = servers
					// value := Cache.GetCache(domain, server, int(dns.TypeA))
					nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
					for _, tempserver := range tempvalue.IP {
						tempnodenum, flag := Graph.NodeNum(domain, dns.TypeA, tempserver)
						if flag {
							//fmt.Println("两个节点图里面都有了")
							//fmt.Println("NNNNNNNNNNNNNN", nodenum, tempnodenum)
							Graph.AddNode(nodenum, tempnodenum)
							continue
						} else {
							//fmt.Println("只有一个节点图里面有")
							Graph.AddNode(nodenum, tempnodenum)
						}
					}
					Cache.Add(domain, server, dns.TypeA, tempvalue)
					//用于DEBUG
					//fmt.Println("成功插入cache")
					//fmt.Println("打印cache", cacheFIX)
				}

				// fmt.Println("XXXXXXXXXXXXXXXXX",serv)
				//fmt.Println("Servers查BUG", servers)
				for _, value := range servers {
					//fmt.Println("递归查询：", index)
					//递归查询
					d.Resolver(domain, dns.TypeA, value)
				}
				return nil, nil
			} else {
				//fmt.Println("*****************************************")
				//fmt.Println("msg.Authoritative RESPONSE  TRUE", responses)
				var tempvalue Cache.Cachevalue
				//tempvalue.IP = R
				for _, value := range msg.Answer {
					//处理CNAME
					if value.Header().Rrtype == dns.TypeCNAME {
						ns, ok := value.(*dns.CNAME)
						nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
						//把这一条放到缓存里

						Cache.Add(domain, server, dns.TypeA, tempvalue)
						fmt.Println("CNAME放到cache里了")
						if ok {
							for _, value := range root46servers {
								//CNAME插入图中
								Graph.SetNodetype(domain, msgType, value, 52)
								tempnodenum, flag := Graph.NodeNum(ns.Target, msgType, value)
								if flag {
									//fmt.Println("两个节点图里面都有了")
									Graph.AddNode(nodenum, tempnodenum)
									continue
								} else {
									//fmt.Println("只有一个节点图里面有")
									Graph.AddNode(nodenum, tempnodenum)
								}
								//d.ResolverIP(ns.Target, dns.TypeCNAME, value)
								d.Resolver(ns.Target, msgType, value)
							}
						}
						return responses, nil
						//return responses, nil
					}
					//打印结果A记录
					if value.Header().Rrtype == dns.TypeA {
						ns, _ := value.(*dns.A)
						fmt.Println("打印A记录", ns.A, domain)
						nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
						Graph.SetNodetype(domain, msgType, string(ns.A), 50)
						tempnodenum, flag := Graph.NodeNum(domain, msgType, ns.A.String())
						if flag {
							//fmt.Println("两个节点图里面都有了")
							Graph.AddNode(nodenum, tempnodenum)
							continue
						} else {
							//fmt.Println("只有一个节点图里面有")
							Graph.AddNode(nodenum, tempnodenum)
						}
						// return responses, nil
					}

					//打印结果AAAA记录
					if value.Header().Rrtype == dns.TypeAAAA {
						ns, _ := value.(*dns.AAAA)
						fmt.Println("打印AAAA记录", ns.AAAA, domain)

						nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
						//Graph.Setflag(domain, msgType, string(ns.AAAA), 6)
						Graph.SetNodetype(domain, msgType, string(ns.AAAA), 51)
						tempnodenum, flag := Graph.NodeNum(domain, msgType, ns.AAAA.String())
						if flag {
							//fmt.Println("两个节点图里面都有了")
							Graph.AddNode(nodenum, tempnodenum)
							continue
						} else {
							//fmt.Println("只有一个节点图里面有")
							Graph.AddNode(nodenum, tempnodenum)
						}
						// return responses, nil
					}
				}
				//fmt.Println("------------------------------------Return")
				return responses, nil
			}
		}
	}
}

// Trace  类似于 dig +trace,把所有根都遍历一遍
func (d *Dig) Trace(domain string, Qtype uint16) ([]TraceResponse, error) {
	// for _, value := range root46servers {
	// 	num, _ := Graph.NodeNum(domain, Qtype, value)
	// 	Graph.AddNode(0, num)
	// }
	//Graph.Dump()
	var trace = make([]TraceResponse, 0)
	for index, value := range root46servers {
		fmt.Println("ROOT：", index)

		//画图
		num, _ := Graph.NodeNum(domain, Qtype, value)
		Graph.AddNode(0, num)
		race, _ := d.Resolver(domain, dns.TypeA, value)
		trace = append(trace, race...)
	}
	return trace, nil
	// d.init()
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

// ResolverIP

func (d *Dig) ResolverIP(domain string, msgType uint16, distination string) ([]TraceResponse, error) {
	//var responses = make([]TraceResponse, 0)
	var servers = make([]string, 0)
	//cacheFIX := make(map[cachekey]cachevalue)
	var NsNotGlueIP = make([]string, 0)
	server := distination
	//死循环
	for {
		if err := d.SetDNS(server); err != nil {
			fmt.Println("IP本身有问题", err)
			Graph.SetNodetype(domain, msgType, distination, 8)
			return nil, nil
		}
		if Cache.Has(domain, server, dns.TypeA) {
			//先得到cache里的内容
			value := Cache.GetCache(domain, server, dns.TypeA)
			nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
			//Cache.AddEdge(1)
			for _, tempserver := range value.IP {
				//Cache.AddEdge(1)
				tempnodenum, flag := Graph.NodeNum(domain, dns.TypeA, tempserver)
				if flag {
					//fmt.Println("两个节点图里面都有了")
					Graph.AddNode(nodenum, tempnodenum)
					continue
				} else {
					//fmt.Println("只有一个节点图里面有")
					Graph.AddNode(nodenum, tempnodenum)
				}
			}
			//处理指向节点
			return nil, nil
		} else {
			flag := Graph.Getflag(domain, dns.TypeA, server)
			if flag != 0 {
				fmt.Println("该节点已访问")
				return nil, nil
			}
			msg, err, num := d.GetMsg(msgType, domain) //GetMsg
			switch num {
			case 100: //正常节点
				Graph.SetNodetype(domain, msgType, server, 100)
			case 61:
				Graph.SetNodetype(domain, msgType, server, 61)
			case 62:
				Graph.SetNodetype(domain, msgType, server, 62)
			case 63:
				Graph.SetNodetype(domain, msgType, server, 63)
			case 9:
				Graph.SetNodetype(domain, msgType, server, 9)
			}
			Graph.Setflag(domain, msgType, server, 1) //节点标记为已访问
			if err != nil {
				Cache.AddERROR6()
				fmt.Println("出现错误  ", fmt.Errorf("%s,%v", server, err))
				//Graph.Setflag(domain, msgType, server, 5) //节点标记为错误节点
				return responses, fmt.Errorf("%s:%v", server, err)
			}
			//Graph.Setflag(domain, msgType, server, 1) //节点标记为已访问
			var rsp TraceResponse
			rsp.Server = server
			rsp.ServerIP = server
			rsp.Msg = msg
			responses = append(responses, rsp)

			switch msg.Rcode {
			case 0:
				//fmt.Println("NOERROR")
				Cache.AddERROR0()
			case 1:
				fmt.Println("出现错误  格式错误")

				Graph.SetNodetype(domain, msgType, server, 1)
				Cache.AddERROR1()
				return nil, nil
			case 2:
				fmt.Println("出现错误  Server Failure")
				Graph.SetNodetype(domain, msgType, server, 2)
				Cache.AddERROR2()
				return nil, nil
			case 3:
				fmt.Println("出现错误  Name Error")
				Graph.SetNodetype(domain, msgType, server, 3)
				Cache.AddERROR3()
				return nil, nil
			case 4:
				fmt.Println("出现错误  不支持查询类型")
				Graph.SetNodetype(domain, msgType, server, 4)
				Cache.AddERROR4()
				return nil, nil
			case 5:
				fmt.Println("出现错误  Refused")
				Graph.SetNodetype(domain, msgType, server, 5)
				Cache.AddERROR5()
				return nil, nil
			}
			//answernum:=msg.MsgHdr.

			//fmt.Println("我想要的东西", msg.Answer)

			if len(msg.Answer) == 0 {
				//fmt.Println("==================================")
				servers = servers[:0]
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
					fmt.Println("NSNOTGLUEIP", NsNotGlueIP)

					//====================把NS不带IP的记录也放到缓存里面
					var cachevalue Cache.Cachevalue
					cachevalue.IP = NsNotGlueIP
					Cache.Add(domain, server, dns.TypeA, cachevalue)

					nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
					for _, tempserver := range NsNotGlueIP {
						tempnodenum, flag := Graph.NodeNum(tempserver, dns.TypeA, "")
						if flag {
							Graph.AddNode(nodenum, tempnodenum)
							continue
						} else {
							//fmt.Println("只有一个节点图里面有")
							Graph.AddNode(nodenum, tempnodenum)
						}
					}
					//用于DEBUG
					//fmt.Println("成功插入cache")
					//fmt.Println("打印cache", cacheFIX)

					//完全没有NS记录，无法继续进行查询
					if len(NsNotGlueIP) == 0 {
						Cache.AddERROR7()
						Graph.SetNodetype(domain, msgType, server, 7)
						fmt.Println("出现错误  ", fmt.Sprintf("从%s得到的数据没有NS记录", server))
						return nil, nil
					} else {
						//nodenum, _ := Graph.NodeNum(domain, int(dns.TypeA), server)
						//把所有NS不带IP的情况都给从根部开始重新遍历一遍
						for _, NS := range NsNotGlueIP {
							nodenum, _ := Graph.NodeNum(NS, dns.TypeA, "")
							for _, value := range root46servers {
								//nodenum, _ := Graph.NodeNum(domain, int(dns.TypeA), server)
								tempnodenum, flag := Graph.NodeNum(NS, dns.TypeA, value)
								if flag {
									//fmt.Println("两个节点图里面都有了")
									Graph.AddNode(nodenum, tempnodenum)
									continue
								} else {
									//fmt.Println("只有一个节点图里面有")
									Graph.AddNode(nodenum, tempnodenum)
								}
								d.ResolverIP(NS, dns.TypeA, value)
							}

						}
						//serverstemp, _ := d.TraceIP(NsNotGlueIP[0])
						//把解析得到的IP供之前的查询继续进行下去
						//把解析得到的IP供之前的查询继续进行下去

						for index, _ := range NSgetIP {
							servers = append(servers, index)
						}
						//servers = append(servers)

						fmt.Println("成功解决Ns不附带IP的情况", NSgetIP)

						if len(servers) == 0 {
							fmt.Println("没有可以继续访问的节点")
						}

						domainnode, _ := Graph.NodeNum(domain, dns.TypeA, server)

						for index, value := range NSgetIP {
							tempnodenum, flag := Graph.NodeNum(index, dns.TypeA, value)
							if flag {
								//fmt.Println("两个节点图里面都有了")
								Graph.AddNode(domainnode, tempnodenum)
							} else {
								//fmt.Println("只有一个节点图里面有")
								Graph.AddNode(domainnode, tempnodenum)
							}
							d.ResolverIP(domain, dns.TypeA, value)
						}

						return nil, nil
						//return nil, nil
						//DEBUG
					}
				} else {
					var tempvalue Cache.Cachevalue
					tempvalue.IP = servers
					// value := Cache.GetCache(domain, server, int(dns.TypeA))
					nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
					for _, tempserver := range tempvalue.IP {
						tempnodenum, flag := Graph.NodeNum(domain, dns.TypeA, tempserver)
						if flag {
							//fmt.Println("两个节点图里面都有了")
							//fmt.Println("NNNNNNNNNNNNNN", nodenum, tempnodenum)
							Graph.AddNode(nodenum, tempnodenum)
							continue
						} else {
							//fmt.Println("只有一个节点图里面有")
							Graph.AddNode(nodenum, tempnodenum)
						}
					}
					Cache.Add(domain, server, dns.TypeA, tempvalue)
					//用于DEBUG
					//fmt.Println("成功插入cache")
					//fmt.Println("打印cache", cacheFIX)
				}

				// fmt.Println("XXXXXXXXXXXXXXXXX",serv)
				//fmt.Println("Servers查BUG", servers)
				for _, value := range servers {
					//fmt.Println("递归查询：", index)
					//递归查询
					d.Resolver(domain, dns.TypeA, value)
				}
				return nil, nil
			} else {
				//fmt.Println("*****************************************")
				//fmt.Println("msg.Authoritative RESPONSE  TRUE", responses)
				var tempvalue Cache.Cachevalue
				//tempvalue.IP = R
				for _, value := range msg.Answer {
					//处理CNAME
					if value.Header().Rrtype == dns.TypeCNAME {
						ns, ok := value.(*dns.CNAME)
						nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
						//把这一条放到缓存里

						Cache.Add(domain, server, dns.TypeA, tempvalue)
						fmt.Println("CNAME放到cache里了")
						if ok {
							for _, value := range root46servers {
								//CNAME插入图中
								Graph.SetNodetype(domain, msgType, value, 52)
								tempnodenum, flag := Graph.NodeNum(ns.Target, msgType, value)
								if flag {
									//fmt.Println("两个节点图里面都有了")
									Graph.AddNode(nodenum, tempnodenum)
									continue
								} else {
									//fmt.Println("只有一个节点图里面有")
									Graph.AddNode(nodenum, tempnodenum)
								}
								//d.ResolverIP(ns.Target, dns.TypeCNAME, value)
								d.ResolverIP(ns.Target, msgType, value)
							}
						}
						return responses, nil
						//return responses, nil
					}
					//打印结果A记录
					if value.Header().Rrtype == dns.TypeA {
						ns, _ := value.(*dns.A)
						fmt.Println("打印A记录", ns.A, domain)
						NSgetIP[domain] = string(ns.A)
						nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
						Graph.SetNodetype(domain, msgType, string(ns.A), 50)
						tempnodenum, flag := Graph.NodeNum(domain, msgType, ns.A.String())
						if flag {
							//fmt.Println("两个节点图里面都有了")
							Graph.AddNode(nodenum, tempnodenum)
							continue
						} else {
							//fmt.Println("只有一个节点图里面有")
							Graph.AddNode(nodenum, tempnodenum)
						}
						// return responses, nil
					}

					//打印结果AAAA记录
					if value.Header().Rrtype == dns.TypeAAAA {
						ns, _ := value.(*dns.AAAA)
						fmt.Println("打印AAAA记录", ns.AAAA, domain)
						NSgetIP[domain] = string(ns.AAAA)

						nodenum, _ := Graph.NodeNum(domain, dns.TypeA, server)
						//Graph.Setflag(domain, msgType, string(ns.AAAA), 6)
						Graph.SetNodetype(domain, msgType, string(ns.AAAA), 51)
						tempnodenum, flag := Graph.NodeNum(domain, msgType, ns.AAAA.String())
						if flag {
							//fmt.Println("两个节点图里面都有了")
							Graph.AddNode(nodenum, tempnodenum)
							continue
						} else {
							//fmt.Println("只有一个节点图里面有")
							Graph.AddNode(nodenum, tempnodenum)
						}
						// return responses, nil
					}
				}
				//fmt.Println("------------------------------------Return")
				return responses, nil
			}
		}
	}
}
