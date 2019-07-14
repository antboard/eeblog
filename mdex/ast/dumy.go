package ast

import (
	"io"

	svg "github.com/ajstarks/svgo"
)

// DumySvgBlock 获取ic位置如果获取不到, 传输一个假的默认0,0
type DumySvgBlock struct{}

// ToSvg ...
func (d *DumySvgBlock) ToSvg(canvas *svg.SVG, w io.Writer) {

}

// GetIdxName 获得块的名字
func (d *DumySvgBlock) GetIdxName() string {
	return ""
}

// GetPin 获得引脚的位置
func (d *DumySvgBlock) GetPin(idx int) (int, int) {
	return 0, 0
}
