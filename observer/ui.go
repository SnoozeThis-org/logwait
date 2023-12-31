package main

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/Jille/convreq"
	"github.com/Jille/convreq/respond"
	"github.com/Jille/genericz/mapz"
	pb "github.com/SnoozeThis-org/logwait/proto"
)

//go:embed template.html
var rawTemplate string

var tpl = template.Must(template.New("").Parse(rawTemplate))

var commonTestFilterErrorRe = regexp.MustCompile(`^invalid regexp for field (\S+): (.+)$`)

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
	o := &pb.Observable{}
	v := templateVars{
		Values: map[string]string{},
		Errors: map[string]string{},
	}
	fields := map[string]struct{}{}
	for f, vs := range req.Form {
		if len(vs) == 0 || len(vs[0]) == 0 {
			continue
		}
		v.Values[f] = vs[0]
		fields[f] = struct{}{}
		o.Filters = append(o.Filters, &pb.Filter{
			Field:  f,
			Regexp: vs[0],
		})
	}

	var hasErrors bool
	if len(o.Filters) > 0 {
		if resp, ok := s.askTestFilters(ctx, &pb.TestFiltersRequest{Filters: o.Filters}); !ok {
			v.Warnings = append(v.Warnings, "Failed to communicate with your scanners to verify your filters")
		} else {
			for _, e := range resp.Errors {
				hasErrors = true
				if m := commonTestFilterErrorRe.FindStringSubmatch(e); m != nil {
					v.Errors[m[1]] = m[2]
				} else {
					v.Warnings = append(v.Warnings, e)
				}
			}
		}
	}

	s.mtx.Lock()
	base := s.registeredUrl
	numScanners := len(s.connectedScanners)
	for _, fs := range s.fieldsPerScanner {
		for _, f := range fs {
			fields[f] = struct{}{}
		}
	}
	s.mtx.Unlock()
	if base == "" {
		v.Warnings = append(v.Warnings, "We haven't been connected to SnoozeThis since startup yet and don't know the registered domain name yet. Making a guess that it is "+req.Host)
		base = "https://" + req.Host + "/"
	}
	if numScanners == 0 {
		v.Warnings = append(v.Warnings, "There are currently no scanners connected, so we aren't seeing any log lines at all.")
	}
	if len(fields) > 0 {
		v.Fields = mapz.KeysSorted(fields)
	} else {
		v.Fields = []string{"message"}
	}

	if len(v.Values) > 0 && !hasErrors {
		if *signingKey != "" {
			o.Signature = calculateSignature(o)
			req.Form.Set("snooze_signature", o.Signature)
		}
		v.URL = base + "?" + req.Form.Encode()
	}
	return respond.RenderTemplate(tpl, v)
}

func (s *service) askTestFilters(ctx context.Context, req *pb.TestFiltersRequest) (*pb.TestFiltersResponse, bool) {
	ch := make(chan *pb.TestFiltersResponse, 1)
	s.mtx.Lock()
	s.nextTestId++
	id := s.nextTestId
	req.TestId = id
	s.pendingTestFilters[id] = ch

	var sentRequest bool
	for scanner := range s.connectedScanners {
		if err := scanner.Send(&pb.ObserverToScanner{
			Msg: &pb.ObserverToScanner_TestFilters{
				TestFilters: req,
			},
		}); err == nil { // err == nil
			sentRequest = true
			break
		}
	}
	s.mtx.Unlock()

	defer mapz.DeleteWithLock(&s.mtx, s.pendingTestFilters, id)

	if !sentRequest {
		return nil, false
	}

	select {
	case <-ctx.Done():
		return nil, false
	case <-time.After(time.Second):
		return nil, false
	case resp := <-ch:
		return resp, true
	}
}

func (s *service) handleTestFiltersResponse(resp *pb.TestFiltersResponse) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	ch, ok := s.pendingTestFilters[resp.GetTestId()]
	if !ok {
		return
	}
	select {
	case ch <- resp:
	default:
	}
}
