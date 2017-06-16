package route

import (
	"reflect"
	"regexp"
	"testing"
	"github.com/millisecond/olb/model"
)

func TestParse(t *testing.T) {
	tests := []struct {
		desc string
		in   string
		out  []*model.RouteDef
		fail bool
	}{
		// error flows
		{"FailEmpty", ``, nil, false},
		{"FailNoRoute", `bang`, nil, true},
		{"FailRouteNoCmd", `route x`, nil, true},
		{"FailRouteAddNoService", `route add`, nil, true},
		{"FailRouteAddNoSrc", `route add svc`, nil, true},

		// happy flows
		{
			desc: "RouteAddService",
			in:   `route add svc /prefix http://1.2.3.4/`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/"}},
		},
		{
			desc: "RouteAddTCPService",
			in:   `route add svc :1234 tcp://1.2.3.4:5678`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: ":1234", Dst: "tcp://1.2.3.4:5678"}},
		},
		{
			desc: "RouteAddServiceWeight",
			in:   `route add svc /prefix http://1.2.3.4/ weight 1.2`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/", Weight: 1.2}},
		},
		{
			desc: "RouteAddServiceWeightTags",
			in:   `route add svc /prefix http://1.2.3.4/ weight 1.2 tags "a,b"`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/", Weight: 1.2, Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteAddServiceWeightOpts",
			in:   `route add svc /prefix http://1.2.3.4/ weight 1.2 opts "foo=bar baz=bang blimp"`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/", Weight: 1.2, Opts: map[string]string{"foo": "bar", "baz": "bang", "blimp": ""}}},
		},
		{
			desc: "RouteAddServiceWeightTagsOpts",
			in:   `route add svc /prefix http://1.2.3.4/ weight 1.2 tags "a,b" opts "foo=bar baz=bang blimp"`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/", Weight: 1.2, Tags: []string{"a", "b"}, Opts: map[string]string{"foo": "bar", "baz": "bang", "blimp": ""}}},
		},
		{
			desc: "RouteAddServiceWeightTagsOptsMoreSpaces",
			in:   ` route  add  svc  /prefix  http://1.2.3.4/  weight  1.2  tags  " a , b "  opts  "foo=bar  baz=bang  blimp" `,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/", Weight: 1.2, Tags: []string{"a", "b"}, Opts: map[string]string{"foo": "bar", "baz": "bang", "blimp": ""}}},
		},
		{
			desc: "RouteAddTags",
			in:   `route add svc /prefix http://1.2.3.4/ tags "a,b"`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/", Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteAddTagsOpts",
			in:   `route add svc /prefix http://1.2.3.4/ tags "a,b" opts "foo=bar baz=bang blimp"`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/", Tags: []string{"a", "b"}, Opts: map[string]string{"foo": "bar", "baz": "bang", "blimp": ""}}},
		},
		{
			desc: "RouteAddOpts",
			in:   `route add svc /prefix http://1.2.3.4/ opts "foo=bar baz=bang blimp"`,
			out:  []*model.RouteDef{{Cmd: model.RouteAddCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/", Opts: map[string]string{"foo": "bar", "baz": "bang", "blimp": ""}}},
		},
		{
			desc: "RouteDelTags",
			in:   `route del tags "a,b"`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteDelTagsMoreSpaces",
			in:   `route  del  tags  " a , b "`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteDelService",
			in:   `route del svc`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Service: "svc"}},
		},
		{
			desc: "RouteDelServiceTags",
			in:   `route del svc tags "a,b"`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Service: "svc", Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteDelServiceTagsMoreSpaces",
			in:   `route  del  svc  tags  " a , b "`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Service: "svc", Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteDelServiceSrc",
			in:   `route del svc /prefix`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Service: "svc", Src: "/prefix"}},
		},
		{
			desc: "RouteDelTCPServiceSrc",
			in:   `route del svc :1234`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Service: "svc", Src: ":1234"}},
		},
		{
			desc: "RouteDelServiceSrcDst",
			in:   `route del svc /prefix http://1.2.3.4/`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/"}},
		},
		{
			desc: "RouteDelTCPServiceSrcDst",
			in:   `route del svc :1234 tcp://1.2.3.4:5678`,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Service: "svc", Src: ":1234", Dst: "tcp://1.2.3.4:5678"}},
		},
		{
			desc: "RouteDelServiceSrcDstMoreSpaces",
			in:   ` route  del  svc  /prefix  http://1.2.3.4/ `,
			out:  []*model.RouteDef{{Cmd: model.RouteDelCmd, Service: "svc", Src: "/prefix", Dst: "http://1.2.3.4/"}},
		},
		{
			desc: "RouteWeightServiceSrc",
			in:   `route weight svc /prefix weight 1.2`,
			out:  []*model.RouteDef{{Cmd: model.RouteWeightCmd, Service: "svc", Src: "/prefix", Weight: 1.2}},
		},
		{
			desc: "RouteWeightServiceSrcTags",
			in:   `route weight svc /prefix weight 1.2 tags "a,b"`,
			out:  []*model.RouteDef{{Cmd: model.RouteWeightCmd, Service: "svc", Src: "/prefix", Weight: 1.2, Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteWeightServiceSrcTagsMoreSpaces",
			in:   ` route  weight  svc  /prefix  weight  1.2  tags  " a , b " `,
			out:  []*model.RouteDef{{Cmd: model.RouteWeightCmd, Service: "svc", Src: "/prefix", Weight: 1.2, Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteWeightSrcTags",
			in:   `route weight /prefix weight 1.2 tags "a,b"`,
			out:  []*model.RouteDef{{Cmd: model.RouteWeightCmd, Src: "/prefix", Weight: 1.2, Tags: []string{"a", "b"}}},
		},
		{
			desc: "RouteWeightSrcTagsMoreSpaces",
			in:   ` route  weight  /prefix  weight  1.2  tags  " a , b " `,
			out:  []*model.RouteDef{{Cmd: model.RouteWeightCmd, Src: "/prefix", Weight: 1.2, Tags: []string{"a", "b"}}},
		},
	}

	reSyntaxError := regexp.MustCompile(`syntax error`)

	deref := func(def []*model.RouteDef) (defs []model.RouteDef) {
		for _, d := range def {
			defs = append(defs, *d)
		}
		return
	}

	run := func(in string, def []*model.RouteDef, fail bool, parseFn func(string) ([]*model.RouteDef, error)) {
		out, err := parseFn(in)
		switch {
		case err == nil && fail:
			t.Errorf("got error nil want fail")
			return
		case err != nil && !fail:
			t.Errorf("got error %v want nil", err)
			return
		case err != nil:
			if !reSyntaxError.MatchString(err.Error()) {
				t.Errorf("got error %q want 'syntax error.*'", err)
			}
			return
		}
		if got, want := out, def; !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot  %#v\nwant %#v", deref(got), deref(want))
		}
	}

	for _, tt := range tests {
		t.Run("Parse-"+tt.desc, func(t *testing.T) { run(tt.in, tt.out, tt.fail, Parse) })
	}
}
