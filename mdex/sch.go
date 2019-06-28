package mdex

import (
	"log"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"

	last "github.com/antboard/eeblog/mdex/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

//* $$$ U10-P8-NSTC12[1:VCC,8:GND] $$$

type schParser struct {
}

var defaultSchParser = &schParser{}

// NewSchParser 创建一个解析器
func NewSchParser() parser.BlockParser {
	return defaultSchParser
}

func (b *schParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	log.Println("1")
	line, segment := reader.PeekLine()
	pos, padding := util.IndentPosition(line, reader.LineOffset(), 4)
	if pos < 0 {
		return nil, parser.NoChildren
	}
	node := last.NewSchBlock()
	reader.AdvanceAndSetPadding(pos, padding)
	_, segment = reader.PeekLine()
	// node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return node, parser.NoChildren
}

func (b *schParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	log.Println("2")
	line, segment := reader.PeekLine()
	if util.IsBlank(line) {
		return parser.Continue | parser.NoChildren
	}
	pos, padding := util.IndentPosition(line, reader.LineOffset(), 4)
	if pos < 0 {
		return parser.Close
	}
	reader.AdvanceAndSetPadding(pos, padding)
	_, segment = reader.PeekLine()
	// node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return parser.Continue | parser.NoChildren
}

func (b *schParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	log.Println("close")
}

func (b *schParser) CanInterruptParagraph() bool {
	log.Println("interrupt")
	return false
}

func (b *schParser) CanAcceptIndentedLine() bool {
	log.Println("indented")
	return false
}

type schExt struct {
}

// SchExt 扩展绑定
var SchExt = &schExt{}

func (e *schExt) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewSchParser(), 0),
	))
	// m.Renderer().AddOptions(renderer.WithNodeRenderers(
	// 	util.Prioritized(NewTaskCheckBoxHTMLRenderer(), 500),
	// ))

}
