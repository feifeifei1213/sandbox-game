package service

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"regexp"
	"strconv"
	"strings"

	"sandbox-game/internal/enum"
)

var (
	orderSheetNamePattern = regexp.MustCompile(`^([1-9]\d*)年(本地|区域|全国|全球)$`)
	cellRefPattern        = regexp.MustCompile(`^([A-Z]+)(\d+)$`)
)

type ParsedOrderWorkbook struct {
	Orders   []ParsedOrderCard `json:"orders"`
	Warnings []string          `json:"warnings"`
}

type ParsedOrderCard struct {
	YearNo          int     `json:"yearNo"`
	MarketCode      string  `json:"marketCode"`
	MarketName      string  `json:"marketName"`
	OrderType       string  `json:"orderType"`
	OrderTypeName   string  `json:"orderTypeName"`
	OrderAmount     float64 `json:"orderAmount"`
	OrderQuantity   float64 `json:"orderQuantity"`
	UnitPrice       float64 `json:"unitPrice"`
	AccountTerm     int     `json:"accountTerm"`
	SourceSheetName string  `json:"sourceSheetName"`
	SourceCell      string  `json:"sourceCell"`
	SourceRowIndex  int     `json:"sourceRowIndex"`
}

type xlsxWorkbook struct {
	Sheets []xlsxSheet `xml:"sheets>sheet"`
}

type xlsxSheet struct {
	Name    string `xml:"name,attr"`
	SheetID string `xml:"sheetId,attr"`
	RelID   string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr"`
}

type xlsxRelationships struct {
	Relationships []xlsxRelationship `xml:"Relationship"`
}

type xlsxRelationship struct {
	ID     string `xml:"Id,attr"`
	Target string `xml:"Target,attr"`
}

type xlsxSharedStrings struct {
	Items []xlsxSharedStringItem `xml:"si"`
}

type xlsxSharedStringItem struct {
	Text string              `xml:"t"`
	Runs []xlsxSharedRunText `xml:"r"`
}

type xlsxSharedRunText struct {
	Text string `xml:"t"`
}

type xlsxWorksheet struct {
	Rows []xlsxRow `xml:"sheetData>row"`
}

type xlsxRow struct {
	Cells []xlsxCell `xml:"c"`
}

type xlsxCell struct {
	Ref     string `xml:"r,attr"`
	Type    string `xml:"t,attr"`
	Value   string `xml:"v"`
	Formula string `xml:"f"`
}

type xlsxSheetMeta struct {
	Name   string
	Target string
}

func ParseOrderExcel(data []byte) (*ParsedOrderWorkbook, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open xlsx zip: %w", err)
	}

	entries := mapZipEntries(reader)
	sharedStrings, err := readSharedStrings(entries)
	if err != nil {
		return nil, err
	}
	sheets, err := readWorkbookSheets(entries)
	if err != nil {
		return nil, err
	}

	result := &ParsedOrderWorkbook{
		Orders:   make([]ParsedOrderCard, 0, 256),
		Warnings: make([]string, 0),
	}
	for _, sheet := range sheets {
		yearNo, marketCode, ok := parseOrderSheetName(sheet.Name)
		if !ok {
			continue
		}
		cellValues, err := readWorksheetCells(entries, sheet.Target, sharedStrings)
		if err != nil {
			return nil, fmt.Errorf("read sheet %s: %w", sheet.Name, err)
		}
		orders := extractOrderCardsFromSheet(sheet.Name, yearNo, marketCode, cellValues)
		result.Orders = append(result.Orders, orders...)
	}

	if len(result.Orders) == 0 {
		result.Warnings = append(result.Warnings, "未从年度市场 sheet 中解析到有效订单卡片")
	}
	return result, nil
}

func mapZipEntries(reader *zip.Reader) map[string]*zip.File {
	entries := make(map[string]*zip.File, len(reader.File))
	for _, file := range reader.File {
		entries[file.Name] = file
	}
	return entries
}

