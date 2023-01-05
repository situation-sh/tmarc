// Code generated by xgen. DO NOT EDIT.

package main

import (
	"encoding/xml"
)

// DateRangeType ...
type DateRangeType struct {
	Begin int `xml:"begin"`
	End   int `xml:"end"`
}

// ReportMetadataType ...
type ReportMetadataType struct {
	Orgname          string         `xml:"org_name"`
	Email            string         `xml:"email"`
	Extracontactinfo string         `xml:"extra_contact_info"`
	Reportid         string         `xml:"report_id"`
	Daterange        *DateRangeType `xml:"date_range"`
	Error            []string       `xml:"error"`
}

// AlignmentType ...
type AlignmentType string

// DispositionType ...
type DispositionType string

// PolicyPublishedType ...
type PolicyPublishedType struct {
	Domain string `xml:"domain"`
	Adkim  string `xml:"adkim"`
	Aspf   string `xml:"aspf"`
	P      string `xml:"p"`
	Sp     string `xml:"sp"`
	Pct    int    `xml:"pct"`
}

// DMARCResultType ...
type DMARCResultType string

// PolicyOverrideType ...
type PolicyOverrideType string

// PolicyOverrideReason ...
type PolicyOverrideReason struct {
	Type    string `xml:"type"`
	Comment string `xml:"comment"`
}

// PolicyEvaluatedType ...
type PolicyEvaluatedType struct {
	Disposition string                  `xml:"disposition"`
	Dkim        string                  `xml:"dkim"`
	Spf         string                  `xml:"spf"`
	Reason      []*PolicyOverrideReason `xml:"reason"`
}

// IPAddress ...
type IPAddress string

// RowType ...
type RowType struct {
	Sourceip        string               `xml:"source_ip"`
	Count           int                  `xml:"count"`
	Policyevaluated *PolicyEvaluatedType `xml:"policy_evaluated"`
}

// IdentifierType ...
type IdentifierType struct {
	Envelopeto string `xml:"envelope_to"`
	Headerfrom string `xml:"header_from"`
}

// DKIMResultType ...
type DKIMResultType string

// DKIMAuthResultType ...
type DKIMAuthResultType struct {
	Domain      string `xml:"domain"`
	Selector    string `xml:"selector"`
	Result      string `xml:"result"`
	Humanresult string `xml:"human_result"`
}

// SPFResultType ...
type SPFResultType string

// SPFAuthResultType ...
type SPFAuthResultType struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
}

// AuthResultType ...
type AuthResultType struct {
	Dkim []*DKIMAuthResultType `xml:"dkim"`
	Spf  []*SPFAuthResultType  `xml:"spf"`
}

// RecordType ...
type RecordType struct {
	Row         *RowType        `xml:"row"`
	Identifiers *IdentifierType `xml:"identifiers"`
	Authresults *AuthResultType `xml:"auth_results"`
}

// Feedback ...
type Feedback struct {
	XMLName         xml.Name             `xml:"feedback"`
	Reportmetadata  *ReportMetadataType  `xml:"report_metadata"`
	Policypublished *PolicyPublishedType `xml:"policy_published"`
	Record          []*RecordType        `xml:"record"`
}
