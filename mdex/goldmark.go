package mdex

import (
	"github.com/yuin/goldmark"
)

// MD 初始化全局渲染器
var MD = goldmark.New(
	goldmark.WithExtensions(SchExt),
	// goldmark.WithParserOptions(
	// 	parser.WithAutoHeadingID(),
	// ),
	// goldmark.WithRendererOptions(
	// 	html.WithHardWraps(),
	// 	html.WithXHTML(),
	// ),
)

func init() {
	// log.SetFlags(log.Llongfile)
	// src := `$(30,20)
	// U10-P8-NSTC12[1:VCC,8:GND](1,2,3)
	// U11-P4-NEEPROM[1:VCC,4:GND](10,12,3)
	// $`
	// var buf bytes.Buffer
	// if err := MD.Convert([]byte(src), &buf); err != nil {
	// 	panic(err)
	// }
	// log.Println(buf.String())
}
