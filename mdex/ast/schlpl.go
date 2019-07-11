package ast

import (
	"io"
	"log"
	"regexp"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

// 电感
// L10-V40n-H(1,2)
// 编号10 阻值40k, 横向放置
// 起始坐标点为1号引脚, 向右或者向下画电阻框和2号引脚

func init() {
	schParsers = append(schParsers, new(LBlock))
}

// LBlock 电阻块
type LBlock struct {
	Index  string
	Value  string
	Layout string
	X      int
	Y      int
}

// CanParse 类型检查
func (lb *LBlock) CanParse(desc string) bool {
	// 如果有Ln出现就是电阻
	lx := regexp.MustCompile(`L([0-9]+)-`)
	l := lx.FindStringSubmatch(desc)
	if len(l) > 1 {
		log.Println("parse L success ", desc)
		return true
	}
	return false
}

// ParseLine 解析块
func (lb *LBlock) ParseLine(desc string) SvgBlock {
	lx := regexp.MustCompile(`L([0-9]+)-`)
	l := lx.FindStringSubmatch(desc)
	if len(l) > 1 {
		cur := new(LBlock)
		cur.Index = l[1]
		desc = desc[len(l[0]):]
		vx := regexp.MustCompile(`V([0-9a-zA-Z]+)-`)
		v := vx.FindStringSubmatch(desc)
		if len(v) > 1 {
			cur.Value = v[1]
			desc = desc[len(v[0]):]
		}
		lcx := regexp.MustCompile(`([HV])\(([0-9]+),([0-9]+)\)`)
		lc := lcx.FindStringSubmatch(desc)
		if len(lc) > 3 {
			cur.Layout = lc[1]
			cur.X, _ = strconv.Atoi(lc[2])
			cur.Y, _ = strconv.Atoi(lc[3])
		}
		return cur
	}
	return nil
}

// ToSvg 生成svg 这里是拼的图形.不是教科书级设计了.path以后研究
func (lb *LBlock) ToSvg(canvas *svg.SVG, w io.Writer) {
	// log.Printf("%#v\n", rb)
	// 暂时不用旋转,直接画
	if lb.Layout == "H" {
		// 1号腿
		canvas.Line(lb.X*div, lb.Y*div, (lb.X+1)*div, lb.Y*div, "stroke:#737375;")
		// 2号腿
		canvas.Line((lb.X+4)*div, lb.Y*div, (lb.X+5)*div, lb.Y*div, "stroke:#737375;")
		// 波浪
		canvas.Circle((lb.X+1)*div+div/2, lb.Y*div, div/2, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		canvas.Circle((lb.X+2)*div+div/2, lb.Y*div, div/2, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		canvas.Circle((lb.X+3)*div+div/2, lb.Y*div, div/2, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		canvas.Rect((lb.X+1)*div, lb.Y*div+1, 3*div+5, div/2+3, "fill:#e0e0e2;stroke:#737375;stroke-width:0pt;")
		// 编号
		canvas.Text((lb.X+1)*div, lb.Y*div-div, "R"+lb.Index, "font-size:"+strconv.Itoa(div)+"px;")
		// 阻值
		canvas.Text((lb.X+1)*div, lb.Y*div+div/2*3, lb.Index, "font-size:"+strconv.Itoa(div)+"px;")
	} else {
		// 1号腿
		canvas.Line(lb.X*div, lb.Y*div, (lb.X)*div, (lb.Y+1)*div, "stroke:#737375;")
		// 2号腿
		canvas.Line((lb.X)*div, (lb.Y+4)*div, (lb.X)*div, (lb.Y+5)*div, "stroke:#737375;")
		// 波浪
		canvas.Circle(lb.X*div, (lb.Y+1)*div+div/2, div/2, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		canvas.Circle(lb.X*div, (lb.Y+2)*div+div/2, div/2, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		canvas.Circle(lb.X*div, (lb.Y+3)*div+div/2, div/2, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		canvas.Rect((lb.X)*div-div/2-3, (lb.Y+1)*div, div/2+2, 3*div+5, "fill:#e0e0e2;stroke:#737375;stroke-width:0pt;")
		// 编号
		canvas.Text((lb.X)*div+div, (lb.Y+2)*div, "R"+lb.Index, "font-size:"+strconv.Itoa(div)+"px;")
		// 阻值
		canvas.Text((lb.X)*div+div, (lb.Y+3)*div, lb.Index, "font-size:"+strconv.Itoa(div)+"px;")

	}
}

// GetName 获取元件名
func (lb *LBlock) GetName() string {
	return "L" + lb.Index
}

// GetPin 获取引脚位置
func (lb *LBlock) GetPin(i int) (x int, y int) {
	if i == 1 {
		return lb.X, lb.Y
	}
	if lb.Layout == "H" {
		return lb.X + div*3, lb.Y
	}
	return lb.X, lb.Y + div*3
}
