package content

import (
	"context"
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/epub"
	"robbykansas/another-novel-scraper/cmd/models"
	"robbykansas/another-novel-scraper/cmd/testbenchmark"
	"robbykansas/another-novel-scraper/cmd/ui/progressbar"
	"robbykansas/another-novel-scraper/cmd/ui/spinner"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func GetContent(content string, folder string, title string) {

	// Track memory usage before the benchmark
	beforeMemStats := testbenchmark.GetMemStats()

	// Start timer for benchmarking with pooling
	start := time.Now()

	list := models.ListContent(content, title)
	wp := models.WorkerPoolContent{
		List:        list.Data,
		Concurrency: 10,
	}
	var SortContent []models.ListChapter

	//  Spinner UI
	spinnerModel := spinner.InitialModel()
	s := tea.NewProgram(spinnerModel)

	go func() {
		if _, err := s.Run(); err != nil {
			cobra.CheckErr(err)
		}
	}()

	s.Send(tea.QuitMsg{})

	// Progressbar UI
	progressbarModel := progressbar.InitialModel(len(wp.List))
	p := tea.NewProgram(progressbarModel)

	go func() {
		for {
			content, ok := <-wp.Res
			if ok {
				copy := *content
				SortContent = append(SortContent, copy)
				content.Reset()
				wp.Pool.Put(content)
				p.Send(progressbar.ProgressMsg{})
			} else {
				time.Sleep(1000 * time.Millisecond)
				p.Send(tea.Quit())
				break
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if _, err := p.Run(); err != nil {
					log.Fatalf("error running progressbar message: %v", err)
				}
			}
		}
	}(ctx)

	wp.Run(content, title)

	cancel()

	sort.Slice(SortContent, func(i, j int) bool {
		return SortContent[i].Order < SortContent[j].Order
	})

	list.Data = SortContent

	elapsed := time.Since(start)
	fmt.Printf("Time taken (with pool): %s\n", elapsed)
	time.Sleep(1 * time.Second)
	fmt.Println(len(SortContent))
	// Track memory usage after the benchmark
	afterMemStats := testbenchmark.GetMemStats()

	// Print memory statistics before and after
	fmt.Printf("Memory Before: Alloc = %v, TotalAlloc = %v, Sys = %v, HeapAlloc = %v, HeapSys = %v\n",
		beforeMemStats.Alloc, beforeMemStats.TotalAlloc, beforeMemStats.Sys, beforeMemStats.HeapAlloc, beforeMemStats.HeapSys)
	fmt.Printf("Memory After: Alloc = %v, TotalAlloc = %v, Sys = %v, HeapAlloc = %v, HeapSys = %v\n",
		afterMemStats.Alloc, afterMemStats.TotalAlloc, afterMemStats.Sys, afterMemStats.HeapAlloc, afterMemStats.HeapSys)

	// Memory difference
	memDiff := afterMemStats.Alloc - beforeMemStats.Alloc
	fmt.Printf("Memory difference: %v bytes\n", memDiff)

	// Track garbage collection (GC) stats
	fmt.Printf("GC Cycles: %d\n", afterMemStats.NumGC-beforeMemStats.NumGC)
	fmt.Printf("GC Pause (ns): %d\n", afterMemStats.PauseTotalNs-beforeMemStats.PauseTotalNs)

	epub.SetEpub(folder, list)
}
