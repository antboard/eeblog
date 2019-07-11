package ast

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

func init() {
	schParsers = append(schParsers, new(ICBlock))
}

/*
* 第一步先渲染一个芯片
* $ U10-P8-NSTC12[1:VCC,8:GND](3,2,8) $
 */
type point struct {
	X, Y int
}

// ICBlock 芯片描述
type ICBlock struct {
	Index    string //U10
	Pins     int    //P8
	Name     string // NSTC
	PinNames map[string]string
	Y        int
	X        int
	W        int
	PinLc    []*point
}

// CanParse 可解析判断
func (ic *ICBlock) CanParse(desc string) bool {
	ux := regexp.MustCompile(`U([0-9]+)-`)
	u := ux.FindStringSubmatch(desc)
	if len(u) > 1 {
		return true
	}
	return false
}

// ParseLine 根据行解析出目标块
func (ic *ICBlock) ParseLine(desc string) SvgBlock {
	ux := regexp.MustCompile(`U([0-9]+)-`)
	u := ux.FindStringSubmatch(desc)
	if len(u) > 1 {
		icb := new(ICBlock)
		icb.PinNames = make(map[string]string)

		icb.Index = u[1]
		// log.Println(n.Index, u[1])
		desc = desc[len(u[0]):]
		//拆出引脚数量
		px := regexp.MustCompile(`P([0-9]+)-`)
		p := px.FindStringSubmatch(desc)
		if len(p) > 1 {
			icb.Pins, _ = strconv.Atoi(p[1])
			icb.PinLc = make([]*point, icb.Pins, icb.Pins)
			desc = desc[len(p[0]):]
		}
		// 拆出芯片
		nx := regexp.MustCompile(`N([A-Za-z0-9]+)`)
		names := nx.FindStringSubmatch(desc)
		if len(names) > 1 {
			icb.Name = names[1]
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
		lc := regexp.MustCompile(`\(([0-9]+),([0-9]+),([0-9]+)\)`)
		lsl := lc.FindStringSubmatch(desc)
		if len(lsl) > 3 {
			icb.X, _ = strconv.Atoi(lsl[1])
			icb.Y, _ = strconv.Atoi(lsl[2])
			icb.W, _ = strconv.Atoi(lsl[3])
			desc = desc[len(lsl[0]):]
		}
		return icb
	}
	return nil
}

// ToSvg 转成图形
func (ic *ICBlock) ToSvg(canvas *svg.SVG, w io.Writer) {
	// 画中间框
	canvas.Rect(ic.X*div, ic.Y*div, ic.W*div, ic.Pins*div/2+1*div, "fill:#cdcdcf;stroke:#737375;stroke-width:1pt;")

	// 芯片编号
	canvas.Text(ic.X*div, ic.Y*div-div/2, "U"+ic.Index)
	// 画芯片引脚,先用左右方式
	// 左侧
	for i := 0; i < ic.Pins/2; i++ {
		// 引脚编号
		strpin := strconv.Itoa(i + 1)
		canvas.Text((ic.X-2)*div, (ic.Y+1+i)*div, strpin, "font-size:"+strconv.Itoa(div)+"px;")
		// 引脚名称
		name, ok := ic.PinNames[strpin]
		if ok {
			canvas.Text((ic.X)*div+2, (ic.Y+1+i)*div+div/3, name, "font-size:"+strconv.Itoa(div)+"px;")
		}
		// 引脚线
		canvas.Line(ic.X*div, (ic.Y+1+i)*div, (ic.X-2)*div, (ic.Y+1+i)*div, "stroke:#737375;")
		pt := &point{
			X: (ic.X - 2) * div,
			Y: (ic.Y + 1 + i) * div,
		}
		ic.PinLc[i] = pt
	}
	// 右侧
	for i := 0; i < ic.Pins/2; i++ {
		// U型编号
		strpin := strconv.Itoa(ic.Pins - i)
		canvas.Text((ic.X+ic.W+1)*div, (ic.Y+1+i)*div, strpin, "font-size:"+strconv.Itoa(div)+"px;")
		// 引脚名称
		name, ok := ic.PinNames[strpin]
		if ok {
			canvas.Text((ic.X+ic.W)*div, (ic.Y+1+i)*div+div/3, name, "font-size:"+strconv.Itoa(div)+"px;text-anchor: end")
		}
		// 引脚线
		canvas.Line(ic.X*div+ic.W*div, (ic.Y+1+i)*div, (ic.X+ic.W+2)*div, (ic.Y+1+i)*div, "stroke:#737375;")
		pt := &point{
			X: (ic.X + ic.W + 2) * div,
			Y: (ic.Y + 1 + i) * div,
		}
		ic.PinLc[ic.Pins-i-1] = pt
	}

	// 芯片名称
	canvas.Text(ic.X*div, ic.Y*div+ic.Pins*div/2+2*div, ic.Name, "font-size:"+strconv.Itoa(div)+"px;")

}

// GetName 获取元件编号
func (ic *ICBlock) GetName() string {
	return "U" + ic.Index
}

// GetPin 获取引脚位置
func (ic *ICBlock) GetPin(i int) (x int, y int) {
	return ic.PinLc[i-1].X, ic.PinLc[i-1].Y
}
