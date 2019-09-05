package wave

import (
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strconv"

	svg "github.com/ajstarks/svgo"
	gast "github.com/yuin/goldmark/ast"
)

const (
	div = 12 // 1div = 12 pt
)

// 保存波形解析器
var waveParsers = make([]WaveParser, 0, 10)

// WaveBlock 波形块
type WaveBlock struct {
	gast.BaseBlock
	Waves []SvgBlock
	PageW int
	PageH int
}

// Dump 继承
func (w *WaveBlock) Dump(source []byte, level int) {
	m := make(map[string]string)
	bic, _ := json.Marshal(w.Waves)
	m["ics"] = string(bic)
	gast.DumpHelper(w, source, level, m, nil)
}

// KindWaveBlock 原理图描述类
var KindWaveBlock = gast.NewNodeKind("WaveBlock")

// Kind implements Node.Kind.
func (w *WaveBlock) Kind() gast.NodeKind {
	return KindWaveBlock
}

// InitByLine 初始化画布
func (w *WaveBlock) InitByLine(desc string) {
	// 读取画布大小
	pageszre := regexp.MustCompile(`~\$\(([0-9]+),([0-9]+)\)`)
	pagesz := pageszre.FindStringSubmatch(desc)
	log.Println(desc, pagesz)
	if len(pagesz) == 3 {
		w.PageW, _ = strconv.Atoi(pagesz[1])
		w.PageH, _ = strconv.Atoi(pagesz[2])
		return
	}
}

// AddLine 添加一个行描述符
func (w *WaveBlock) AddLine(desc string) int {
	log.Println(desc)
	for _, v := range waveParsers {
		if v.CanParse(desc) {
			lp := v.ParseLine(w, desc)
			if lp != nil {
				w.Waves = append(w.Waves, lp)
			}
			break
		}
	}
	return len(desc)
}

// ToSvg 输出svg
func (w *WaveBlock) ToSvg(wr io.Writer) {
	width := w.PageW * div
	height := w.PageH * div
	canvas := svg.New(wr)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:#e0e0e2;")
	for _, v := range w.Waves {
		v.ToSvg(canvas, wr)
	}

	canvas.End()
}

// NewWaveBlock 解析出一个新波形
func NewWaveBlock() *WaveBlock {
	return &WaveBlock{PageW: 50, PageH: 50}
}
