package release

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

type RefType int

const (
	Bugzilla RefType = iota
	Jira
)

const (
	maxRetries = 2

	bugzillaBrowseURL = "https://bugzilla.redhat.com/show_bug.cgi?id=%s"
	bugzillaRestURL   = "https://bugzilla.redhat.com/rest/bug"
	jiraBrowseURL     = "https://issues.redhat.com/browse/%s"
	JiraRestURL       = "https://issues.redhat.com/rest/api/latest/issue/%s"
)

var (
	// find strings like "Bug: 123,456 789: some commit description"
	reBugzillaBug = regexp.MustCompile(`^[B|b]ug(?:s)?\s([\d\s,]+)(?:-|:)\s*`)

	// find strings like "XYZ-123,ABC-456 DEF-789: some commit description"
	reJiraRef = regexp.MustCompile(`^([A-Z]+-\d+[\s,]*)+(?:-|:)\s*`)
)

// RefList stores Bug list
type RefList struct {
	Refs []Ref
}

// Ref stores ref ids and
// source of ref type like Bugzilla or Jira.
type Ref struct {
	ID     string
	Source RefType
}

// RefRemoteList stores RefRemoteInfo list
type RefRemoteList struct {
	Refs []RefRemoteInfo `json:"refs"`
}

// RefRemoteInfo stores the detail information of bug retrieved
// from remote url like bugzilla or jira.
type RefRemoteInfo struct {
	ID       string  `json:"id"`
	Status   string  `json:"status"`
	Priority string  `json:"priority"`
	Summary  string  `json:"summary"`
	Source   RefType `json:"source"`
}

type JiraRemoteRef struct {
	Key    string           `json:"key"`
	Fields JiraRemoteFields `json:"fields"`
}

type JiraRemoteFields struct {
	Summary  string             `json:"summary"`
	Status   JiraRemoteStatus   `json:"status"`
	Priority JiraRemotePriority `json:"priority"`
}

type JiraRemoteStatus struct {
	Name string `json:"name"`
}

type JiraRemotePriority struct {
	Name string `json:"name"`
}

// RetrieveRefs retrieves ref details and fills RefRemoteList
// by sending request to Bugzilla and Jira according to the source type of ref.
func RetrieveRefs(refs []Ref) (*RefRemoteList, error) {
	var brl RefRemoteList
	var bugzilla, jira []Ref
	for _, v := range refs {
		if v.Source == Bugzilla {
			bugzilla = append(bugzilla, v)
		} else {
			jira = append(jira, v)
		}
	}

	var lastErr error
	bz, err := retrieveRefsBugzila(bugzilla)
	if err != nil {
		lastErr = err
	}
	if bz != nil {
		brl.Refs = append(brl.Refs, bz.Refs...)
	}

	jr, err := retrieveRefsJira(jira)
	if err != nil {
		lastErr = err
	}
	if jr != nil {
		brl.Refs = append(brl.Refs, jr.Refs...)
	}

	return &brl, lastErr
}

func retrieveRefsJira(refs []Ref) (*RefRemoteList, error) {
	if len(refs) == 0 {
		return nil, nil
	}
	var brl RefRemoteList
	client := http.DefaultClient
	var lastErr error
	for _, b := range refs {
		u, err := url.Parse(fmt.Sprintf(JiraRestURL, b.ID))
		if err != nil {
			return nil, err
		}
		resp, err := client.Get(u.String())
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("jira server responded with %d", resp.StatusCode)
			continue
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("unable to get body contents: %v", err)
			continue
		}
		resp.Body.Close()
		var jrb JiraRemoteRef
		if err := json.Unmarshal(data, &jrb); err != nil {
			lastErr = fmt.Errorf("unable to parse issue list: %v", err)
			continue
		}

		refInfo := convertJiraToRefRemoteInfo(jrb)
		brl.Refs = append(brl.Refs, refInfo)
	}

	if len(brl.Refs) == 0 && lastErr != nil {
		return nil, lastErr
	}

	return &brl, nil
}

func convertJiraToRefRemoteInfo(jrb JiraRemoteRef) RefRemoteInfo {
	var bri RefRemoteInfo
	bri.ID = jrb.Key
	bri.Priority = jrb.Fields.Priority.Name
	bri.Summary = jrb.Fields.Summary
	bri.Status = jrb.Fields.Status.Name
	bri.Source = Jira
	return bri
}

