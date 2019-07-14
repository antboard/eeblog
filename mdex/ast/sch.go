package ast

import (
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strconv"

	svg "github.com/ajstarks/svgo"
	gast "github.com/yuin/goldmark/ast"
)

const (
	div = 12 // 1div = 12 pt
)

// 保存解析器
var schParsers = make([]Lineparser, 0, 10)

// SchBlock 原理图
type SchBlock struct {
	gast.BaseBlock
	Ics   []SvgBlock
	PageW int
	PageH int
}

// Dump 继承
func (n *SchBlock) Dump(source []byte, level int) {
	m := make(map[string]string)
	bic, _ := json.Marshal(n.Ics)
	m["ics"] = string(bic)
	gast.DumpHelper(n, source, level, m, nil)
}

// KindSchBlock 原理图描述类
var KindSchBlock = gast.NewNodeKind("SchBlock")

// Kind implements Node.Kind.
func (n *SchBlock) Kind() gast.NodeKind {
	return KindSchBlock
}

// InitByLine 初始化画布
func (n *SchBlock) InitByLine(desc string) {
	// 读取画布大小
	pageszre := regexp.MustCompile(`\$\(([0-9]+),([0-9]+)\)`)
	pagesz := pageszre.FindStringSubmatch(desc)
	if len(pagesz) == 3 {
		n.PageW, _ = strconv.Atoi(pagesz[1])
		n.PageH, _ = strconv.Atoi(pagesz[2])
		return
	}
}

// GetIcByIndex 根据名字找到ic
func (n *SchBlock) GetIcByIndex(idx string) SvgBlock {
	for _, ic := range n.Ics {
		if ic.GetIdxName() == idx {
			return ic
		}
	}
	return new(DumySvgBlock)
}

// AddLine 添加一个行描述符
func (n *SchBlock) AddLine(desc string) int {
	log.Println(desc)

	for _, v := range schParsers {
		if v.CanParse(desc) {
			lp := v.ParseLine(n, desc)
			if lp != nil {
				n.Ics = append(n.Ics, lp)
			}
			break
		}
	}
	return len(desc)
}

// ToSvg 输出svg
func (n *SchBlock) ToSvg(w io.Writer) {
	width := n.PageW * div
	height := n.PageH * div
	canvas := svg.New(w)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:#e0e0e2;")
	for i := div; i < width; i += div {
		for j := div; j < height; j += div {
			if ((i/div)%5 == 0) || ((j/div)%5 == 0) {
				canvas.Circle(i, j, 1) // fill: gray;
			} else {
				canvas.Circle(i, j, 1, "stroke:#e0e0e2;") // fill: gray;
			}
		}
	}

	for _, v := range n.Ics {
		v.ToSvg(canvas, w)
	}

	canvas.End()
}

// NewSchBlock 解析出一个新芯片
func NewSchBlock() *SchBlock {
	return &SchBlock{PageW: 50, PageH: 50}
}
