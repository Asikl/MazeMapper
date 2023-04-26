package Graph

import (
	"fmt"

	//"hash/maphash"
	"github.com/yourbasic/graph"
)

const (
	NoVisit        = 0
	Corrupt        = 1
	Serverfailure  = 2
	NXDOMAIN       = 3
	NotImplemented = 4
	Refused        = 5
	Timeout        = 6
	NoNsrecord     = 7
	NsNotGlueIP    = 71
	IPerror        = 8
	IDMisMatch     = 9
	LeaveA         = 10
	LeaveAAAA      = 11
	LeaveCNAME     = 12
	Common         = 100

	CorruptNode        = 980
	ServerfailureNode  = 981
	NXDOMAINNode       = 982
	NotImplementedNode = 983
	RefusedNode        = 984
	TimeoutNode        = 985
	NoNsrecordNode     = 986
	NsNotGlueIPNode    = 987
	IPerrorNode        = 988
	IDMisMatchNode     = 989
	LeaveANode         = 990
	LeaveAAAANode      = 991
	LeaveCNAMENode     = 992
	OneCircleNode      = 993
)

type GraphStruct struct {
	Domain       string
	Num          int
	GraphMap     map[Key]Value
	GraphReverse map[int]KeyRevserse
	Domaingraph  *graph.Mutable
}

type Key struct {
	Domain string
	Ip     string // compatible with ipv4/ipv6, could be a user-defined ip type
	Qtype  uint16
}

type KeyRevserse struct {
	Domain string
	Ip     string // compatible with ipv4/ipv6, could be a user-defined ip type
	Qtype  uint16
	Flag   int
}

type Value struct {
	ID   int
	flag int //标记该节点访问了没，以及该节点的状态
}

type DomainGraph struct {
	Num          int //图中节点的序号
	GraphMap     map[Key]Value
	GraphReverse map[int]Key
}

// var Num int
// var GraphMap = make(map[Key]Value, 0)
// var GraphReverse = make(map[int]Key, 0)
// var Domaingraph = graph.New(1000)

// Map[domain] graph
// Map[domain] id
// Hash domain ip qtype

func Init(g *GraphStruct, Domain string) {
	var Beginkey Key
	var BeginValue Value
	var BeginkeyReverse KeyRevserse
	g.Domain = Domain
	g.GraphMap = make(map[Key]Value, 0)
	g.GraphReverse = make(map[int]KeyRevserse, 0)
	g.Domaingraph = graph.New(1000)

	Beginkey.Domain = ""
	BeginValue.ID = 0

	BeginkeyReverse.Domain = ""

	//先插入一个开始节点
	g.GraphMap[Beginkey] = BeginValue
	g.GraphReverse[0] = BeginkeyReverse
	g.Num = 0
}

// 获取 Domaingraph 字段的值的方法
func (g *GraphStruct) GetDomaingraph() *graph.Mutable {
	return g.Domaingraph
}

func (g *GraphStruct) AddNum() {
	//图中节点的个数+1
	g.Num++
}

func (g *GraphStruct) SubNum() {
	//图中节点的个数+1
	g.Num--
}
func (g *GraphStruct) GetNum() (num int) {
	//图中节点的个数+1
	return g.Num
}

func (g *GraphStruct) AddNode(num1 int, num2 int) {
	//fmt.Println("在图中插入节点", num1, num2)

	//Cache.AddEdge(1)
	g.Domaingraph.Add(num1, num2) //Add inserts a directed edge from v to w with zero cost. It removes the previous cost if this edge already exists.
}

func (g *GraphStruct) Getflag(domain string, Qtype uint16, Ip string) (flag int) {
	var key Key
	key.Domain = domain
	key.Qtype = Qtype
	key.Ip = Ip
	if value, ok := g.GraphMap[key]; ok {
		return value.flag
	} else {
		return 0 //没有该节点暂时认为是为访问状态
	}
	//return GraphMap[key].flag
}

// func SetNodetype(domain string, Qtype uint16, Ip string, Nodetype int) {
// 	var temp Key
// 	temp.Domain = domain
// 	temp.Ip = Ip
// 	temp.Qtype = Qtype
// 	Value := GraphMap[temp]
// 	Value.Nodetype = Nodetype
// 	GraphMap[temp] = Value
// }

func (g *GraphStruct) Setflag(domain string, Qtype uint16, Ip string, flag int) {
	var temp Key
	temp.Domain = domain
	temp.Ip = Ip
	temp.Qtype = Qtype
	Value := g.GraphMap[temp]
	Value.flag = flag
	g.GraphMap[temp] = Value
	//GraphMap[temp].flag = flag
}

func (g *GraphStruct) SetRerverseflag(nodenum int, flag int) {

	temp := g.GraphReverse[nodenum]
	temp.Flag = flag
	g.GraphReverse[nodenum] = temp
	//GraphMap[temp].flag = flag
}

func (g *GraphStruct) NodeNum(domain string, Qtype uint16, Ip string) (Nodenum int) {
	var Nodeint int
	var temp Key
	var tempReverse KeyRevserse
	temp.Domain = domain
	temp.Ip = Ip
	temp.Qtype = Qtype

	tempReverse.Domain = domain
	tempReverse.Ip = Ip
	tempReverse.Qtype = Qtype
	//记录总边的个数  Edge++

	if value, ok := g.GraphMap[temp]; ok {
		Nodeint = value.ID
		//fmt.Println("图中已有该节点，只需要加一条边")
		return Nodeint
	} else {
		//fmt.Println("图中没有该节点，需要加一个点")
		var value Value
		g.AddNum()
		value.ID = g.GetNum()
		//fmt.Println("value.IDKKKKKKKKKKKKK", GetNum(), value.ID)

		g.GraphReverse[value.ID] = tempReverse
		value.flag = 0 //该节点未访问
		g.GraphMap[temp] = value
		return value.ID
	}
}

func (g *GraphStruct) Dump() {
	fmt.Println("把图输出", g.Domaingraph)

}

func (g *GraphStruct) DumpGraphMap() {
	fmt.Println("把图输出", g.GraphMap)
}

func (g *GraphStruct) DumpGraphReverse() {
	fmt.Println("把逆图输出", g.GraphReverse)
}
