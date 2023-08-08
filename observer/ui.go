package main

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"
	"regexp"

	"github.com/Jille/convreq"
	"github.com/Jille/convreq/respond"
	pb "github.com/SnoozeThis-org/logwait/proto"
	"golang.org/x/exp/slices"
)

//go:embed template.html
var rawTemplate string

var tpl = template.Must(template.New("").Parse(rawTemplate))

func (s *service) registerUI() {
	http.Handle("/", convreq.Wrap(s.handleHttp))
}

type templateVars struct {
	Warnings []string
	Fields   []string
	Values   map[string]string
	Errors   map[string]string
	URL      string
}

func (s *service) handleHttp(ctx context.Context, req *http.Request) convreq.HttpResponse {
	if err := req.ParseForm(); err != nil {
		return respond.BadRequest(err.Error())
	}
	v := templateVars{}
	v.Fields = []string{"message"}
	v.Values = map[string]string{}
	v.Errors = map[string]string{}
	for f, vs := range req.Form {
		if len(vs) == 0 {
			continue
		}
		v.Values[f] = vs[0]
		if !slices.Contains(v.Fields, f) {
			v.Fields = append(v.Fields, f)
		}
		if _, err := regexp.Compile(vs[0]); err != nil {
			v.Errors[f] = err.Error()
		}
	}
	if len(v.Errors) == 0 && len(v.Values) > 0 {
		o := &pb.Observable{}
		for f, r := range v.Values {
			o.Filters = append(o.Filters, &pb.Filter{
				Field:  f,
				Regexp: r,
			})
		}
		o.Signature = calculateSignature(o)
		s.mtx.Lock()
		base := s.registeredUrl
		s.mtx.Unlock()
		if base == "" {
			v.Warnings = append(v.Warnings, "We haven't been connected to SnoozeThis since startup yet and don't know the registered domain name yet. Making a guess that it is "+req.Host)
			base = "https://" + req.Host + "/"
		}
		req.Form.Set("snooze_signature", o.Signature)
		v.URL = base + "?" + req.Form.Encode()
	}
	return respond.RenderTemplate(tpl, v)
}
