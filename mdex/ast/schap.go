package ast

import (
	"io"
	"log"
	"regexp"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

// 锚点
// 在此是提供视觉效果悬浮的点.

func init() {
	schParsers = append(schParsers, new(AcPoint))
}

// AcPoint 锚点
type AcPoint struct {
	X, Y int
}

// CreatAcPoint ...
func CreatAcPoint() *AcPoint {
	return &AcPoint{}
}

// CanParse 类型检查
func (ac *AcPoint) CanParse(desc string) bool {
	// 如果有A出现就是锚点
	nx := regexp.MustCompile(`^[\s]*A`)
	n := nx.FindStringSubmatch(desc)
	if len(n) > 0 {
		log.Println("parse Anchor Point success ", desc)
		return true
	}

	return false
}

// ParseLine 解析行定义
func (ac *AcPoint) ParseLine(b *SchBlock, desc string) SvgBlock {
	nx := regexp.MustCompile(`^[\s]*A\(([0-9]+),([0-9]+)\)`)
	n := nx.FindStringSubmatch(desc)
	if len(n) > 1 {
		cur := CreatAcPoint()
		cur.X, _ = strconv.Atoi(n[1])
		cur.Y, _ = strconv.Atoi(n[2])
		return cur
	}
	return nil
}

// ToSvg ToSvg
func (ac *AcPoint) ToSvg(canvas *svg.SVG, w io.Writer) {
	log.Println(ac.X, ac.Y)
	canvas.Circle(ac.X*div, ac.Y*div, 2)
}

// GetIdxName 获取芯片名称
func (ac *AcPoint) GetIdxName() string {
	return ""
}

// GetPin 获取引脚位置
func (ac *AcPoint) GetPin(i int) (x int, y int) {
	return 0, 0
}
