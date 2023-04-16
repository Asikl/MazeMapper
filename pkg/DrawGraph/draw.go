package DrawGraph

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	// "os"
	// "os/exec"
	"hello/pkg/Graph"

	"github.com/yourbasic/graph"
)

func Test() {
	fmt.Println("TESTTESTTEST")
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 使用bufio.NewReader创建缓冲读取器
	reader := bufio.NewReader(file)

	// 逐行读取文件内容
	for {
		// 使用ReadString('\n')方法读取一行文本
		line, err := reader.ReadString('\n')
		if err != nil {
			// 如果遇到文件末尾或者读取错误，则退出循环
			break
		}
		// 打印读取到的一行文本
		fmt.Print(line)
	}
}

type DrawStruct struct {
	path string
}

func Visual(domain string, kk *Graph.GraphStruct) {
	fmt.Println("开始画图")
	//Path := "\"DomainPicture\""
	Path := "./DomainPicture/"
	g := kk.Domaingraph
	// 将有向图导出为 Dot 格式的图形描述
	dot := "digraph G {\n"
	// 遍历节点
	for v := 0; v < g.Order(); v++ {
		// if v > Graph.Num {
		// 	break
		// }
		//fmt.Print(v, " -> ")
		aborted := graph.Sort(g).Visit(v, func(w int, c int64) (skip bool) {
			//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)

			switch kk.GraphReverse[w].Flag {
			case Graph.LeaveANode:
				str := fmt.Sprintf("%d", w)
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "blue" + "];\n"
				dot += str1
				return
			case Graph.LeaveAAAANode:
				str := fmt.Sprintf("%d", w)
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "blue" + "];\n"
				dot += str1
				return
			case Graph.LeaveCNAMENode:
				str := fmt.Sprintf("%d", w)
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "green" + "];\n"
				dot += str1
				return
			case Graph.NsNotGlueIPNode:
				str := fmt.Sprintf("%d", w)
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "yellow" + "];\n"
				dot += str1
				return
			}

			switch w {
			case Graph.RefusedNode:
				str := "Refused"
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "red" + "];\n"
				dot += str1
			case Graph.NameErrorNode:
				str := "NameError"
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "red" + "];\n"
				dot += str1
			case Graph.TimeoutNode:
				str := "Timeout"
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "red" + "];\n"
				dot += str1

			case Graph.CorruptNode:
				str := "Corrupt"
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "red" + "];\n"
				dot += str1

			case Graph.IPerrorNode:
				str := "IPerror"
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "red" + "];\n"
				dot += str1
			case Graph.NotImplementedNode:
				str := "NotImplemented"
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "red" + "];\n"
				dot += str1

			case Graph.IDMisMatchNode:
				str := "IDMisMatch"
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "red" + "];\n"
				dot += str1
			case Graph.NoNsrecordNode:
				str := "NoNsrecord"
				dot += fmt.Sprintf("\t%d -> \"%s\";\n", v, str)
				//node1 [label="Node 1", color="blue"];
				//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
				str1 := str + " [" + "color=" + "red" + "];\n"
				dot += str1

			default:
				dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
			}
			return
		})
		if aborted {
			break
		}
		//fmt.Println()
	}
	dot += "}"
	//fmt.Println(dot)
	// 创建一个文件，将 Dot 格式的图形描述写入该文件
	//fmt.Println("Hello")

	str := Path + domain + "directed_graph.dot"
	strr := Path + domain + "directed_graph.png"
	fmt.Println(str, strr)
	file, err := os.Create(str)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
	defer file.Close()

	if _, err := file.WriteString(dot); err != nil {
		log.Fatal(err)
	}

	fmt.Println("已生成 directed_graph.dot 文件")

	// 使用 Graphviz 将 Dot 文件渲染为图像
	cmd := exec.Command("dot", "-Tpng", str, "-o", strr)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("已生成 directed_graph.png 图片")
}

func Visual1(domain string, kk *Graph.GraphStruct) {
	fmt.Println("开始画图")
	// // // 创建一个有向图
	g := kk.Domaingraph
	// 将有向图导出为 Dot 格式的图形描述
	dot := "digraph G {\n"
	// 遍历节点
	for v := 0; v < g.Order(); v++ {
		// if v > Graph.Num {
		// 	break
		// }
		//fmt.Print(v, " -> ")
		aborted := graph.Sort(g).Visit(v, func(w int, c int64) (skip bool) {
			//dot += fmt.Sprintf("\t%d -> %d;\n", v, w)

			dot += fmt.Sprintf("\t \"%d\" -> \"%d\";\n", kk.GraphReverse[v], kk.GraphReverse[w])
			//fmt.Print(w, "  ")
			return
		})
		if aborted {
			break
		}
		//fmt.Println()
	}
	dot += "}"
	//fmt.Println(dot)
	// 创建一个文件，将 Dot 格式的图形描述写入该文件
	//fmt.Println("Hello")
	file, err := os.Create("Reverse.dot")
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
	defer file.Close()

	if _, err := file.WriteString(dot); err != nil {
		log.Fatal(err)
	}

	fmt.Println("已生成 directed_graph.dot 文件")

	// 使用 Graphviz 将 Dot 文件渲染为图像
	cmd := exec.Command("dot", "-Tpng", "Reverse.dot", "-o", "Reverse.png")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("已生成 directed_graph.png 图片")
}
