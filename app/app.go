package app

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/common-nighthawk/go-figure"
)


func Fuzzer() *cli.App {

	//Application Banner
	banner := figure.NewColorFigure("Fuzzer", "larry3d", "cyan", true)
	banner2 := figure.NewColorFigure("web...", "larry3d", "cyan", true)
	banner.Print()
	banner2.Print()
	fmt.Println()

	
	
	application := &cli.App{

		Name: "Fuzzing webs",
		Usage: "Application to enumerate web directories",

		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:"threads",
				Value: 1,
				Usage: "Select the number of threads you want to use",
			},
			&cli.StringFlag{
				Name: "host",
				Usage: "Url website ex: https://yoursite.com/FUZZER",
			},
			&cli.StringFlag{
				Name: "file",
				Usage: "Path to wordlist",
			},
			&cli.IntSliceFlag{
				Name: "hdc",
				Usage: "Hidden unwanted status code  ex: 201,404",
			},
			&cli.Int64SliceFlag{
				Name: "hcl",
				Usage: "Hidden unwanted content lenght  ex: 5213,63234",
			},
		},

		Action: func(ctx *cli.Context) error {
			
			threads 			:= ctx.Int("threads")
			host 				:= ctx.String("host")
			filePath 			:= ctx.String("file")
			hiddenStatusCode 	:= ctx.IntSlice("hdc")
			hiddenContentLenght := ctx.Int64Slice("hcl")

			

			// Start the timer
			startTime := time.Now()

			//Verify hostname
			if !ctx.IsSet("host"){
				return errors.New("missing host parameter, use --help")
			}else if !strings.Contains(host, "FUZZER"){
				return errors.New("missing FUZZER word at hostname parameter, use --help")
			}

			fmt.Printf("HOST: %s\n\n", strings.Replace(host, "FUZZER", "", -1))

			paths, err := fileHandler(filePath)
			if err != nil {
				return err
			}

			//instantiates the channels
			jobs := make(chan string, len(paths))
			results := make(chan string, len(paths))

			//Create a WaitGroup
			var wg sync.WaitGroup

			//Start the workers
			for i := 1; i <= threads; i++ {
				wg.Add(1)
				go worker(i, jobs, results, &wg, host, hiddenStatusCode, hiddenContentLenght)
			}

			//Send jobs to the workes
			for _, path := range paths{
				jobs <- path
			}
			close(jobs)

			// Goroutine to print the results
			go func() {
				for result := range results {
					fmt.Println(result)
				}
			}()


			// Wait for all workers to finish and close the chan
			wg.Wait()
			close(results)

			// Calculate the execution time 
			duration := time.Since(startTime)
			fmt.Printf("Execution time: %v\n", duration)

			return nil
		},
	
	}

	return application


}


/*
	FileHandler processes the file as follows:
		- Checks if the file can be opened
		- Starts reading line by line, adding each item to a slice
		- Returns:
			- An error if any occurred during the process
			- The slice with the read values
*/
func fileHandler(filePath string) ([]string, error) {
	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var paths []string
	linesFile := bufio.NewScanner(file)

	for linesFile.Scan(){
		paths = append(paths, linesFile.Text())
	}

	if err := linesFile.Err(); err != nil{
		return nil, fmt.Errorf("error reading file: %v", err)
	}else{
		return paths, nil
	}

}

// isStatusCodeHidden responsible to remove the unwanted status code
func isStatusCodeHidden(statusCode int, hiddenStatusCodes []int) bool {
	for _, code := range hiddenStatusCodes {
		if code == statusCode {
			return true
		}
	}
	return false
}

// isContentLengthHidden responsible to remove the unwanted content lenght
func isContentLengthHidden(contentLength int64, hiddenContentLengths []int64) bool {
	for _, length := range hiddenContentLengths {
		if length == contentLength {
			return true
		}
	}
	return false
}

/*
	worker starts the execution of the fuzzer system.

	- Receives a channel `jobs` containing paths that will replace the "FUZZ" placeholder in the base URL `host`.
	
	- Path:
		- Replaces "FUZZ" in the base URL with the current path name.
		- Executes the HTTP request for the generated URL.
		- Extracts the Status Code and Content Length from the response.
		- Verifies if the Status Code or Content Length should be hidden based on user-defined parameters.
		- Sends the result of the request back through the `results` channel if it meets the criteria.

	- Uses colored output to highlight thread ID, URL, status code, and content length.
	- Signals the completion of the worker's task to the `WaitGroup`.
*/
func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup, host string, hiddenStatusCodes []int, hiddenContentLengths []int64){

	defer wg.Done() //wait to the function finish to execute sinalize that is already done

	for path := range jobs{

		url := strings.Replace(host, "FUZZER", path, 1)

		//Execute the request and get the response
		response, err := http.Get(url)

		if err != nil{
			// results <- fmt.Sprintf("Thread %d: Error fetching %s: %v", id, url, err)
			continue
		}

		statusCode := response.StatusCode
		contentLength := response.ContentLength

		if !isStatusCodeHidden(statusCode, hiddenStatusCodes) && !isContentLengthHidden(contentLength, hiddenContentLengths) {
			threadColor := color.New(color.FgCyan).SprintFunc()
			urlColor := color.New(color.FgBlue).SprintFunc()
			contentLengthColor := color.New(color.FgYellow).SprintFunc()
			var statusColor func(validade ...interface{}) string
			
			if statusCode != 200{
				statusColor = color.New(color.FgRed).SprintFunc()
			}else{
				statusColor = color.New(color.FgGreen).SprintFunc()
			}

			// Format results with colors output
			
			result := fmt.Sprintf(
				"Thread %s: %s - Status: %s, Size: %s bytes",
				threadColor(fmt.Sprintf("%d", id)),
				urlColor(path),
				statusColor(fmt.Sprintf("%d", statusCode)),
				contentLengthColor(fmt.Sprintf("%d", contentLength)),
			)

			results <- result
		}

		response.Body.Close()

	}

	
}

