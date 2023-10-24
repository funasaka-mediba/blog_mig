package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: ./program csv_file_path input_xml_file_path output_xml_file_path presence_report_path")
		return
	}

	csvPath := os.Args[1]
	xmlInputPath := os.Args[2]
	xmlOutputPath := os.Args[3]
	reportPath := os.Args[4]

	// CSVファイルを開く
	fmt.Println("Opening CSV file...")
	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("Error opening CSV:", err)
		return
	}
	defer csvFile.Close()

	// CSVを読み込む
	fmt.Println("Reading from CSV...")
	r := csv.NewReader(csvFile)
	replacements := make(map[string]string)
	presenceMap := make(map[string]int)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV:", err)
			return
		}
		replacements[record[0]] = record[1]
		presenceMap[record[0]] = 0
	}

	// XML入力ファイルを開く
	fmt.Println("Opening XML input file...")
	xmlInputFile, err := os.Open(xmlInputPath)
	if err != nil {
		fmt.Println("Error opening input XML:", err)
		return
	}
	defer xmlInputFile.Close()

	// XML出力ファイルを作成
	fmt.Println("Creating XML output file...")
	xmlOutputFile, err := os.Create(xmlOutputPath)
	if err != nil {
		fmt.Println("Error creating output XML:", err)
		return
	}
	defer xmlOutputFile.Close()

	writer := bufio.NewWriter(xmlOutputFile)
	scanner := bufio.NewScanner(xmlInputFile)
	count := 0

	fmt.Println("Processing XML file...")
	for scanner.Scan() {
		line := scanner.Text()
		for oldURL := range replacements {
			if strings.Contains(line, oldURL) {
				presenceMap[oldURL] = 1
				count++
			}
		}
		for oldURL, newURL := range replacements {
			line = strings.ReplaceAll(line, oldURL, newURL)
		}
		writer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input XML:", err)
	}
	writer.Flush()

	// Create a report
	fmt.Println("Generating report...")
	reportFile, err := os.Create(reportPath)
	if err != nil {
		fmt.Println("Error creating report file:", err)
		return
	}
	defer reportFile.Close()

	reportWriter := csv.NewWriter(reportFile)
	reportWriter.Write([]string{"該当", "変換前URL", "変換後URL"})

	for oldURL, newURL := range replacements {
		reportWriter.Write([]string{fmt.Sprintf("%d", presenceMap[oldURL]), oldURL, newURL})
	}
	reportWriter.Flush()

	fmt.Printf("Total of %d matching URLs found.\n", count)
	fmt.Println("Process completed!")
}
