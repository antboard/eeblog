package wave

import (
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

func init() {
	waveParsers = append(waveParsers, new(Square))

}

// Square 方波描述
// SQ(x,y,c)-010101010101
// 0: 低电平
// 1: 高电平
// 2,4,6,8: 2x,4x,6x,8x
type Square struct {
	X, Y, C   int
	Cfg, Name string
}

// CreateSquare ...
func CreateSquare() *Square {
	return &Square{}
}

// CanParse 类型检查
func (s *Square) CanParse(desc string) bool {
	// SQ 方波
	nx := regexp.MustCompile(`^[\s]*SQ`)
	n := nx.FindStringSubmatch(desc)
	if len(n) > 0 {
		log.Println("parse Square wave success ", desc)
		return true
	}
	return false
}

// ParseLine 解析行
func (s *Square) ParseLine(b *WaveBlock, desc string) SvgBlock {
	// SQ 方波
	nx := regexp.MustCompile(`^[\s]*SQ([a-zA-z0-9]+)\(([0-9]+),([0-9]+),([0-9]+)\)-`)
	n := nx.FindStringSubmatch(desc)
	if len(n) > 0 {
		// 方波正确就可以切到前缀
		prelen := len(n[0])
		desc = desc[prelen:]
		cur := CreateSquare()
		cur.Name = n[1]
		cur.X, _ = strconv.Atoi(n[2])
		cur.Y, _ = strconv.Atoi(n[3])
		cur.C, _ = strconv.Atoi(n[4])
		cur.Cfg = strings.TrimSuffix(desc, "\n")
		log.Printf("%#v\n", cur)
		return cur
	}
	return nil
}

// ToSvg ToSvg
func (s *Square) ToSvg(canvas *svg.SVG, w io.Writer) {
	startX := s.X * div
	Y := 0
	hi := div * 3
	lastplus := false
	canvas.Text(s.X*div, s.Y*div+24, s.Name, `opacity=".3" font-size=24px`)
	for idx := 0; idx < len(s.Cfg); idx++ {
		c := s.Cfg[idx] - '0'
		plus := false
		if (c % 2) == 1 {
			if lastplus == false {
				plus = true
			}
			lastplus = true
			Y = s.Y * div
		} else {
			if lastplus == true {
				plus = true
			}
			lastplus = false
			Y = s.Y*div + hi
		}
		c = c / 2
		endX := startX + int(c+1)*div
		canvas.Line(startX, Y, endX, Y, "stroke:#737375;")
		if plus {
			canvas.Line(startX, s.Y*div, startX, s.Y*div+hi, "stroke:#737375;")
		}
		startX = endX
	}
}
