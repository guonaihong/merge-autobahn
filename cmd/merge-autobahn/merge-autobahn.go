package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
)

var htmlTemplate = `<!DOCTYPE html>
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

		/* 添加一个名为 'grey-column' 的类，用于控制第一列的灰色背景 */
        .grey-column td:first-child {
            background-color: #666666;
        }
    </style>

</head>
<body>
    <h1>Autobahn Testsuite Report</h1>
    {{$printedGroupHeader := false}}
	<table>
	<tr>
		<th>ID</th>
		<th>Behavior</th>
		<th>Close Behavior</th>
		<th>Duration</th>
		<th>Close Code</th>
		<th>Report File</th>
	</tr>
    {{range $caseID := .CaseIDs}}
        {{$group := findGroupTitle $caseID $.GroupTitles}}
        {{with $group}}
            {{if .ParentTitle}}
				<tr> 
					<td>{{.ParentTitle}} </td>
				</tr>
            {{end}}
            {{if .Title}}
				<tr> 
					<td> {{.Title}} </td>
				</tr>
            {{end}}
        {{end}}
        {{with $testCase := index $.Suite.TestCases $caseID}}
                <tr class="grey-column">
                    <td>Case {{$caseID}}</td>
                    <td>{{$testCase.Behavior}}</td>
                    <td>{{$testCase.BehaviorClose}}</td>
                    <td>{{$testCase.Duration}}</td>
                    <td>{{$testCase.RemoteCloseCode}}</td>
                    <td>{{$testCase.ReportFile}}</td>
                </tr>
        {{end}}
    {{else}}
        <p>No test cases found.</p>
    {{end}}
	</table>
</body>
</html>
`

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
	Title       string `json:"title"`
	ParentTitle string `json:"parentTitle"`
}

// findGroupTitle tries to find the most specific group title for the given caseID.
func findGroupTitle(caseID string, groups map[string]GroupTitle, status map[string]bool) GroupTitle {
	// Traverse keys in descending order to match the longest prefix.
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	for _, prefix := range keys {
		if strings.HasPrefix(caseID, prefix) {
			if status[prefix] {
				return GroupTitle{}
			}

			status[prefix] = true
			return groups[prefix]
		}
	}
	return GroupTitle{}
}

