package release

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

type BugType int

const (
	Bugzilla BugType = iota
	Jira
)

const (
	maxRetries = 2

	bugzillaBrowseURL = "https://bugzilla.redhat.com/show_bug.cgi?id=%d"
	bugzillaRestURL   = "https://bugzilla.redhat.com/rest/bug"
	jiraBrowseURL     = "https://issues.redhat.com/browse/OCPBUGS-%d"
	JiraRestURL       = "https://issues.redhat.com/rest/api/latest/issue/OCPBUGS-%d"
)

var (
	reBugzillaBug = regexp.MustCompile(`^[B|b]ug(?:s)?\s([\d\s,]+)(?:-|:)\s*`)
	reJiraBug     = regexp.MustCompile(`^OCPBUGS-([\d\s,]+)(?:-|:)\s*`)
)

// BugList stores Bug list
type BugList struct {
	Bugs []Bug
}

// Bug stores referenced bug ids and
// source of bug type like Bugzilla or Jira.
type Bug struct {
	ID     int
	Source BugType
}

// BugRemoteList stores BugRemoteInfo list
type BugRemoteList struct {
	Bugs []BugRemoteInfo `json:"bugs"`
}

// BugRemoteInfo stores the detail information of bug retrieved
// from remote url like bugzilla or jira.
type BugRemoteInfo struct {
	ID       int     `json:"id"`
	Status   string  `json:"status"`
	Priority string  `json:"priority"`
	Summary  string  `json:"summary"`
	Source   BugType `json:"source"`
}

type JiraRemoteBug struct {
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

// RetrieveBugs retrieves bug details and fills BugRemoteList
// by sending request to Bugzilla and Jira according to the source type of bug.
func RetrieveBugs(bugs []Bug) (*BugRemoteList, error) {
	var brl BugRemoteList
	var bugzilla, jira []Bug
	for _, v := range bugs {
		if v.Source == Bugzilla {
			bugzilla = append(bugzilla, v)
		} else {
			jira = append(jira, v)
		}
	}

	var lastErr error
	bz, err := retrieveBugsBugzilla(bugzilla)
	if err != nil {
		lastErr = err
	}
	if bz != nil {
		brl.Bugs = append(brl.Bugs, bz.Bugs...)
	}

	jr, err := retrieveBugsJira(jira)
	if err != nil {
		lastErr = err
	}
	if jr != nil {
		brl.Bugs = append(brl.Bugs, jr.Bugs...)
	}

	return &brl, lastErr
}

func retrieveBugsJira(bugs []Bug) (*BugRemoteList, error) {
	if len(bugs) == 0 {
		return nil, nil
	}
	var brl BugRemoteList
	client := http.DefaultClient
	var lastErr error
	for _, b := range bugs {
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
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("unable to get body contents: %v", err)
			continue
		}
		resp.Body.Close()
		var jrb JiraRemoteBug
		if err := json.Unmarshal(data, &jrb); err != nil {
			lastErr = fmt.Errorf("unable to parse bug list: %v", err)
			continue
		}

		bugInfo := convertJiraToBugRemoteInfo(jrb)
		brl.Bugs = append(brl.Bugs, bugInfo)
	}

	if len(brl.Bugs) == 0 && lastErr != nil {
		return nil, lastErr
	}

	return &brl, nil
}

func convertJiraToBugRemoteInfo(jrb JiraRemoteBug) BugRemoteInfo {
	var bri BugRemoteInfo
	bri.ID, _ = strconv.Atoi(strings.Replace(jrb.Key, "OCPBUGS-", "", 1))
	bri.Priority = jrb.Fields.Priority.Name
	bri.Summary = jrb.Fields.Summary
	bri.Status = jrb.Fields.Status.Name
	bri.Source = Jira
	return bri
}

func retrieveBugsBugzilla(bugs []Bug) (*BugRemoteList, error) {
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
		q.Add("id", strconv.Itoa(b.ID))
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
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("unable to get body contents: %v", err)
			continue
		}
		resp.Body.Close()
		var bugList BugRemoteList
		if err := json.Unmarshal(data, &bugList); err != nil {
			lastErr = fmt.Errorf("unable to parse bug list: %v", err)
			continue
		}

		for i := range bugList.Bugs {
			bugList.Bugs[i].Source = Bugzilla
		}

		return &bugList, nil
	}
	return nil, lastErr
}

// PrintBugs prints bugs in a formatted way by also handling
// its correct url whether it is bugzilla bug or jira bug.
func (b BugList) PrintBugs(out io.Writer) {
	for i, bug := range b.Bugs {
		if i == 0 {
			if bug.Source == Bugzilla {
				fmt.Fprintf(out, " [Bug %d](%s)", bug.ID, fmt.Sprintf(bugzillaBrowseURL, bug.ID))
			} else {
				fmt.Fprintf(out, " [OCPBUGS-%d](%s)", bug.ID, fmt.Sprintf(jiraBrowseURL, bug.ID))
			}
		} else {
			if bug.Source == Bugzilla {
				fmt.Fprintf(out, ", [%d](%s)", bug.ID, fmt.Sprintf(bugzillaBrowseURL, bug.ID))
			} else {
				fmt.Fprintf(out, ", [%d](%s)", bug.ID, fmt.Sprintf(jiraBrowseURL, bug.ID))
			}
		}
	}
	if len(b.Bugs) > 0 {
		fmt.Fprintf(out, ":")
	}
}

// GetBugList converts Bug map to bug list after sorting
// according to the ID field.
func GetBugList(b map[string]Bug) []Bug {
	var bs []Bug
	for _, v := range b {
		bs = append(bs, v)
	}

	sort.Slice(bs, func(i, j int) bool {
		return bs[i].ID < bs[j].ID
	})

	return bs
}

func generateBugKey(bt BugType, id int) string {
	if bt == Bugzilla {
		return fmt.Sprintf("bugzilla-%d", id)
	}

	return fmt.Sprintf("jira-%d", id)
}

// extractBugs parses git log output and extracts bug list
// compatible with Bugzilla or Jira prefixes.
func extractBugs(msg string) (BugList, string) {
	var b BugList
	if msg == "" {
		return b, ""
	}

	bt := Bugzilla
	msg = rePrefix.ReplaceAllString(strings.TrimSpace(msg), "")
	parsedMsg := reBugzillaBug.FindStringSubmatch(msg)
	if parsedMsg == nil {
		parsedMsg = reJiraBug.FindStringSubmatch(msg)
		if parsedMsg == nil {
			return b, msg
		}
		bt = Jira
	}

	msg = msg[len(parsedMsg[0]):]
	for _, part := range strings.Split(parsedMsg[1], ",") {
		for _, subpart := range strings.Split(part, " ") {
			subpart = strings.TrimSpace(subpart)
			if len(subpart) == 0 {
				continue
			}
			bug, err := strconv.Atoi(subpart)
			if err != nil {
				klog.V(5).Infof("unable to parse numbers from %q: %v", part, err)
				continue
			}
			b.Bugs = append(b.Bugs, Bug{
				ID:     bug,
				Source: bt,
			})
		}
	}

	sort.Slice(b.Bugs, func(i, j int) bool {
		return b.Bugs[i].ID < b.Bugs[j].ID
	})

	return b, msg
}