func retrieveRefsBugzila(bugs []Ref) (*RefRemoteList, error) {
	if len(bugs) == 0 {
		return nil, nil
	}
	u, err := url.Parse(bugzillaRestURL)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	q := url.Values{}
	for _, b := range bugs {
		q.Add("id", b.ID)
	}
	u.RawQuery = q.Encode()
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(u.String())
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("server responded with %d", resp.StatusCode)
			continue
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("unable to get body contents: %v", err)
			continue
		}
		resp.Body.Close()
		var bugList RefRemoteList
		if err := json.Unmarshal(data, &bugList); err != nil {
			lastErr = fmt.Errorf("unable to parse bug list: %v", err)
			continue
		}

		for i := range bugList.Refs {
			bugList.Refs[i].Source = Bugzilla
		}

		return &bugList, nil
	}
	return nil, lastErr
}

// PrintRefs prints refs in a formatted way by also handling
// its correct url whether it is bugzilla bug or jira reference.
func (b RefList) PrintRefs(out io.Writer) {
	for i, ref := range b.Refs {
		if i == 0 {
			if ref.Source == Bugzilla {
				fmt.Fprintf(out, " [Bug %s](%s)", ref.ID, fmt.Sprintf(bugzillaBrowseURL, ref.ID))
			} else {
				fmt.Fprintf(out, " [%s](%s)", ref.ID, fmt.Sprintf(jiraBrowseURL, ref.ID))
			}
		} else {
			if ref.Source == Bugzilla {
				fmt.Fprintf(out, ", [%s](%s)", ref.ID, fmt.Sprintf(bugzillaBrowseURL, ref.ID))
			} else {
				fmt.Fprintf(out, ", [%s](%s)", ref.ID, fmt.Sprintf(jiraBrowseURL, ref.ID))
			}
		}
	}
	if len(b.Refs) > 0 {
		fmt.Fprintf(out, ":")
	}
}

// GetRefsForSource returns a map of URLs, sorted by ID, for the specified RefType
func (b RefList) GetRefsForSource(source RefType) map[string]string {
	refs := make(map[string]string)
	for _, ref := range b.Refs {
		if ref.Source == source {
			switch ref.Source {
			case Bugzilla:
				refs[ref.ID] = fmt.Sprintf(bugzillaBrowseURL, ref.ID)
			case Jira:
				refs[ref.ID] = fmt.Sprintf(jiraBrowseURL, ref.ID)
			}
		}
	}
	return refs
}

// GetRefList converts Ref map to Ref list after sorting
// according to the ID field.
func GetRefList(b map[string]Ref) []Ref {
	var bs []Ref
	for _, v := range b {
		bs = append(bs, v)
	}

	sort.Slice(bs, func(i, j int) bool {
		return bs[i].ID < bs[j].ID
	})

	return bs
}

// extractRefs parses git log output and extracts bugzilla+jira references
// compatible with Bugzilla or Jira prefixes.
func extractRefs(msg string) (RefList, string) {
	var b RefList
	if msg == "" {
		return b, ""
	}

	bt := Bugzilla
	msg = rePrefix.ReplaceAllString(strings.TrimSpace(msg), "")
	parsedMsg := reBugzillaBug.FindStringSubmatch(msg)
	if parsedMsg == nil {
		parsedMsg = reJiraRef.FindStringSubmatch(msg)
		if parsedMsg == nil {
			return b, msg
		}
		bt = Jira
	}

	refs := parsedMsg[1]
	if bt == Jira {
		refs = strings.TrimSpace(parsedMsg[0])
		refs = strings.TrimSuffix(refs, ":")
		refs = strings.TrimSuffix(refs, "-")
	}

	// trim all the bug/jira refs from the commit message
	msg = msg[len(parsedMsg[0]):]

	for _, part := range strings.Split(refs, ",") {
		for _, subpart := range strings.Split(part, " ") {
			subpart = strings.TrimSpace(subpart)
			if len(subpart) == 0 {
				continue
			}
			b.Refs = append(b.Refs, Ref{
				ID:     subpart,
				Source: bt,
			})
		}
	}

	sort.Slice(b.Refs, func(i, j int) bool {
		return b.Refs[i].ID < b.Refs[j].ID
	})

	return b, msg
}
