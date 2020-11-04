package merge

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"

	"github.com/unikumAB/logmerge/formats"
)

func Merge(outputFileName string, inputFiles []string) {
	parser, err := formats.NewVCombinedParser()
	if err != nil {
		log.Fatalf("Cannot create parser: %v", err)
	}
	wg := &sync.WaitGroup{}
	outChan := writeFileLines(outputFileName, wg)
	var inChans []<-chan formats.Line
	for _, arg := range inputFiles {
		inChans = append(inChans, readfile(arg, parser, wg))
	}
	lines := make([]*formats.Line, len(inChans))
	var line *formats.Line
	for {
		lines, inChans = fillLines(inChans, lines)
		line, lines = popOldestLine(lines)
		if line == nil {
			break
		}
		outChan <- line.Source
	}
	close(outChan)
	wg.Wait()
}

func popOldestLine(lines []*formats.Line) (*formats.Line, []*formats.Line) {
	oldestIndex := 0
	var oldestLine *formats.Line
	for i := range lines {
		line := lines[i]
		if oldestLine == nil {
			oldestIndex = i
			oldestLine = line
			continue
		}
		if line == nil {
			continue
		}
		if oldestLine.Timestamp.After(line.Timestamp) {
			oldestIndex = i
			oldestLine = line
		}
	}
	lines[oldestIndex] = nil
	return oldestLine, lines
}

func fillLines(inChans []<-chan formats.Line, lines []*formats.Line) ([]*formats.Line, []<-chan formats.Line) {
	for i, inChan := range inChans {
		if lines[i] == nil {
			line, ok := <-inChan
			if !ok {
				lines[i] = nil
				inChans = append(inChans[:i], inChans[i+1:]...)
				continue
			}
			lines[i] = &line
		}

	}
	return lines, inChans
}

func readfile(arg string, parser formats.LineReader, wg *sync.WaitGroup) <-chan formats.Line {
	inChan := make(chan formats.Line)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(inChan)
		inputFile, err := os.Open(arg)
		if err != nil {
			log.Printf("Failed to open input file %v: %v", arg, err)
			return
		}
		defer checkedClose(inputFile)
		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			text := scanner.Text()
			line, err := parser.ParseLine(text)
			if err != nil {
				log.Printf("Failed to parse line: %v", err)
			}
			inChan <- line
		}
	}()
	return inChan
}

func checkedClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Fatalf("Cannot close normally: %s", err)
	}
}

func writeFileLines(outputFileName string, wg *sync.WaitGroup) chan<- string {
	outChan := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			log.Fatalf("Failed to open output file %v: %s", outputFileName, err)
		}
		defer checkedClose(outputFile)
		writer := bufio.NewWriter(outputFile)
		for line := range outChan {
			count, err := writer.WriteString(line + "\n")
			if err != nil {
				log.Printf("Wrote %d chars but should have more: %v", count, err)
			}
		}
		err = writer.Flush()
		if err != nil {
			log.Fatalf("Failed to Flush output: %s", err)
		}
	}()
	return outChan
}
