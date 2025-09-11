package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Table struct {
	Headers []string
	Values  [][]any
}

func NewTable(headers []string) Table {
	return Table{
		Headers: headers,
		Values:  [][]any{},
	}
}

func (t *Table) AddRow(values ...any) {
	var valuesCopy = make([]any, len(values))
	copy(valuesCopy, values)
	t.Values = append(t.Values, valuesCopy)
}

func (t *Table) Render() string {
	var lengths = make([]int, len(t.Headers))
	for i := range t.Headers {
		lengths[i] = len(t.Headers[i]) + 2
	}
	for i := range t.Values {
		for j := range t.Values[i] {
			var formattedValue = formatUuuid(t.Values[i][j])
			if lengths[j] < len(fmt.Sprint(formattedValue))+2 {
				lengths[j] = len(fmt.Sprint(formattedValue)) + 2
			}
		}
	}

	var table = ""

	table += t.Separator(lengths, 0)

	var headers = make([]string, len(t.Headers))
	for i := range t.Headers {
		headers[i] = fmt.Sprintf("%-"+strconv.Itoa(lengths[i])+"s", t.Headers[i])
	}

	table += "┃" + strings.Join(headers, "┃") + "┃\n"

	table += t.Separator(lengths, 1)

	rows := make([]string, len(t.Values))
	for i := range t.Values {
		var row = make([]string, len(t.Headers))
		for j := range t.Values[i] {
			var formattedValue = formatUuuid(t.Values[i][j])
			row[j] = fmt.Sprintf("%-"+strconv.Itoa(lengths[j])+"v", formattedValue)
		}
		rows[i] = "┃" + strings.Join(row, "┃") + "┃\n"
	}

	table += strings.Join(rows, t.Separator(lengths, 1))

	table += t.Separator(lengths, 2)
	return table
}

func (t *Table) Separator(n []int, position int) string {
	var start, mid, end string
	if position == 0 {
		start = "┏"
		mid = "┳"
		end = "┓"
	} else if position == 1 {
		start = "┣"
		mid = "╋"
		end = "┫"
	} else {
		start = "┗"
		mid = "┻"
		end = "┛"
	}
	var segments = make([]string, len(n))
	for i := range n {
		segments[i] = strings.Repeat("━", n[i])
	}

	return start + strings.Join(segments, mid) + end + "\n"
}

func formatUuuid(s any) any {
	uuidRegex := `^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`
	st, ok := s.(string)
	if !ok {
		return s
	}
	matched, err := regexp.MatchString(uuidRegex, st)
	if err != nil || !matched {
		return s
	}
	return st[:4] + "..." + st[len(st)-4:]
}
