package merge

import (
	"bufio"

	"github.com/UnikumAB/logmerge/compress"
	"github.com/UnikumAB/logmerge/utils"
	gzip "github.com/klauspost/pgzip"
	//"compress/gzip"
	"io"
	"log"
	"os"
	"sync"

	"github.com/UnikumAB/logmerge/formats"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

func Merge(outputFileName string, inputFiles []string) {
	parser, err := formats.NewVCombinedParser()
	if err != nil {
		log.Fatalf("Cannot create parser: %v", err)
	}
	wg := &sync.WaitGroup{}
	p := mpb.New(mpb.WithWaitGroup(wg))
	outChan := writeFileLines(outputFileName, wg)
	var inChans []<-chan formats.Line
	for _, filename := range inputFiles {
		bar := p.AddBar(0,
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(filename),
				// decor.DSyncWidth bit enables column width synchronization
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				// replace ETA decorator with "done" message, OnComplete event
				decor.OnComplete(
					// ETA decorator with ewma age of 60
					decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
				),
			),
			mpb.BarRemoveOnComplete(),
		)

		inChans = append(inChans, readfile(filename, parser, wg, bar))
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

func readfile(filename string, parser formats.LineReader, wg *sync.WaitGroup, bar *mpb.Bar) <-chan formats.Line {
	inChan := make(chan formats.Line)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(inChan)
		isGzip, size, err := compress.DetectGzip(filename)
		bar.SetTotal(size, false)
		if err != nil {
			log.Printf("Failed to detect content type: %v", err)
			return
		}
		inputFile, err := os.Open(filename)
		if err != nil {
			log.Printf("Failed to open input file %v: %v", filename, err)
			return
		}
		defer utils.CheckedClose(inputFile)
		var inputReader io.Reader = inputFile
		if isGzip {
			reader, err := gzip.NewReader(inputFile)
			if err != nil {
				log.Printf("Failed to create gzip reader for %v: %v", filename, err)
			}
			defer utils.CheckedClose(reader)
			inputReader = reader
		}

		scanner := bufio.NewScanner(inputReader)
		for scanner.Scan() {
			text := scanner.Text()
			bar.IncrBy(len([]byte(text)))
			line, err := parser.ParseLine(text)
			if err != nil {
				log.Printf("Failed to parse line: %v", err)
			}
			inChan <- line
		}
		bar.SetTotal(bar.Current(), true)
	}()
	return inChan
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
		defer utils.CheckedClose(outputFile)
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
