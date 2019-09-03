package mdex

import (
	"log"
	"strings"

	"github.com/antboard/eeblog/mdex/wave"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type waveParser struct {
}

var defaultWaveParser = &waveParser{}

// NewWaveParser 创建一个波形解析器
func NewWaveParser() parser.BlockParser {
	return defaultWaveParser
}

// Open 确认处理波形图
func (w *waveParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	log.Println("w1:", string(line))
	// pos := pc.BlockOffset()
	if !strings.HasPrefix(string(line), "~$") {
		return nil, parser.NoChildren
	}

	node := wave.NewWaveBlock()
	node.InitByLine(string(line))
	return node, parser.NoChildren
}

func (w *waveParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	log.Println("w2:", string(line))
	if util.IsBlank(line) {
		return parser.Continue | parser.NoChildren
	}
	// 如果 是结束符就返回close
	// if line[0] == '$' {
	if strings.HasPrefix(string(line), "~$") {
		log.Println("wend")
		reader.Advance(2)
		return parser.Close | parser.NoChildren
	}
	cur, ok := node.(*wave.WaveBlock)
	if !ok {
		log.Println("w no node")
		log.Printf("w %#v\n", node)
		return parser.Close
	}

	// log.Printf("%#v", string(line))
	cur.AddLine(string(line))
	// reader.Advance(len(line))
	return parser.Continue | parser.NoChildren
}

func (w *waveParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// log.Println("close")
}

func (w *waveParser) CanInterruptParagraph() bool {
	log.Println("w interrupt")
	return false
}

func (w *waveParser) CanAcceptIndentedLine() bool {
	log.Println("w indented")
	return false
}

// WaveHTMLRenderer 渲染器
type WaveHTMLRenderer struct {
	html.Config
}

// NewWaveHTMLRenderer Wave渲染器
func NewWaveHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &WaveHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs 注册渲染函数
func (s *WaveHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(wave.KindWaveBlock, s.renderWaveBlock)
}

func (s *WaveHTMLRenderer) renderWaveBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*wave.WaveBlock)
	n.ToSvg(w)
	// b, _ := json.Marshal(n)
	// w.WriteString(string(b))
	return ast.WalkContinue, nil
}

type waveExt struct {
}

// WaveExt 扩展绑定
var WaveExt = &waveExt{}

func (e *waveExt) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewWaveParser(), 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewWaveHTMLRenderer(), 500),
	))

}
