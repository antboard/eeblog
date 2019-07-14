package ast

import (
	"io"

	svg "github.com/ajstarks/svgo"
)

// Lineparser 行解析器
// 按行分类, 进行解析保存,最后输出svg
type Lineparser interface {
	CanParse(str string) bool
	ParseLine(b *SchBlock, str string) SvgBlock
}

// SvgBlock 用来生成svg图
type SvgBlock interface {
	ToSvg(canvas *svg.SVG, w io.Writer)
	GetIdxName() string        // 获得块的名字
	GetPin(idx int) (int, int) // 获得引脚的位置
}
