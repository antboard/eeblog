package ast

import (
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

// 自定义芯片
// 封闭点线框Point
// 线段Line
// 引脚Pin
// 椭圆Cir

func init() {
	schParsers = append(schParsers, new(UserDefine))
}

type line struct {
	Start *point
	End   *point
}

type cir struct {
	x1, y1, d int
}

type pinCfg struct {
	X, Y    int
	Dir     string // 方向
	Dist    int    // 距离
	PinIdx  int    // 引脚索引
	PinName string // 引脚名称
}

// UserDefine 自定义模块
type UserDefine struct {
	Index  string
	Points []*point
	Cir    *cir
	Lines  []*line
	Pins   []*pinCfg
	Lable  point // 原件索引显示位置
	Value  point // 原件数值显示位置
}

// CanParse 类型检查
func (ud *UserDefine) CanParse(desc string) bool {
	eix := regexp.MustCompile(`^[\s]*Ei([0-9]+)-`)
	ei := eix.FindStringSubmatch(desc)
	if len(ei) > 1 {
		log.Println("parse Ei success ", desc)
		return true
	}
	return false
}

// ParseLine 解析块
func (ud *UserDefine) ParseLine(b *SchBlock, desc string) SvgBlock {
	cur := new(UserDefine)
	// Ei35-
	eix := regexp.MustCompile(`^[\s]*Ei([0-9]+)-`)
	ei := eix.FindStringSubmatch(desc)
	if len(ei) > 1 {
		cur.Index = ei[1]
		desc = desc[len(ei[0]):]
	}
	for {
		desc = strings.TrimPrefix(desc, "-")
		// log.Println(desc)
		// Point[(1,1);(1,5);(5,1);(5,5)]
		if strings.HasPrefix(desc, "Points[") {
			pos := strings.Index(desc, "]")
			ptsStr := desc[len("Point["):pos]
			pts := strings.Split(ptsStr, ";")
			for _, ptstr := range pts {
				ptx := regexp.MustCompile(`\(([0-9]+),([0-9]+)\)`)
				lpt := ptx.FindStringSubmatch(ptstr)
				if len(lpt) > 1 {
					pt := new(point)
					pt.X, _ = strconv.Atoi(lpt[1])
					pt.Y, _ = strconv.Atoi(lpt[2])
					cur.Points = append(cur.Points, pt)
				}
			}
			desc = desc[pos+1:]
			continue
		}
		// Line[(0,1)-(1,1);(0,2)-(1,2);(0,5)-(1,5)]-
		if strings.HasPrefix(desc, "Lines[") {
			pos := strings.Index(desc, "]")
			linesstr := desc[len("Line["):pos]
			lines := strings.Split(linesstr, ";")
			for _, linestr := range lines {
				linex := regexp.MustCompile(`\(([0-9]+),([0-9]+)\)-\(([0-9]+),([0-9]+)\)`)
				lpts := linex.FindStringSubmatch(linestr)
				if len(lpts) > 1 {
					start := new(point)
					start.X, _ = strconv.Atoi(lpts[1])
					start.Y, _ = strconv.Atoi(lpts[2])
					end := new(point)
					end.X, _ = strconv.Atoi(lpts[3])
					end.Y, _ = strconv.Atoi(lpts[4])
					line := new(line)
					line.Start = start
					line.End = end
					cur.Lines = append(cur.Lines, line)
				}
			}
			desc = desc[pos+1:]
			continue
		}
		//Cir(3,3,1)- 暂时支持圆.预留椭圆
		cirx := regexp.MustCompile(`^Cir\(([0-9]+),([0-9]+),([0-9]+),([0-9]*)\)`)
		cirs := cirx.FindStringSubmatch(desc)
		if len(cirs) > 1 {
			c := new(cir)
			c.x1, _ = strconv.Atoi(cirs[1])
			c.y1, _ = strconv.Atoi(cirs[2])
			c.d, _ = strconv.Atoi(cirs[3])
			// c.ry, _ = strconv.Atoi(cirs[4])
			cur.Cir = c
			desc = desc[len(cirs[0]):]
			continue
		}
		//Pin[(0,1,+1,12,vcc);(0,2,+1,13,gnd)]
		if strings.HasPrefix(desc, "Pins[") {
			pos := strings.Index(desc, "]")
			pinStr := desc[len("Pins["):pos]
			pins := strings.Split(pinStr, ";")
			for _, pinstr := range pins {
				pinx := regexp.MustCompile(`\(([0-9]+),([0-9]+),([\+\-/\*])([0-9]+),([0-9]*),([a-zA-Z0-9]*)\)`)
				pins := pinx.FindStringSubmatch(pinstr)
				if len(pins) > 1 {
					pt := new(pinCfg)
					pt.X, _ = strconv.Atoi(pins[1])
					pt.Y, _ = strconv.Atoi(pins[2])
					pt.Dir = pins[3]
					pt.Dist, _ = strconv.Atoi(pins[4])
					pt.PinIdx, _ = strconv.Atoi(pins[5])
					pt.PinName = pins[6]
					cur.Pins = append(cur.Pins, pt)
				}
			}
			desc = desc[pos+1:]
			continue
		}
		// Label(3,3)
		labelx := regexp.MustCompile(`^Label\(([0-9]+),([0-9]+)\)`)
		labels := labelx.FindStringSubmatch(desc)
		if len(labels) > 1 {
			cur.Lable.X, _ = strconv.Atoi(labels[1])
			cur.Lable.Y, _ = strconv.Atoi(labels[2])
			desc = desc[len(labels[0]):]
			continue
		}
		// Value(5,3)
		valuex := regexp.MustCompile(`^Value\(([0-9]+),([0-9]+)\)`)
		values := valuex.FindStringSubmatch(desc)
		if len(values) > 1 {
			cur.Value.X, _ = strconv.Atoi(values[1])
			cur.Value.Y, _ = strconv.Atoi(values[2])
			desc = desc[len(values[0]):]
			continue
		}
		break
	}

	return cur

}

// ToSvg 生成svg
func (ud *UserDefine) ToSvg(canvas *svg.SVG, w io.Writer) {
}

// LayoutSvg 生成svg
func (ud *UserDefine) LayoutSvg(canvas *svg.SVG, w io.Writer, offsetX, offsetY int, label string, val string) {
	if len(ud.Points) > 0 {
		str := "M" + strconv.Itoa((ud.Points[0].X+offsetX)*div) + " " + strconv.Itoa((ud.Points[0].Y+offsetY)*div)
		for i := 1; i < len(ud.Points); i++ {
			str += " L" + strconv.Itoa((ud.Points[i].X+offsetX)*div) + " " + strconv.Itoa((ud.Points[i].Y+offsetY)*div)
		}
		str += " Z"
		canvas.Path(str, `style="fill:#cdcdcf;stroke:#737375;stroke-width:1pt;"`)
	}

	if ud.Cir != nil {
		canvas.Circle((ud.Cir.x1+offsetX)*div, (ud.Cir.y1+offsetY)*div, ud.Cir.d*div/2, `style="fill:#cdcdcf;stroke:#737375;stroke-width:1pt;"`)
	}
	for _, vline := range ud.Lines {
		canvas.Line(vline.Start.X*div/3+offsetX*div, vline.Start.Y*div/3+offsetY*div, vline.End.X*div/3+offsetX*div, vline.End.Y*div/3+offsetY*div, "stroke:#737375;")
	}
	for _, pt := range ud.Pins {
		canvas.Circle((pt.X+offsetX)*div, (pt.Y+offsetY)*div, 2, "fill:#e0e0e2;stroke:#737375;")
		destX, destY := pt.X, pt.Y
		ex := ""
		fixX := 0
		fixY := 0
		rot := false
		switch pt.Dir {
		case "+":
			destX += pt.Dist
			fixY = div / 3
		case "-":
			destX -= pt.Dist
			ex = "text-anchor: end;"
			fixY = div / 3
		case "*":
			destY -= pt.Dist
			ex = "text-anchor: end;"
			fixY = div / 3
			rot = true
		case "/":
			destY += pt.Dist
			fixY = div / 3
			rot = true
		}
		canvas.Line((pt.X+offsetX)*div, (pt.Y+offsetY)*div, (destX+offsetX)*div, (destY+offsetY)*div, "stroke:#737375;")
		txtX, txtY := (destX+offsetX)*div+fixX, (destY+offsetY)*div+fixY
		strrot := ""
		if rot {
			strrot = `transform="rotate(90,` + strconv.Itoa(txtX) + " " + strconv.Itoa(txtY) + `)"`
		}
		canvas.Text(txtX, txtY, pt.PinName, "font-size:"+strconv.Itoa(div)+"px;"+ex, strrot)
	}
	if label != "" {
		canvas.Text((ud.Lable.X+offsetX)*div, (ud.Lable.Y+offsetY)*div, label, "text-anchor: middle;")
	}
	if val != "" {
		canvas.Text((ud.Value.X+offsetX)*div, (ud.Value.Y+offsetY)*div, label, "text-anchor: middle;")
	}
}

// GetIdxName 获取元件名
func (ud *UserDefine) GetIdxName() string {
	return "Ei" + ud.Index
}

// GetPin 获取引脚位置
func (ud *UserDefine) GetPin(i int) (x int, y int) {
	for _, pin := range ud.Pins {
		if pin.PinIdx == i {
			return pin.X, pin.Y
		}
	}
	return 0, 0
}
