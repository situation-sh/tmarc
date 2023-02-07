package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"time"
)

var columns = []string{
	"end", "source", "header_from", "count", "dkim", "spf",
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

// <record>
//
//	  <row>
//	    <source_ip>209.90.90.209</source_ip>  	// The source IP is the address the message was sent from.
//	    <count>1</count> 						// The sum of the number of messages in that particular report subset
//	    <policy_evaluated>
//	      <disposition>none</disposition>		// Tells you what happened to the messages: none, quarantine or reject
//	      <dkim>pass</dkim>
//	      <spf>fail</spf> 						// the SPF alignment test, which verifies that both the "From" field in the message header & the RFC 5321 "MAIL FROM" are from the same domain
//	    </policy_evaluated>
//	  </row>
//	  <identifiers>
//	    <header_from>evilcorp.com</header_from> // The "From:" header field indicates who the author of the message is. This is the address that is usually visible in the email client as sender. This information can be spoofed
//	  </identifiers>
//	  <auth_results>
//	    <dkim>
//	      <domain>evilcorp.com</domain>
//	      <result>pass</result>
//	      <selector>gm1</selector>
//	    </dkim>
//	    <spf>
//	      <domain>magic.io</domain> 			// "MAIL FROM" field in the SMTP envelope
//	      <result>pass</result> 				// tests whether or not the sending Mail Transfer Agent (MTA) is an authorized sender for the domain in the domain in the RFC 5321 "MAIL FROM".
//	    </spf>
//	  </auth_results>
//	</record>
//
// When spf -> result pass while policy_evaluated -> spf failed, it is a identifier misalignement.
// This failure can intentional if the sending service wants to handle bounce messages [1] for you, which prevents
// SPF from aligning.
// [1] When a mail server fails to deliver a message, it should send a so-called bounce message to the sender=
// to notify them about the failed delivery. Since bounce messages are automatic responses, they must be
// sent to the MAIL FROM address of the envelope.
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
	DKIMResult string `json:"dkim"`
	SPFResult  string `json:"spf"`
	Reason     string `json:"reason"`
	XML        []byte `json:"xml"`
}

func (r *FeedbackResult) Columns() []string {
	return columns
}

func (r *FeedbackResult) ToRow() []string {
	b, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	cols := r.Columns()
	out := make([]string, len(cols))
	for i, c := range columns {
		// fallback to ip
		if c == "source" && m[c] == "" {
			m[c] = m["source_ip"]
		}
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

// type domainResult struct {
// 	dkim string
// 	spf  string
// }

func parseFeedback(feedback *Feedback, sourceFile string) FeedbackResults {
	results := make([]*FeedbackResult, 0)
	var source string
	reportID := feedback.Reportmetadata.Reportid
	orgName := feedback.Reportmetadata.Orgname
	begin := Date(time.Unix(int64(feedback.Reportmetadata.Daterange.Begin), 0))

	end := Date(time.Unix(int64(feedback.Reportmetadata.Daterange.End), 0))

	for _, record := range feedback.Record {
		raw, err := xml.MarshalIndent(record, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(raw))
		sourceIP := net.ParseIP(record.Row.Sourceip)
		names, err := net.LookupAddr(record.Row.Sourceip)
		if err == nil {
			source = names[0]
		} else {
			source = ""
		}

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
				Count:      record.Row.Count,
				EnvelopeTo: record.Identifiers.Envelopeto,
				HeaderFrom: record.Identifiers.Headerfrom,
				// Domain:     d,
				SPFResult:  record.Row.Policyevaluated.Spf,
				DKIMResult: record.Row.Policyevaluated.Dkim,
				XML:        raw,
			},
		)
		// spf

		// m[row.Identifiers.Headerfrom] = &domainResult{}
		// for _, spf := range row.Authresults.Spf {
		// 	r, exists := m[spf.Domain]
		// 	if exists {
		// 		r.spf = spf.Result
		// 	} else {
		// 		m[spf.Domain] = &domainResult{spf: spf.Result}
		// 	}
		// }
		// // dkim
		// for _, dkim := range row.Authresults.Dkim {
		// 	r, exists := m[dkim.Domain]
		// 	if exists {
		// 		r.dkim = dkim.Result
		// 	} else {
		// 		m[dkim.Domain] = &domainResult{dkim: dkim.Result}
		// 	}
		// }

		// for d, r := range m {
		// 	results = append(
		// 		results,
		// 		&FeedbackResult{
		// 			SourceFile: sourceFile,
		// 			OrgName:    orgName,
		// 			ReportID:   reportID,
		// 			Begin:      begin,
		// 			End:        end,
		// 			SourceIP:   sourceIP,
		// 			Source:     source,
		// 			Count:      row.Row.Count,
		// 			EnvelopeTo: row.Identifiers.Envelopeto,
		// 			HeaderFrom: row.Identifiers.Headerfrom,
		// 			Domain:     d,
		// 			SPFResult:  r.spf,
		// 			DKIMResult: r.dkim,
		// 		},
		// 	)
		// }

	}

	return results
}
