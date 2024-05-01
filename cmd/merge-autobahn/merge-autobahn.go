package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"text/template"
)

// TestCase represents the structure of a test case from the JSON.
type TestCase struct {
	Behavior        string `json:"behavior"`
	BehaviorClose   string `json:"behaviorClose"`
	Duration        int    `json:"duration"`
	RemoteCloseCode int    `json:"remoteCloseCode"`
	ReportFile      string `json:"reportfile"`
}

// TestSuite represents the structure of the entire test suite.
type TestSuite struct {
	TestCases map[string]TestCase `json:"non-tls"`
}

// GroupTitle represents the title for a group of test cases.
type GroupTitle struct {
	Title string `json:"title"`
}

// findGroupTitle attempts to find the most specific group title for a given case ID.
func findGroupTitle(caseID string, groups map[string]GroupTitle) *GroupTitle {
	groupID := strings.Split(caseID, ".")[0]
	if title, ok := groups[groupID]; ok {
		return &title
	}
	return nil
}

func main() {
	// Read the JSON file
	jsonData, err := ioutil.ReadFile("index.json")
	if err != nil {
		fmt.Printf("Error reading JSON file: %s\n", err)
		return
	}

	// Unmarshal the JSON data into a TestSuite struct
	var suite TestSuite
	err = json.Unmarshal(jsonData, &suite)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %s\n", err)
		return
	}

	// Define the group titles based on the provided information
	groupTitles := map[string]GroupTitle{
		"1": {"Framing"},
		"2": {"Pings/Pongs"},
		// ... Add other group titles as needed
	}

	// Sort the test case IDs
	var sortedCaseIDs []string
	for caseID := range suite.TestCases {
		sortedCaseIDs = append(sortedCaseIDs, caseID)
	}
	sort.Strings(sortedCaseIDs)

	// Define an HTML template for the output with inline CSS
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Autobahn Testsuite Report</title>
	<style>
		body { font-family: Arial, sans-serif; }
		table { width: 100%; border-collapse: collapse; }
		th, td { border: 1px solid #ddd; padding: 8px; }
		th { background-color: #f2f2f2; }
		tr:nth-child(even) { background-color: #f9f9f9; }
		h1, h2 { color: #333; }
	</style>
</head>
<body>
	<h1>Autobahn Testsuite Report</h1>
	{{$printedTitle := false}}
	{{range $caseID := .CaseIDs}}
		{{with $group := findGroupTitle $caseID $.GroupTitles}}
            {{if .Title}}
                {{if not $printedTitle}}
                    <h1>{{getGroupTitle $.GroupTitles "1"}}</h1>
                    {{$printedTitle = true}}
                {{end}}
                <h2>{{.Title}}</h2>
            {{end}}
        {{end}}
		{{with $testCase := index $.Suite.TestCases $caseID}}
			<table>
				<thead>
					<tr>
						<th>ID</th>
						<th>Behavior</th>
						<th>Close Behavior</th>
						<th>Duration</th>
						<th>Close Code</th>
						<th>Report File</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>{{$caseID}}</td>
						<td>{{$testCase.Behavior}}</td>
						<td>{{$testCase.BehaviorClose}}</td>
						<td>{{$testCase.Duration}}</td>
						<td>{{$testCase.RemoteCloseCode}}</td>
						<td>{{$testCase.ReportFile}}</td>
					</tr>
				</tbody>
			</table>
		{{end}}
	{{else}}
		<p>No test cases found.</p>
	{{end}}
</body>
</html>
`
	// Create an HTML template
	t := template.Must(template.New("testsuite").Funcs(template.FuncMap{
		"findGroupTitle": func(caseID string, groups map[string]GroupTitle) *GroupTitle {
			return findGroupTitle(caseID, groups)
		},
		"index": func(m map[string]TestCase, key string) TestCase {
			if v, ok := m[key]; ok {
				return v
			}
			return TestCase{}
		},
		"getGroupTitle": func(m map[string]GroupTitle, key string) string {
			if v, ok := m[key]; ok {
				return v.Title
			}
			return ""
		},
	}).Parse(htmlTemplate))

	// Execute the template and write to an HTML file
	var htmlData bytes.Buffer
	err = t.ExecuteTemplate(&htmlData, "testsuite", struct {
		CaseIDs     []string
		Suite       TestSuite
		GroupTitles map[string]GroupTitle
	}{sortedCaseIDs, suite, groupTitles})
	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		return
	}

	err = ioutil.WriteFile("index.html", htmlData.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %s\n", err)
	}
}
