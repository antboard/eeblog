package ast

import (
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

// 电容
// C10-V40n-H(1,2)-P
// 编号10 阻值40k, 横向放置, 有极性
// 起始坐标点为1号引脚, 向右或者向下画电阻框和2号引脚, 极性电容标记+, 和弯符号

func init() {
	schParsers = append(schParsers, new(CBlock))
}

// CBlock 电阻块
type CBlock struct {
	Index    string
	Value    string
	Layout   string
	Polarity bool
	X        int
	Y        int
}

// CanParse 类型检查
func (cb *CBlock) CanParse(desc string) bool {
	// 如果有Cn出现就是电阻
	cx := regexp.MustCompile(`C([0-9]+)-`)
	c := cx.FindStringSubmatch(desc)
	if len(c) > 1 {
		log.Println("parse C success ", desc)
		return true
	}
	return false
}

// ParseLine 解析块
func (cb *CBlock) ParseLine(desc string) SvgBlock {
	cx := regexp.MustCompile(`C([0-9]+)-`)
	c := cx.FindStringSubmatch(desc)
	if len(c) > 1 {
		cur := new(CBlock)
		cur.Index = c[1]
		desc = desc[len(c[0]):]
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
			desc = desc[len(lc[0]):]
		}
		// 因为有回车,这里用前缀.
		cur.Polarity = strings.HasPrefix(desc, "-P")
		return cur
	}
	return nil
}

// ToSvg 生成svg
func (cb *CBlock) ToSvg(canvas *svg.SVG, w io.Writer) {
	log.Printf("%#v\n", cb)
	if cb.Layout == "H" {
		// 1号腿
		canvas.Line(cb.X*div, cb.Y*div, (cb.X+1)*div, cb.Y*div, "stroke:#737375;")
		// 2号腿
		canvas.Line((cb.X+2)*div, cb.Y*div, (cb.X+3)*div, cb.Y*div, "stroke:#737375;")
		// 竖线1
		canvas.Line((cb.X+1)*div, cb.Y*div-div, (cb.X+1)*div, cb.Y*div+div, "stroke:#737375;")
		if !cb.Polarity {
			canvas.Line((cb.X+2)*div, cb.Y*div-div, (cb.X+2)*div, cb.Y*div+div, "stroke:#737375;")
		} else {
			canvas.Arc((cb.X+2)*div, cb.Y*div-div, div, div, div, false, false, (cb.X+2)*div, cb.Y*div+div, "stroke:#737375;fill:#cdcdcf;")
		}
		// 编号
		canvas.Text((cb.X+1)*div, cb.Y*div-div/2*3, "R"+cb.Index, "font-size:"+strconv.Itoa(div)+"px;")
		// 阻值
		canvas.Text((cb.X+1)*div, cb.Y*div+div*2, cb.Index, "font-size:"+strconv.Itoa(div)+"px;")
	} else {
		// 1号腿
		canvas.Line(cb.X*div, cb.Y*div, (cb.X)*div, (cb.Y+1)*div, "stroke:#737375;")
		// 2号腿
		canvas.Line((cb.X)*div, (cb.Y+2)*div, (cb.X)*div, (cb.Y+3)*div, "stroke:#737375;")
		// 横线1
		canvas.Line((cb.X-1)*div, (cb.Y+1)*div, (cb.X+1)*div, (cb.Y+1)*div, "stroke:#737375;")
		if !cb.Polarity {
			canvas.Line((cb.X-1)*div, (cb.Y+2)*div, (cb.X+1)*div, (cb.Y+2)*div, "stroke:#737375;")
		} else {
			canvas.Arc((cb.X-1)*div, (cb.Y+2)*div, div, div, div, true, true, (cb.X+1)*div, (cb.Y+2)*div, "stroke:#737375;fill:#cdcdcf;")
		}
		// 编号
		canvas.Text((cb.X)*div+div/2*3, (cb.Y+1)*div+div/2, "R"+cb.Index, "font-size:"+strconv.Itoa(div)+"px;")
		// 阻值
		canvas.Text((cb.X)*div+div/2*3, (cb.Y+2)*div+div/2, cb.Index, "font-size:"+strconv.Itoa(div)+"px;")

	}
}
