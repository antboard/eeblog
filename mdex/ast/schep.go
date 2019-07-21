package ast

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

// 放置自定义器件
// EpU12-Ref(Ei35)-(50, 20)

func init() {
	schParsers = append(schParsers, new(EpBlock))
}

// EpBlock 放置自定义块
type EpBlock struct {
	Index string // 这个索引可以加自定义字符大写一个
	Ref   string // 引用
	Ud    *UserDefine
	X     int
	Y     int
}

// CanParse 类型检查
func (epb *EpBlock) CanParse(desc string) bool {
	epx := regexp.MustCompile(`^[\s]*Ep([A-Z0-9]+)-`)
	ep := epx.FindStringSubmatch(desc)
	if len(ep) > 1 {
		log.Println("parse Ep success ", desc)
		return true
	}
	return false
}

// ParseLine 解析块
func (epb *EpBlock) ParseLine(b *SchBlock, desc string) SvgBlock {
	epx := regexp.MustCompile(`^[\s]*Ep([A-Z0-9]+)-`)
	ep := epx.FindStringSubmatch(desc)
	if len(ep) > 1 {
		cur := new(EpBlock)
		cur.Index = ep[1]
		desc = desc[len(ep[0]):]
		// Ref(Ei35)-
		refx := regexp.MustCompile(`Ref\(([0-9a-zA-Z]+)\)-`)
		refsstr := refx.FindStringSubmatch(desc)
		if len(refsstr) > 1 {
			cur.Ref = refsstr[1]
			desc = desc[len(refsstr[0]):]
		}
		// (50, 20)
		fmt.Println(desc)
		lcx := regexp.MustCompile(`\(([0-9]+),([0-9]+)\)`)
		lcsstr := lcx.FindStringSubmatch(desc)
		if len(lcsstr) > 1 {
			cur.X, _ = strconv.Atoi(lcsstr[1])
			cur.Y, _ = strconv.Atoi(lcsstr[2])
		}
		fmt.Println(lcsstr)
		// 找到引用的芯片直接引用
		cur.Ud = b.GetIcByIndex(cur.Ref).(*UserDefine)
		if cur.Ud == nil {
			fmt.Println("no ref", cur.Ref)
		}
		fmt.Printf("%#v\n", cur.Ud)
		fmt.Printf("%#v\n", cur)
		return cur
	}
	return nil
}

// ToSvg 生成svg
func (epb *EpBlock) ToSvg(canvas *svg.SVG, w io.Writer) {
	// 基于当前位置,进行相关画图
	epb.Ud.LayoutSvg(canvas, w, epb.X, epb.Y, epb.Index, "")
}

// GetIdxName 获取元件名
func (epb *EpBlock) GetIdxName() string {
	return epb.Index
}

// GetPin 获取引脚位置
func (epb *EpBlock) GetPin(i int) (x int, y int) {
	udx, udy := epb.Ud.GetPin(i)
	return epb.X + udx, epb.Y + udy
}
