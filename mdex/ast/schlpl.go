package ast

import (
	"io"

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
	Index  int
	Value  string
	Layout string
	X      int
	Y      int
}

// CanParse 类型检查
func (r *LBlock) CanParse(desc string) bool {
	return false
}

// ParseLine 解析块
func (r *LBlock) ParseLine(str string) SvgBlock {
	return nil
}

// ToSvg 生成svg
func (r *LBlock) ToSvg(canvas *svg.SVG, w io.Writer) {

}
