package loadtest

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func worker(jobs <-chan int, results chan<- int, client *http.Client, url string, wg *sync.WaitGroup) {
	defer wg.Done()
	for range jobs {
		resp, err := client.Get(url)
		if err != nil {
			results <- 0
			continue
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		results <- resp.StatusCode
	}
}

// Run executes the CLI with given args and writes report to out. Returns exit code.
func Run(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("loadtest", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	url := fs.String("url", "", "URL do serviço a ser testado")
	requests := fs.Int("requests", 1, "Número total de requisições a serem realizadas")
	concurrency := fs.Int("concurrency", 1, "Número de chamadas simultâneas")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(out, "erro ao parsear flags:", err)
		return 1
	}

	if *url == "" {
		fmt.Fprintln(out, "--url é obrigatório")
		return 1
	}
	if *requests <= 0 || *concurrency <= 0 {
		fmt.Fprintln(out, "--requests e --concurrency devem ser maiores que 0")
		return 1
	}
	if *concurrency > *requests {
		*concurrency = *requests
	}

	client := &http.Client{Timeout: 15 * time.Second}
	jobs := make(chan int, *requests)
	results := make(chan int, *requests)

	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(jobs, results, client, *url, &wg)
	}

	for i := 0; i < *requests; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(results)

	elapsed := time.Since(start)

	total := 0
	counts := make(map[int]int)
	for status := range results {
		counts[status]++
		total++
	}

	fmt.Fprintln(out, "===== Relatório do Teste =====")
	fmt.Fprintf(out, "Tempo total: %s\n", elapsed)
	fmt.Fprintf(out, "Total de requests: %d\n", total)
	fmt.Fprintf(out, "200 OK: %d\n", counts[200])
	fmt.Fprintln(out, "Distribuição de códigos HTTP:")
	for code, c := range counts {
		fmt.Fprintf(out, "%d: %d\n", code, c)
	}

	return 0
}

func main() {
	code := Run(os.Args[1:], os.Stdout)
	os.Exit(code)
}
