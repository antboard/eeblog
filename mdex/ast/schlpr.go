package ast

import (
	"io"

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
	Index  int
	Value  string
	Layout string
	X      int
	Y      int
}

// CanParse 类型检查
func (r *RBlock) CanParse(desc string) bool {
	return false
}

// ParseLine 解析块
func (r *RBlock) ParseLine(str string) SvgBlock {
	return nil
}

// ToSvg 生成svg
func (r *RBlock) ToSvg(canvas *svg.SVG, w io.Writer) {

}