func versionCompare(version1, version2 string) int {
	i, j := 0, 0
	for i < len(version1) || j < len(version2) {
		var num1, num2 int
		for i < len(version1) && version1[i] != '.' {
			num1 = num1*10 + int(version1[i]-'0')
			i++
		}
		for j < len(version2) && version2[j] != '.' {
			num2 = num2*10 + int(version2[j]-'0')
			j++
		}
		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
		i++
		j++
	}
	return 0
}
func main() {
	// Read the JSON file
	jsonData, err := os.ReadFile("index.json")
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
		"1.1":  {Title: "1.1 Text Messages", ParentTitle: "1 Framing"},
		"1.2":  {Title: "1.2 Binary Messages", ParentTitle: "1 Framing"},
		"2":    {Title: "Pings/Pongs", ParentTitle: ""},
		"3":    {Title: "Reserved Bits", ParentTitle: ""},
		"4.1":  {Title: "4.1 Non-control Opcodes", ParentTitle: "4 Opcodes"},
		"4.2":  {Title: "4.2 Control Opcodes", ParentTitle: "4 Opcodes"},
		"5":    {Title: "Fragmentation", ParentTitle: ""},
		"6.1":  {Title: "6.1 Valid UTF-8 with zero payload fragments", ParentTitle: "6 UTF-8 Handling"},
		"6.2":  {Title: "6.2 Valid UTF-8 unfragmented, fragmented on code-points and within code-points", ParentTitle: "6 UTF-8 Handling"},
		"6.3":  {Title: "6.3 Invalid UTF-8 differently fragmented", ParentTitle: "6 UTF-8 Handling"},
		"6.4":  {Title: "6.4 Fail-fast on invalid UTF-8", ParentTitle: "6 UTF-8 Handling"},
		"6.5":  {Title: "6.5 Some valid UTF-8 sequences", ParentTitle: "6 UTF-8 Handling"},
		"6.6":  {Title: "6.6 All prefixes of a valid UTF-8 string that contains multi-byte code points", ParentTitle: "6 UTF-8 Handling"},
		"6.7":  {Title: "6.7 First possible sequence of a certain length", ParentTitle: "6 UTF-8 Handling"},
		"6.8":  {Title: "6.8 First possible sequence length 5/6 (invalid codepoints)", ParentTitle: "6 UTF-8 Handling"},
		"6.9":  {Title: "6.9 Last possible sequence of a certain length", ParentTitle: "6 UTF-8 Handling"},
		"6.10": {Title: "6.10 Last possible sequence length 4/5/6 (invalid codepoints)", ParentTitle: "6 UTF-8 Handling"},
		"6.11": {Title: "6.11 Other boundary conditions", ParentTitle: "6 UTF-8 Handling"},
		"6.12": {Title: "6.12 Unexpected continuation bytes", ParentTitle: "6 UTF-8 Handling"},
		"6.13": {Title: "6.13 Lonely start characters", ParentTitle: "6 UTF-8 Handling"},
		"6.14": {Title: "6.14 Sequences with last continuation byte missing", ParentTitle: "6 UTF-8 Handling"},
		"6.15": {Title: "6.15 Concatenation of incomplete sequences", ParentTitle: "6 UTF-8 Handling"},
		"6.16": {Title: "6.16 Impossible bytes", ParentTitle: "6 UTF-8 Handling"},
		"6.17": {Title: "6.17 Examples of an overlong ASCII character", ParentTitle: "6 UTF-8 Handling"},
		"6.18": {Title: "6.18 Maximum overlong sequences", ParentTitle: "6 UTF-8 Handling"},
		"6.19": {Title: "6.19 Overlong representation of the NUL character", ParentTitle: "6 UTF-8 Handling"},
		"6.20": {Title: "6.20 Single UTF-16 surrogates", ParentTitle: "6 UTF-8 Handling"},
		"6.21": {Title: "6.21 Paired UTF-16 surrogates", ParentTitle: "6 UTF-8 Handling"},
		"6.22": {Title: "6.22 Non-character code points (valid UTF-8)", ParentTitle: "6 UTF-8 Handling"},
		"6.23": {Title: "6.23 Unicode specials (i.e. replacement char)", ParentTitle: "6 UTF-8 Handling"},
		"7.1":  {Title: "7.1 Basic close behavior (fuzzer initiated)", ParentTitle: "7 Close Handling"},
		"7.3":  {Title: "7.3 Close frame structure: payload length (fuzzer initiated)", ParentTitle: "7 Close Handling"},
		"7.5":  {Title: "7.5 Close frame structure: payload value (fuzzer initiated)", ParentTitle: "7 Close Handling"},
		"7.7":  {Title: "7.7 Close frame structure: valid close codes (fuzzer initiated)", ParentTitle: "7 Close Handling"},
		"7.9":  {Title: "7.9 Close frame structure: invalid close codes (fuzzer initiated)", ParentTitle: "7 Close Handling"},
		"7.13": {Title: "7.13 Informational close information (fuzzer initiated)", ParentTitle: "7 Close Handling"},
		"9.1":  {Title: "9.1 Text Message (increasing size)", ParentTitle: "9 Limits/Performance"},
		"9.2":  {Title: "9.2 Binary Message (increasing size)", ParentTitle: "9 Limits/Performance"},
		"9.3":  {Title: "9.3 Fragmented Text Message (fixed size, increasing fragment size)", ParentTitle: "9 Limits/Performance"},
		"9.4":  {Title: "9.4 Fragmented Binary Message (fixed size, increasing fragment size)", ParentTitle: "9 Limits/Performance"},
		"9.5":  {Title: "9.5 Text Message (fixed size, increasing chop size)", ParentTitle: "9 Limits/Performance"},
		"9.6":  {Title: "9.6 Binary Text Message (fixed size, increasing chop size)", ParentTitle: "9 Limits/Performance"},
		"9.7":  {Title: "9.7 Text Message Roundtrip Time (fixed number, increasing size)", ParentTitle: "9 Limits/Performance"},
		"9.8":  {Title: "9.8 Binary Message Roundtrip Time (fixed number, increasing size)", ParentTitle: "9 Limits/Performance"},
		"10.1": {Title: "10.1 Auto-Fragmentation", ParentTitle: "10 Misc"},
		"12.1": {Title: "12.1 Large JSON data file (utf8, 194056 bytes)", ParentTitle: "12 WebSocket Compression (different payloads)"},
		"12.2": {Title: "12.2 Lena Picture, Bitmap 512x512 bw (binary, 263222 bytes)", ParentTitle: "12 WebSocket Compression (different payloads)"},
		"12.3": {Title: "12.3 Human readable text, Goethe's Faust I (German) (binary, 222218 bytes)", ParentTitle: "12 WebSocket Compression (different payloads)"},
		"12.4": {Title: "12.4 Large HTML file (utf8, 263527 bytes)", ParentTitle: "12 WebSocket Compression (different payloads)"},
		"12.5": {Title: "12.5 A larger PDF (binary, 1042328 bytes)", ParentTitle: "12 WebSocket Compression (different payloads)"},
		"13.1": {Title: "13.1 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(False, 0)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(False, 0)]", ParentTitle: "13 WebSocket Compression (different parameters)"},
		"13.2": {Title: "13.2 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(True, 0)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(True, 0)]", ParentTitle: "13 WebSocket Compression (different parameters)"},
		"13.3": {Title: "13.3 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(False, 9)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(False, 9)]", ParentTitle: "13 WebSocket Compression (different parameters)"},
		"13.4": {Title: "13.4 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(False, 15)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(False, 15)]", ParentTitle: "13 WebSocket Compression (different parameters)"},
		"13.5": {Title: "13.5 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(True, 9)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(True, 9)]", ParentTitle: "13 WebSocket Compression (different parameters)"},
		"13.6": {Title: "13.6 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(True, 15)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(True, 15)]", ParentTitle: "13 WebSocket Compression (different parameters)"},
		"13.7": {Title: "13.7 Large JSON data file (utf8, 194056 bytes) - client offers (requestNoContextTakeover, requestMaxWindowBits): [(True, 9), (True, 0), (False, 0)] / server accept (requestNoContextTakeover, requestMaxWindowBits): [(True, 9), (True, 0), (False, 0)]", ParentTitle: "13 WebSocket Compression (different parameters)"},
	}

	// Sort the test case IDs
	var sortedCaseIDs []string
	for caseID := range suite.TestCases {
		sortedCaseIDs = append(sortedCaseIDs, caseID)
	}
	// sort.Strings(sortedCaseIDs)
	sort.Slice(sortedCaseIDs, func(i, j int) bool {
		return versionCompare(sortedCaseIDs[i], sortedCaseIDs[j]) == -1
	})

	status := make(map[string]bool)
	// Define an HTML template for the output with inline CSS
	// Create an HTML template
	t := template.Must(template.New("testsuite").Funcs(template.FuncMap{
		"findGroupTitle": func(caseID string, groups map[string]GroupTitle) GroupTitle {
			return findGroupTitle(caseID, groups, status)
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

	err = os.WriteFile("index.html", htmlData.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %s\n", err)
	}
}
