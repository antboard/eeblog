package mdex

import (
	"log"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
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
	// 判断$前缀标记
	pos := pc.BlockOffset()
	if line[pos] != '$' {
		return nil, parser.NoChildren
	}

	node := last.NewSchBlock()
	remain := node.AddLine(string(line))
	// log.Printf("%#v", node)
	reader.Advance(segment.Len() - remain)
	return node, parser.NoChildren
}

func (b *schParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	log.Println("2")
	line, segment := reader.PeekLine()
	if util.IsBlank(line) {
		return parser.Continue | parser.NoChildren
	}
	// 如果 是结束符就返回close
	if line[0] == '$' {
		return parser.Close
	}
	cur, ok := node.(*last.SchBlock)
	if !ok {
		log.Println("no node")
		log.Printf("%#v\n", node)
		return parser.Close
	}

	// log.Printf("%#v", string(line))
	cur.AddLine(string(line))
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

// SchHTMLRenderer 渲染器
type SchHTMLRenderer struct {
	html.Config
}

// NewSchHTMLRenderer sch渲染器
func NewSchHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &SchHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs 注册渲染函数
func (s *SchHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(last.KindSchBlock, s.renderSchBlock)
}

func (s *SchHTMLRenderer) renderSchBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*last.SchBlock)
	n.ToSvg(w)
	// b, _ := json.Marshal(n)
	// w.WriteString(string(b))
	return ast.WalkContinue, nil
}

type schExt struct {
}

// SchExt 扩展绑定
var SchExt = &schExt{}

func (e *schExt) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewSchParser(), 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewSchHTMLRenderer(), 500),
	))

}
