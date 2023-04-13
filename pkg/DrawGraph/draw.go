package DrawGraph

import (
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
}

func Visual() {
	fmt.Println("开始画图")
	// // // 创建一个有向图
	g := Graph.Domaingraph
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
			dot += fmt.Sprintf("\t%d -> %d;\n", v, w)
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
	file, err := os.Create("directed_graph.dot")
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
	cmd := exec.Command("dot", "-Tpng", "directed_graph.dot", "-o", "directed_graph.png")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("已生成 directed_graph.png 图片")
}
