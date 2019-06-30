package ast

import (
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
	gast "github.com/yuin/goldmark/ast"
)

/*
* 第一步先渲染一个芯片
* $ U10-P8-NSTC12[1:VCC,8:GND] $
 */
type ICBlock struct {
	ICIndex  string //U10
	ICPins   int    //P8
	ICName   string // NSTC
	PinNames map[string]string
	Y        int
	X        int
}

// SchBlock 原理图
type SchBlock struct {
	gast.BaseBlock
	Ics []*ICBlock
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

// AddLine 添加一个行描述符
func (n *SchBlock) AddLine(desc string) int {
	log.Println(desc)
	// 拆出芯片编号->成功则添加一颗芯片
	ux := regexp.MustCompile(`U([0-9]+)-`)
	u := ux.FindStringSubmatch(desc)
	if len(u) > 1 {
		icb := new(ICBlock)
		icb.PinNames = make(map[string]string)
		n.Ics = append(n.Ics, icb)

		icb.ICIndex = u[1]
		// log.Println(n.ICIndex, u[1])
		desc = desc[len(u[0]):]
		//拆出引脚数量
		px := regexp.MustCompile(`P([0-9]+)-`)
		p := px.FindStringSubmatch(desc)
		if len(p) > 1 {
			icb.ICPins, _ = strconv.Atoi(p[1])
			desc = desc[len(p[0]):]
		}
		// 拆出芯片
		nx := regexp.MustCompile(`N([A-Za-z0-9]+)`)
		names := nx.FindStringSubmatch(desc)
		if len(names) > 1 {
			icb.ICName = names[1]
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
				icb.PinNames[apin[0]] = apin[1]
			}
		}
		// 拆出位置信息
		lc := regexp.MustCompile(`\(([0-9]+),([0-9]+)\)`)
		lsl := lc.FindStringSubmatch(desc)
		if len(lsl) >= 3 {
			icb.X, _ = strconv.Atoi(lsl[1])
			icb.Y, _ = strconv.Atoi(lsl[2])
			desc = desc[len(lsl[0]):]
		}
	}
	//
	return len(desc)
}

// ToSvg 输出svg
func (n *SchBlock) ToSvg(w io.Writer) {
	div := 12
	width := 500
	height := 500
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
		// 芯片编号
		canvas.Text(v.X*div, v.Y*div-div/2, "U"+v.ICIndex)
		// 画芯片引脚,先用左右方式
		// 左侧
		for i := 0; i < v.ICPins/2; i++ {
			canvas.Line(v.X*div, (v.Y+1+i)*div, (v.X-2)*div, (v.Y+1+i)*div, "stroke:#737375;")
		}
		// 右侧
		for i := 0; i < v.ICPins/2; i++ {
			canvas.Line(v.X*div+10*div, (v.Y+1+i)*div, (v.X+12)*div, (v.Y+1+i)*div, "stroke:#737375;")
		}

		// 画中间框
		canvas.Rect(v.X*div, v.Y*div, 10*div, v.ICPins*div/2+1*div, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		// 芯片名称
		canvas.Text(v.X*div, v.Y*div+v.ICPins*div/2+2*div, v.ICName, "font-size:"+strconv.Itoa(div)+"px;")

	}

	// canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	canvas.End()
}

// NewSchBlock 解析出一个新芯片
func NewSchBlock() *SchBlock {
	return &SchBlock{}
}
