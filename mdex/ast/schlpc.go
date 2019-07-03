package ast

import (
	"io"

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
	Index    int
	Value    string
	Layout   string
	Polarity bool
	X        int
	Y        int
}

// CanParse 类型检查
func (r *CBlock) CanParse(desc string) bool {
	return false
}

// ParseLine 解析块
func (r *CBlock) ParseLine(str string) SvgBlock {
	return nil
}

// ToSvg 生成svg
func (r *CBlock) ToSvg(canvas *svg.SVG, w io.Writer) {

}
