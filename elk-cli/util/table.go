package util

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Table struct {
	TableName  string
	Columns    []Column
	Rows       []Row
	Space      int
	TotalWidth int
	isFitted   bool
}

// table types

type TableInt = int
type TableEnum struct {
	Values []struct {
		Name  string
		Value int
		Color int
	}
}
type TableString = string

type Column struct {
	Name  string
	Width int //in spaces
}

type Row struct {
	Values []string
}

func NewTable(name string) *Table {
	return &Table{TableName: name, Space: 2, isFitted: false}
}

func NewTableWithCSV(filename string) (*Table, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	t := NewTable(filename)
	t.InitColumns(strings.Split(lines[0], ","))
	for _, line := range lines[1:] {
		t.AddRow(strings.Split(line, ","))
	}
	return t, nil
}

func (t *Table) Clone(name string) *Table {
	newTable := NewTable(name)
	newTable.Columns = slices.Clone(t.Columns)
	newTable.Rows = slices.Clone(t.Rows)
	newTable.Space = t.Space
	newTable.TotalWidth = t.TotalWidth

	return newTable
}

func (t *Table) InitColumns(names []string) {
	for _, name := range names {
		t.AddColumn(name, 0)
	}
}

func (t *Table) GetColumnIdex(colName string) int {
	for i, col := range t.Columns {
		if col.Name == colName {
			return i
		}
	}
	return -1

}

func (t *Table) AddColumnPos(name string, width int, pos int) *Table {
	if pos == 0 {
		t.Columns = append([]Column{{name, width}}, t.Columns...)
		return t
	}
	t.Columns = append(t.Columns[:pos], Column{name, width})
	t.Columns = append(t.Columns, t.Columns[pos:]...)
	return t
}

func (t *Table) AddColumn(name string, width int) *Table {
	t.Columns = append(t.Columns, Column{name, width})
	return t
}

func (t *Table) AddRow(newValues []string) *Table {
	// check the type of new value
	if len(newValues) != len(t.Columns) {
		log.Fatalf("Invalid number of values for row %v", newValues)
	}
	t.Rows = append(t.Rows, Row{newValues})
	return t
}

func (t *Table) RemoveRow(rowIdx int) *Table {
	t.Rows = slices.Delete(t.Rows, rowIdx, rowIdx+1)
	return t
}

func (t *Table) RemoveColumnByID(colIdx int) *Table {
	for i, row := range t.Rows {
		t.Rows[i].Values = slices.Delete(row.Values, colIdx, colIdx+1)
	}
	t.Columns = slices.Delete(t.Columns, colIdx, colIdx+1)
	return t
}

func (t *Table) RemoveColumnByName(colName string) *Table {
	colIdx := t.GetColumnIdex(colName)
	return t.RemoveColumnByID(colIdx)
}

func (t *Table) Filter(colName string, value string) *Table {
	colIdx := t.GetColumnIdex(colName)
	cloneT := t.Clone("filtered_" + colName + "_" + value + "_" + t.TableName)
	if colIdx == -1 {
		log.Fatalf("Column %s not found", colName)
	}
	for i, row := range cloneT.Rows {
		if row.Values[colIdx] != value {
			cloneT.Rows = slices.Delete(cloneT.Rows, i, i+1)
		}
	}
	return cloneT
}

func (t *Table) Sort(isAsc bool, columnName ...string) error {
	targetColIdx := []int{}

	for _, colName := range columnName {
		i := t.GetColumnIdex(colName)
		if i == -1 {
			return fmt.Errorf("Column %s not found", colName)
		}
		targetColIdx = append(targetColIdx, i)
	}
	if len(targetColIdx) == 0 {
		return fmt.Errorf("Column %s not found", columnName)
	}
	sort.Slice(t.Rows, func(i, j int) bool {
		if isAsc {
			for _, colIdx := range targetColIdx {
				if t.Rows[i].Values[colIdx] < t.Rows[j].Values[colIdx] {
					return true
				} else if t.Rows[i].Values[colIdx] > t.Rows[j].Values[colIdx] {
					return false
				}
			}
			return false
		} else {
			for _, colIdx := range targetColIdx {
				if t.Rows[i].Values[colIdx] > t.Rows[j].Values[colIdx] {
					return true
				} else if t.Rows[i].Values[colIdx] < t.Rows[j].Values[colIdx] {
					return false
				}
			}
			return false
		}
	})
	return nil
}

func (t *Table) Serialize() {
	t.AddColumnPos("ID", 2, 0)
	for i, row := range t.Rows {
		t.Rows[i].Values = append([]string{strconv.Itoa(i + 1)}, row.Values...)
	}
}

func (t *Table) fit() {
	if t.isFitted {
		return
	}
	w := 0
	for i := range t.Columns {
		col := &t.Columns[i]
		col.Width = len(col.Name)
		for _, row := range t.Rows {
			if len(row.Values[i]) > col.Width {
				col.Width = len(row.Values[i])
			}
		}
		w += col.Width
	}
	//total width = sum of column widths + spaces + borders + padding
	t.TotalWidth = w + (len(t.Columns) * t.Space) + (len(t.Columns) + 1) + (len(t.Columns))
	t.isFitted = true
}

func (t Table) Print() {
	t.fit()
	fmt.Println("Table Name:", t.TableName)

	fmt.Printf("%s\n", strings.Repeat("-", t.TotalWidth))
	// print header
	for i, col := range t.Columns {
		if i == 0 {
			fmt.Printf("| %-*s%*s|", col.Width, col.Name, t.Space, "")
		} else {
			fmt.Printf(" %-*s%*s|", col.Width, col.Name, t.Space, "")
		}
	}
	fmt.Println()
	fmt.Printf("%s\n", strings.Repeat("-", t.TotalWidth))
	// print rows
	for _, row := range t.Rows {
		for i, col := range row.Values {
			if i == 0 {
				fmt.Printf("| %-*s%*s|", t.Columns[i].Width, col, t.Space, "")
			} else {
				fmt.Printf(" %-*s%*s|", t.Columns[i].Width, col, t.Space, "")
			}
		}
		fmt.Println()
	}
	fmt.Printf("%s\n", strings.Repeat("-", t.TotalWidth))
}