func readSharedStrings(entries map[string]*zip.File) ([]string, error) {
	file := entries["xl/sharedStrings.xml"]
	if file == nil {
		return nil, nil
	}
	raw, err := readZipFile(file)
	if err != nil {
		return nil, fmt.Errorf("read shared strings: %w", err)
	}
	var parsed xlsxSharedStrings
	if err := xml.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("parse shared strings: %w", err)
	}
	values := make([]string, 0, len(parsed.Items))
	for _, item := range parsed.Items {
		if item.Text != "" {
			values = append(values, item.Text)
			continue
		}
		var builder strings.Builder
		for _, run := range item.Runs {
			builder.WriteString(run.Text)
		}
		values = append(values, builder.String())
	}
	return values, nil
}

func readWorkbookSheets(entries map[string]*zip.File) ([]xlsxSheetMeta, error) {
	workbookFile := entries["xl/workbook.xml"]
	relsFile := entries["xl/_rels/workbook.xml.rels"]
	if workbookFile == nil || relsFile == nil {
		return nil, fmt.Errorf("xlsx missing workbook metadata")
	}

	workbookRaw, err := readZipFile(workbookFile)
	if err != nil {
		return nil, fmt.Errorf("read workbook: %w", err)
	}
	relsRaw, err := readZipFile(relsFile)
	if err != nil {
		return nil, fmt.Errorf("read workbook relationships: %w", err)
	}

	var workbook xlsxWorkbook
	if err := xml.Unmarshal(workbookRaw, &workbook); err != nil {
		return nil, fmt.Errorf("parse workbook: %w", err)
	}
	var rels xlsxRelationships
	if err := xml.Unmarshal(relsRaw, &rels); err != nil {
		return nil, fmt.Errorf("parse workbook relationships: %w", err)
	}
	relTargets := make(map[string]string, len(rels.Relationships))
	for _, rel := range rels.Relationships {
		relTargets[rel.ID] = normalizeWorkbookTarget(rel.Target)
	}

	sheets := make([]xlsxSheetMeta, 0, len(workbook.Sheets))
	for _, sheet := range workbook.Sheets {
		target := relTargets[sheet.RelID]
		if target == "" {
			continue
		}
		sheets = append(sheets, xlsxSheetMeta{Name: sheet.Name, Target: target})
	}
	return sheets, nil
}

func normalizeWorkbookTarget(target string) string {
	target = strings.TrimSpace(target)
	if strings.HasPrefix(target, "/") {
		return strings.TrimPrefix(target, "/")
	}
	if strings.HasPrefix(target, "xl/") {
		return target
	}
	return path.Clean("xl/" + target)
}

func readWorksheetCells(entries map[string]*zip.File, target string, sharedStrings []string) (map[string]string, error) {
	file := entries[target]
	if file == nil {
		return nil, fmt.Errorf("worksheet %s not found", target)
	}
	raw, err := readZipFile(file)
	if err != nil {
		return nil, err
	}
	var sheet xlsxWorksheet
	if err := xml.Unmarshal(raw, &sheet); err != nil {
		return nil, fmt.Errorf("parse worksheet: %w", err)
	}
	values := make(map[string]string)
	for _, row := range sheet.Rows {
		for _, cell := range row.Cells {
			if strings.TrimSpace(cell.Ref) == "" {
				continue
			}
			value := strings.TrimSpace(cell.Value)
			if cell.Type == "s" && value != "" {
				index, err := strconv.Atoi(value)
				if err == nil && index >= 0 && index < len(sharedStrings) {
					value = sharedStrings[index]
				}
			}
			if value == "" {
				continue
			}
			values[cell.Ref] = strings.TrimSpace(value)
		}
	}
	return values, nil
}

