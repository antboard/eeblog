package ast

import (
	"regexp"
	"strings"
	"testing"
)

func TestU(t *testing.T) {
	desc := `U10-P8-NSTC12[1:VCC,8:GND](1,2)`
	// re := regexp.MustCompile(`U([0-9]+)-P([0-9]+)-N([A-Za-z0-9]+)-\[(([0-9]+):[A-Za-z0-9]+,)*\]`)
	// 拆出芯片编号
	ux := regexp.MustCompile(`U([0-9]+)-`)
	u := ux.FindStringSubmatch(desc)
	if len(u) > 1 {
		t.Error("u:", u[1])
		desc = desc[len(u[0]):]
		// t.Error(desc)
	}
	//拆出引脚数量
	px := regexp.MustCompile(`P([0-9]+)-`)
	p := px.FindStringSubmatch(desc)
	if len(p) > 1 {
		t.Error("p", p[1])
		desc = desc[len(p[0]):]
	}
	// 拆出芯片
	nx := regexp.MustCompile(`N([A-Za-z0-9]+)`)
	n := nx.FindStringSubmatch(desc)
	if len(n) > 1 {
		t.Error("n:", n[1])
		desc = desc[len(n[0]):]
	}
	// 如果有[]则解析引脚命名
	nstart := strings.Index(desc, "[")
	if nstart >= 0 {
		nend := strings.Index(desc, "]")
		pinstr := desc[nstart+1 : nend]
		desc = desc[nend+1:]
		pins := strings.Split(pinstr, ",")
		for _, v := range pins {
			apin := strings.Split(v, ":")
			t.Error(apin)
		}
	}
	// 拆出位置信息
	lc := regexp.MustCompile(`\(([0-9]+),([0-9]+)\)`)
	lsl := lc.FindStringSubmatch(desc)
	if len(lsl) >= 3 {
		x := lsl[1]
		y := lsl[2]
		t.Error(x, y)
		desc = desc[len(lsl[0]):]
	}

	t.Error(desc)
}
