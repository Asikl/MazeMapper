package Cache

import (
	// "github.com/miekg/dns"
	"fmt"
)

type Cachekey struct {
	CacheDomain string
	QType       uint16
	CacheIp     string
}
type Cachevalue struct {
	flag bool //来标志该缓存的NS记录带不带IP
	NS   []string
	IP   []string
}

type ERROR struct {
	Edge   int //边的个数
	ERROR0 int //没有错误查询
	ERROR1 int //格式错误
	ERROR2 int //Server failure
	ERROR3 int //Name Error
	ERROR4 int //Not Implemented
	ERROR5 int //Refused
	ERROR6 int //time out
	ERROR7 int //数据包中没有NS记录，无法继续进行下去
	//ERROR8 int
}

// 先定义一个比较low的全局缓存cache
var cacheFIX = make(map[Cachekey]Cachevalue, 0)

// func InitCache() map[Cachekey]Cachevalue {
// 	// 先定义一个比较low的全局缓存cache
// 	var cacheFIX = make(map[Cachekey]Cachevalue, 0)

// 	return cacheFIX
// }

var DomainERROR ERROR

func ERROR0Init() {
	DomainERROR.ERROR0 = 0
}
func GetERROR0() (num int) {
	return DomainERROR.ERROR0
}
func AddERROR0() {
	DomainERROR.ERROR0++
}

func ERROR1Init() {
	DomainERROR.ERROR1 = 0
}
func GetERROR1() (num int) {
	return DomainERROR.ERROR1
}
func AddERROR1() {
	DomainERROR.ERROR1++
}

func ERROR2Init() {
	DomainERROR.ERROR2 = 0
}
func GetERROR2() (num int) {
	return DomainERROR.ERROR2
}
func AddERROR2() {
	DomainERROR.ERROR2++
}

func ERROR3Init() {
	DomainERROR.ERROR3 = 0
}
func GetERROR3() (num int) {
	return DomainERROR.ERROR3
}
func AddERROR3() {
	DomainERROR.ERROR3++
}

func ERROR4Init() {
	DomainERROR.ERROR4 = 0
}
func GetERROR4() (num int) {
	return DomainERROR.ERROR4
}
func AddERROR4() {
	DomainERROR.ERROR4++
}

func ERROR5Init() {
	DomainERROR.ERROR5 = 0
}
func GetERROR5() (num int) {
	return DomainERROR.ERROR5
}
func AddERROR5() {
	DomainERROR.ERROR5++
}

func ERROR6Init() {
	DomainERROR.ERROR6 = 0
}
func GetERROR6() (num int) {
	return DomainERROR.ERROR6
}
func AddERROR6() {
	DomainERROR.ERROR6++
}

func ERROR7Init() {
	DomainERROR.ERROR7 = 0
}
func GetERROR7() (num int) {
	return DomainERROR.ERROR7
}
func AddERROR7() {
	DomainERROR.ERROR7++
}

func EdegeInit() {
	DomainERROR.Edge = 0
}

func AddEdge(num int) {
	DomainERROR.Edge += num
}
func GetEdge() (num int) {
	return DomainERROR.Edge
}

func InitERROR() {
	ERROR1Init()
	ERROR2Init()
	ERROR3Init()
	ERROR4Init()
	ERROR5Init()
	ERROR6Init()
	ERROR7Init()
}

func Add(domain string, server string, Qtype uint16, value Cachevalue) {

	var temp Cachekey
	temp.CacheDomain = domain
	temp.CacheIp = server
	temp.QType = Qtype
	cacheFIX[temp] = value
}

func GetCache(domain string, server string, Qtype uint16) (value Cachevalue) {
	var key Cachekey
	key.CacheDomain = domain
	key.CacheIp = server
	key.QType = Qtype
	return cacheFIX[key]
}

func Has(domain string, server string, Qtype uint16) (flag bool) {
	var temp Cachekey
	temp.CacheDomain = domain
	temp.CacheIp = server
	temp.QType = Qtype
	if _, ok := cacheFIX[temp]; ok {
		//fmt.Println("缓存中已有，是重复节点！")
		return true
		//存在
	} else {
		return false
	}
}

func Dump() {
	fmt.Println(cacheFIX)
}
