package ast

import (
	"io"
	"log"
	"regexp"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

/*
* 连线是最不好写的.一直在犹豫
* 起点,还是选可连线的锚点,比如U1.8 是指芯片U1的pin8
* 拐点,有两种,一种是相对距离, 另一种是绝对点位置.+\-* /对应左右上下,
* 因为加减是x轴的常用方式, y轴为了好分辨,用* /
* 相对位置因为有方向,所以直接写数字, 比如+5, 正向移动5div, sch图不提供更精细的连线.
* 绝对位置不在乎符号,但是应该是加减乘除的一种, 后面跟绝对坐标(x, y),
* 当前版本只支持横竖线, 所以允许(5,)这种简写, 就是移动到x=5, y不变的一条横线.
 */

func init() {
	schParsers = append(schParsers, new(NetBlock))
}

type lineSegment struct {
	Segment string // 原始字符串描述
	Lc      point  // 转换成路径点,容易画图
}

// NetBlock 线网块
type NetBlock struct {
	Lp []*lineSegment
}

// CreatNetBlock 创建一个线段块
func CreatNetBlock() *NetBlock {
	return &NetBlock{make([]*lineSegment, 0, 100)}
}

// CanParse 类型检查
func (nb *NetBlock) CanParse(desc string) bool {
	// 如果有Rn出现就是电阻
	nx := regexp.MustCompile(`^[\s]*N([a-zA-z0-9]+)\.?([0-9]+)?`)
	n := nx.FindStringSubmatch(desc)
	if len(n) > 1 {
		log.Println("parse Net success ", desc)
		return true
	}

	return false
}

// ParseLine 解析块
func (nb *NetBlock) ParseLine(b *SchBlock, desc string) SvgBlock {
	nx := regexp.MustCompile(`^[\s]*N([a-zA-z0-9]+)\.?([0-9]+)?`)
	n := nx.FindStringSubmatch(desc)
	lastPt := point{}
	if len(n) > 1 {
		cur := CreatNetBlock()
		lp := new(lineSegment)
		lp.Segment = n[0]
		partIdxName := n[1] // 原件名
		partPinStr := n[2]  // 引脚名
		log.Println(partIdxName, "--", partPinStr)
		pinIdx, _ := strconv.Atoi(partPinStr)
		lp.Lc.X, lp.Lc.Y = b.GetIcByIndex(partIdxName).GetPin(pinIdx)
		lastPt = lp.Lc
		cur.Lp = append(cur.Lp, lp)
		desc = desc[len(n[0]):]
		// 解析中间点
		for {
			// 正则分析->相对路径
			{
				nx := regexp.MustCompile(`^[\s]*([\+\-/\*])([0-9]+)`)
				n := nx.FindStringSubmatch(desc)
				if len(n) > 1 {
					lp := new(lineSegment)
					lp.Segment = n[0]
					lp.Lc = lastPt
					oper := n[1]
					sz, _ := strconv.Atoi(n[2])
					switch oper {
					case "+":
						lp.Lc.X += sz
					case "-":
						lp.Lc.X -= sz
					case "*":
						lp.Lc.Y -= sz
					case "/":
						lp.Lc.Y += sz
					}
					lastPt = lp.Lc
					cur.Lp = append(cur.Lp, lp)
					desc = desc[len(n[0]):]
					continue
				}
			}
			// 正则分析->绝对路径
			{
				nx := regexp.MustCompile(`^[\s]*([\+\-/\*])\(([0-9]+)?,([0-9]+)?\)`)
				n := nx.FindStringSubmatch(desc)
				if len(n) > 1 {
					lp := new(lineSegment)
					lp.Segment = n[0]
					lp.Lc = lastPt
					// oper := n[1]
					szx, _ := strconv.Atoi(n[2])
					szy, _ := strconv.Atoi(n[3])
					if szx != 0 {
						lp.Lc.X = szx
					}
					if szy != 0 {
						lp.Lc.Y = szy
					}
					lastPt = lp.Lc
					cur.Lp = append(cur.Lp, lp)
					desc = desc[len(n[0]):]
					continue
				}
			}
			// 出错退出->目标位置
			break
		}
		// 解析目标点
		log.Println(desc)
		nx := regexp.MustCompile(`^[\s]*([\+\-/\*])([a-zA-z0-9]+)\.?([0-9]+)?`)
		n := nx.FindStringSubmatch(desc)
		if len(n) > 1 {
			lp := new(lineSegment)
			lp.Segment = n[0]
			partIdxName := n[2] // 原件名
			partPinStr := n[3]  // 引脚名
			log.Println(partIdxName, "--", partPinStr)
			pinIdx, _ := strconv.Atoi(partPinStr)
			lp.Lc.X, lp.Lc.Y = b.GetIcByIndex(partIdxName).GetPin(pinIdx)
			lastPt = lp.Lc
			cur.Lp = append(cur.Lp, lp)
			desc = desc[len(n[0]):]
		}
		return cur
	}
	return nil
}

// ToSvg 生成svg
func (nb *NetBlock) ToSvg(canvas *svg.SVG, w io.Writer) {
	log.Println("netblock")
	for _, v := range nb.Lp {
		log.Printf("%#v\n", v)
	}

	if len(nb.Lp) < 2 {
		return
	}
	last := nb.Lp[0].Lc
	for i := 1; i < len(nb.Lp); i++ {
		cur := nb.Lp[i].Lc
		canvas.Line(last.X*div, last.Y*div, cur.X*div, cur.Y*div, "stroke:#737375;")
		last = cur
	}
}

// GetIdxName 获取芯片名称
func (nb *NetBlock) GetIdxName() string {
	return ""
}

// GetPin 获取引脚位置
func (nb *NetBlock) GetPin(i int) (x int, y int) {
	return 0, 0
}