func extractOrderCardsFromSheet(sheetName string, yearNo int, marketCode string, cells map[string]string) []ParsedOrderCard {
	typeColumns := []string{"B", "H", "N", "T", "Z"}
	valueColumns := map[string]string{
		"B": "C",
		"H": "I",
		"N": "O",
		"T": "U",
		"Z": "AA",
	}
	orders := make([]ParsedOrderCard, 0)
	rowIndex := 1
	for row := 1; row <= 240; row++ {
		for _, col := range typeColumns {
			orderTypeName := strings.TrimSpace(cells[cellRef(col, row)])
			orderType, ok := mapOrderTypeName(orderTypeName)
			if !ok || !strings.Contains(cells[cellRef(col, row+1)], "总额") {
				continue
			}
			valueCol := valueColumns[col]
			amount, amountOK := parseFiniteNumber(cells[cellRef(valueCol, row+1)])
			quantity, quantityOK := parseFiniteNumber(cells[cellRef(valueCol, row+2)])
			unitPrice, unitPriceOK := parseFiniteNumber(cells[cellRef(valueCol, row+3)])
			accountTerm, termOK := parseAccountTerm(cells[cellRef(valueCol, row+4)])
			if !amountOK || !quantityOK || !unitPriceOK || !termOK || amount <= 0 || quantity <= 0 {
				continue
			}
			orders = append(orders, ParsedOrderCard{
				YearNo:          yearNo,
				MarketCode:      marketCode,
				MarketName:      marketName(marketCode),
				OrderType:       orderType,
				OrderTypeName:   orderTypeName,
				OrderAmount:     amount,
				OrderQuantity:   quantity,
				UnitPrice:       unitPrice,
				AccountTerm:     accountTerm,
				SourceSheetName: sheetName,
				SourceCell:      cellRef(col, row),
				SourceRowIndex:  rowIndex,
			})
			rowIndex++
		}
	}
	return orders
}

func readZipFile(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func parseOrderSheetName(name string) (int, string, bool) {
	matches := orderSheetNamePattern.FindStringSubmatch(strings.TrimSpace(name))
	if len(matches) != 3 {
		return 0, "", false
	}
	yearNo, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", false
	}
	marketCode, ok := mapMarketName(matches[2])
	return yearNo, marketCode, ok
}

func mapMarketName(name string) (string, bool) {
	switch strings.TrimSpace(name) {
	case "本地", "本地市场", "01本地市场":
		return enum.MarketCodeLocal, true
	case "区域", "区域市场", "02区域市场":
		return enum.MarketCodeRegional, true
	case "全国", "全国市场", "03全国市场":
		return enum.MarketCodeNational, true
	case "全球", "全球市场", "04全球市场":
		return enum.MarketCodeGlobal, true
	default:
		return "", false
	}
}

func mapOrderTypeName(name string) (string, bool) {
	switch strings.TrimSpace(name) {
	case "代办过检":
		return enum.OrderTypeAgencyInspection, true
	case "两舱贵宾":
		return enum.OrderTypeTwoCabinVIP, true
	case "商务贵宾":
		return enum.OrderTypeBusinessVIP, true
	case "会员定制":
		return enum.OrderTypeMemberCustom, true
	default:
		return "", false
	}
}

func marketName(code string) string {
	switch code {
	case enum.MarketCodeLocal:
		return "本地市场"
	case enum.MarketCodeRegional:
		return "区域市场"
	case enum.MarketCodeNational:
		return "全国市场"
	case enum.MarketCodeGlobal:
		return "全球市场"
	default:
		return code
	}
}

func orderTypeName(code string) string {
	switch code {
	case enum.OrderTypeAgencyInspection:
		return "代办过检"
	case enum.OrderTypeTwoCabinVIP:
		return "两舱贵宾"
	case enum.OrderTypeBusinessVIP:
		return "商务贵宾"
	case enum.OrderTypeMemberCustom:
		return "会员定制"
	default:
		return code
	}
}

func cellRef(col string, row int) string {
	return fmt.Sprintf("%s%d", col, row)
}

func parseFiniteNumber(value string) (float64, bool) {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "#") {
		return 0, false
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func parseAccountTerm(value string) (int, bool) {
	parsed, ok := parseFiniteNumber(value)
	if !ok {
		return 0, false
	}
	return int(parsed), true
}

func splitCellRef(ref string) (string, int, bool) {
	matches := cellRefPattern.FindStringSubmatch(ref)
	if len(matches) != 3 {
		return "", 0, false
	}
	row, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, false
	}
	return matches[1], row, true
}
