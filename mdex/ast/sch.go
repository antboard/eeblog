package ast

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

/*
* 第一步先渲染一个芯片
* $ U10-P8-NSTC12[1:VCC,8:GND] $
 */

// SchBlock 原理图
type SchBlock struct {
	gast.BaseBlock
	ICIndex  string //U10
	ICPins   int    //P8
	ICName   string // NSTC
	PinNames map[string]string
	Y        int
	X        int
}

// Dump 继承
func (n *SchBlock) Dump(source []byte, level int) {
	m := make(map[string]string)
	m["index"] = n.ICIndex
	m["pins"] = strconv.Itoa(n.ICPins)
	m["name"] = n.ICName
	gast.DumpHelper(n, source, level, m, nil)
}

// KindSchBlock 原理图描述类
var KindSchBlock = gast.NewNodeKind("SchBlock")

// Kind implements Node.Kind.
func (n *SchBlock) Kind() gast.NodeKind {
	return KindSchBlock
}

// AddLine 添加一个行描述符
func (n *SchBlock) AddLine(desc string) int {
	log.Println(desc)
	// 拆出芯片编号
	ux := regexp.MustCompile(`U([0-9]+)-`)
	u := ux.FindStringSubmatch(desc)
	if len(u) > 1 {
		n.ICIndex = u[1]
		// log.Println(n.ICIndex, u[1])
		desc = desc[len(u[0]):]
	}
	//拆出引脚数量
	px := regexp.MustCompile(`P([0-9]+)-`)
	p := px.FindStringSubmatch(desc)
	if len(p) > 1 {
		n.ICPins, _ = strconv.Atoi(p[1])
		desc = desc[len(p[0]):]
	}
	// 拆出芯片
	nx := regexp.MustCompile(`N([A-Za-z0-9]+)`)
	names := nx.FindStringSubmatch(desc)
	if len(names) > 1 {
		n.ICName = names[1]
		desc = desc[len(names[0]):]
	}
	// 如果有[]则解析引脚命名
	nstart := strings.Index(desc, "[")
	if nstart >= 0 {
		nend := strings.Index(desc, "]")
		pinstr := desc[nstart+1 : nend]
		desc = desc[nend+1:]
		pins := strings.Split(pinstr, ",")
		for _, v := range pins {
			apin := strings.Split(v, ":")
			n.PinNames[apin[0]] = apin[1]
		}
	}
	// 拆出位置信息
	lc := regexp.MustCompile(`\(([0-9]+),([0-9]+)\)`)
	lsl := lc.FindStringSubmatch(desc)
	if len(lsl) >= 3 {
		n.X, _ = strconv.Atoi(lsl[1])
		n.Y, _ = strconv.Atoi(lsl[2])
		desc = desc[len(lsl[0]):]
	}
	return len(desc)
}

// NewSchBlock 解析出一个新芯片
func NewSchBlock() *SchBlock {
	return &SchBlock{PinNames: make(map[string]string)}
}
