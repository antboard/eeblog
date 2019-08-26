package wave

import (
	"io"

	svg "github.com/ajstarks/svgo"
)

// WaveParser 波形解析原型
type WaveParser interface {
	CanParse(str string) bool
	ParseLine(b *WaveBlock, str string) SvgBlock
}

// SvgBlock 用来生成svg图
type SvgBlock interface {
	ToSvg(canvas *svg.SVG, w io.Writer)
}
