package ast

import (
	"io"
	"log"
	"regexp"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

// 电阻
// R10-V40k-H(1,2)
// 编号10 阻值40k, 横向放置
// 起始坐标点为1号引脚, 向右或者向下画电阻框和2号引脚

func init() {
	schParsers = append(schParsers, new(RBlock))
}

// RBlock 电阻块
type RBlock struct {
	Index  string
	Value  string
	Layout string
	X      int
	Y      int
}

// CanParse 类型检查
func (rb *RBlock) CanParse(desc string) bool {
	// 如果有Rn出现就是电阻
	rx := regexp.MustCompile(`R([0-9]+)-`)
	r := rx.FindStringSubmatch(desc)
	if len(r) > 1 {
		log.Println("parse R success ", desc)
		return true
	}

	return false
}

// ParseLine 解析块
func (rb *RBlock) ParseLine(desc string) SvgBlock {
	rx := regexp.MustCompile(`R([0-9]+)-`)
	r := rx.FindStringSubmatch(desc)
	if len(r) > 1 {
		cur := new(RBlock)
		cur.Index = r[1]
		desc = desc[len(r[0]):]
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

// ToSvg 生成svg
func (rb *RBlock) ToSvg(canvas *svg.SVG, w io.Writer) {
	// log.Printf("%#v\n", rb)
	// 暂时不用旋转,直接画
	if rb.Layout == "H" {
		// 1号腿
		canvas.Line(rb.X*div, rb.Y*div, (rb.X+1)*div, rb.Y*div, "stroke:#737375;")
		// 2号腿
		canvas.Line((rb.X+4)*div, rb.Y*div, (rb.X+5)*div, rb.Y*div, "stroke:#737375;")
		// 方框
		canvas.Rect((rb.X+1)*div, rb.Y*div-div/2, 3*div, 1*div, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		// 编号
		canvas.Text((rb.X+1)*div, rb.Y*div-div, "R"+rb.Index, "font-size:"+strconv.Itoa(div)+"px;")
		// 阻值
		canvas.Text((rb.X+1)*div, rb.Y*div+div/2*3, rb.Index, "font-size:"+strconv.Itoa(div)+"px;")
	} else {
		// 1号腿
		canvas.Line(rb.X*div, rb.Y*div, (rb.X)*div, (rb.Y+1)*div, "stroke:#737375;")
		// 2号腿
		canvas.Line((rb.X)*div, (rb.Y+4)*div, (rb.X)*div, (rb.Y+5)*div, "stroke:#737375;")
		// 方框
		canvas.Rect((rb.X)*div-div/2, (rb.Y+1)*div, div, 3*div, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")
		// 编号
		canvas.Text((rb.X)*div+div, (rb.Y+2)*div, "R"+rb.Index, "font-size:"+strconv.Itoa(div)+"px;")
		// 阻值
		canvas.Text((rb.X)*div+div, (rb.Y+3)*div, rb.Index, "font-size:"+strconv.Itoa(div)+"px;")

	}
}

// GetName 获取元件名
func (rb *RBlock) GetName() string {
	return "R" + rb.Index
}

// GetPin 获取引脚位置
func (rb *RBlock) GetPin(i int) (x int, y int) {
	if i == 1 {
		return rb.X, rb.Y
	}
	if rb.Layout == "H" {
		return rb.X + div*3, rb.Y
	}
	return rb.X, rb.Y + div*3
}
