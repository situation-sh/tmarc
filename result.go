package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var columns = []string{
	"end", "source", "domain", "count", "dkim", "spf",
}

const dateFormat = "Mon, 02 Jan 2006"

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	f := time.Time(d).Format(dateFormat)
	return []byte(fmt.Sprintf(`"%s"`, f)), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(dateFormat, string(b))
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

type FeedbackResult struct {
	SourceFile string `json:"source_file"`
	OrgName    string `json:"org_name"`
	ReportID   string `json:"report_id"`
	Begin      Date   `json:"begin"`
	End        Date   `json:"end"`
	SourceIP   net.IP `json:"source_ip"`
	Source     string `json:"source"`
	Count      int    `json:"count"`
	EnvelopeTo string `json:"envelope_to"`
	HeaderFrom string `json:"header_from"`
	Domain     string `json:"domain"`
	DKIMResult string `json:"dkim"`
	SPFResult  string `json:"spf"`
}

func (r *FeedbackResult) Columns(nbcols int) []string {
	return columns
}

func (r *FeedbackResult) ToRow(nbcols int) []string {
	b, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	cols := r.Columns(nbcols)
	out := make([]string, len(cols))
	for i, c := range columns {
		// if c == "dkim" {
		// 	switch x := m[c].(type) {
		// 	case string:
		// 		if len(x) == 0 {
		// 			out[i] = ""
		// 			continue
		// 		}
		// 		s := lipgloss.NewStyle().
		// 			Foreground(lipgloss.Color("#FF4500")).
		// 			Padding(0).
		// 			Margin(0).
		// 			BorderStyle(lipgloss.RoundedBorder()).
		// 			BorderBackground(lipgloss.Color("63")).
		// 			Width(len(x))
		// 		out[i] = s.Inline(true).Render(x)
		// 		// fmt.Println(len(out[i]))
		// 	}

		// } else {
		// 	// s := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500"))
		// 	// out[i] = s.Render(fmt.Sprintf("%v", m[c]))
		// 	out[i] = fmt.Sprintf("%v", m[c])
		// }
		out[i] = fmt.Sprintf("%v", m[c])
	}
	return out
}

type FeedbackResults []*FeedbackResult

func (r FeedbackResults) Len() int {
	return len(r)
}

func (r FeedbackResults) Files() int {
	m := make(map[string]bool)
	for _, x := range r {
		m[x.SourceFile] = true
	}
	return len(m)
}

func (r FeedbackResults) Less(i, j int) bool {
	ei := time.Time(r[i].End)
	ej := time.Time(r[j].End)
	if ei.Equal(ej) {
		return r[i].Source <= r[j].Source
	} else {
		return ei.Before(ej)
	}

}

func (r FeedbackResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type domainResult struct {
	dkim string
	spf  string
}

func parseFeedback(feedback *Feedback, sourceFile string) FeedbackResults {
	results := make([]*FeedbackResult, 0)
	var source string
	reportID := feedback.Reportmetadata.Reportid
	orgName := feedback.Reportmetadata.Orgname
	begin := Date(time.Unix(int64(feedback.Reportmetadata.Daterange.Begin), 0))

	end := Date(time.Unix(int64(feedback.Reportmetadata.Daterange.End), 0))

	for _, row := range feedback.Record {
		m := make(map[string]*domainResult)
		sourceIP := net.ParseIP(row.Row.Sourceip)
		names, err := net.LookupAddr(row.Row.Sourceip)
		if err == nil {
			source = names[0]
		} else {
			source = ""
		}
		// spf

		for _, spf := range row.Authresults.Spf {
			r, exists := m[spf.Domain]
			if exists {
				r.spf = spf.Result
			} else {
				m[spf.Domain] = &domainResult{spf: spf.Result}
			}
		}
		// dkim
		for _, dkim := range row.Authresults.Dkim {
			r, exists := m[dkim.Domain]
			if exists {
				r.dkim = dkim.Result
			} else {
				m[dkim.Domain] = &domainResult{dkim: dkim.Result}
			}
		}

		for d, r := range m {
			results = append(
				results,
				&FeedbackResult{
					SourceFile: sourceFile,
					OrgName:    orgName,
					ReportID:   reportID,
					Begin:      begin,
					End:        end,
					SourceIP:   sourceIP,
					Source:     source,
					Count:      row.Row.Count,
					EnvelopeTo: row.Identifiers.Envelopeto,
					HeaderFrom: row.Identifiers.Headerfrom,
					Domain:     d,
					SPFResult:  r.spf,
					DKIMResult: r.dkim,
				},
			)
		}

	}

	return results
}
