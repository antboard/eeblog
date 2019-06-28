package ast

import (
	"strconv"

	gast "github.com/yuin/goldmark/ast"
)

/*
* 第一步先渲染一个芯片
* $ U10-P8-NSTC12[1:VCC,8:GND] $
 */

// SchBlock 原理图
type SchBlock struct {
	gast.BaseBlock
	ICIndex  string
	ICPins   int
	ICName   string
	PinNames map[string]string
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

// NewSchBlock 解析出一个新芯片
func NewSchBlock() *SchBlock {
	return &SchBlock{}
}
