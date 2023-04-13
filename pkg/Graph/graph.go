package Graph

import (
	"fmt"
	"hello/pkg/Cache"

	//"hash/maphash"
	"github.com/yourbasic/graph"
)

type Key struct {
	Domain string
	Ip     string // compatible with ipv4/ipv6, could be a user-defined ip type
	Qtype  uint16
}

type Value struct {
	ID       int
	Nodetype int //标记节点状态，timeout,Corrupt等等
	//Corrupt 格式错误 1
	//Server failure 2
	//Name Error 3
	//Not Implemented 4
	//Refused 5
	//timeout 6   发不过去61     收不过来62          根本无法建立连接63
	//数据包中没有NS记录，无法继续进行下去 7
	//目标IP有问题 8
	//发的数据包和收到的数据包ID不匹配  9
	//正常节点  100

	//叶子节点 50  A
	//叶子节点 51  AAAA
	//叶子节点 52  CNAME
	flag int //标记该节点访问了没，以及该节点的状态
}

var Num int
var GraphMap = make(map[Key]Value, 0)
var Domaingraph = graph.New(1000)

// Map[domain] graph
// Map[domain] id
// Hash domain ip qtype

func Init() {
	var Beginkey Key
	var BeginValue Value

	Beginkey.Domain = "Begin"
	BeginValue.ID = 0
	//先插入一个开始节点
	GraphMap[Beginkey] = BeginValue
	Num = 0
}

func AddNum() {
	//图中节点的个数+1
	Num++
}

func SubNum() {
	//图中节点的个数+1
	Num--
}
func GetNum() (num int) {
	//图中节点的个数+1
	return Num
}

func AddNode(num1 int, num2 int) {
	Cache.AddEdge(1)
	Domaingraph.Add(num1, num2) //Add inserts a directed edge from v to w with zero cost. It removes the previous cost if this edge already exists.
}

func Getflag(domain string, Qtype uint16, Ip string) (flag int) {
	var key Key
	key.Domain = domain
	key.Qtype = Qtype
	key.Ip = Ip

	if value, ok := GraphMap[key]; ok {
		return value.flag
	} else {
		return 0 //没有该节点暂时认为是为访问状态
	}
	//return GraphMap[key].flag
}

func SetNodetype(domain string, Qtype uint16, Ip string, Nodetype int) {
	var temp Key
	temp.Domain = domain
	temp.Ip = Ip
	temp.Qtype = Qtype
	Value := GraphMap[temp]
	Value.Nodetype = Nodetype
	GraphMap[temp] = Value

}

func Setflag(domain string, Qtype uint16, Ip string, flag int) {
	var temp Key
	temp.Domain = domain
	temp.Ip = Ip
	temp.Qtype = Qtype
	Value := GraphMap[temp]
	Value.flag = flag
	GraphMap[temp] = Value

	//GraphMap[temp].flag = flag
}

func NodeNum(domain string, Qtype uint16, Ip string) (Nodenum int, flag bool) {
	var Nodeint int
	var temp Key
	temp.Domain = domain
	temp.Ip = Ip
	temp.Qtype = Qtype
	//记录总边的个数  Edge++

	if value, ok := GraphMap[temp]; ok {
		Nodeint = value.ID
		//fmt.Println("图中已有该节点，只需要加一条边")
		return Nodeint, true
	} else {
		//fmt.Println("图中没有该节点，需要加一个点")
		var value Value
		AddNum()
		value.ID = GetNum()
		//fmt.Println("value.IDKKKKKKKKKKKKK", GetNum(), value.ID)

		value.flag = 0 //该节点未访问
		GraphMap[temp] = value
		return value.ID, false
	}
}

func Dump() {
	fmt.Println("把图输出", Domaingraph)

}

func DumpGraphMap() {
	fmt.Println("把图输出", GraphMap)
}
