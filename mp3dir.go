package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type TransformAction string

const (
	ConvertAction = TransformAction("CONVERT")
	CopyAction    = TransformAction("COPY")
)

type TransformJob struct {
	action TransformAction
	source string
	dest   string
}

func (t *TransformJob) print() {
	fmt.Printf("Job (%v):\n <- %s\n -> %s\n", t.action, t.source, t.dest)
}

func rebasePathWithSuffix(path string, srclib string, dstlib string, suffix string) (string, error) {
	var result string

	// get track path relative to source folder
	result, err := filepath.Rel(srclib, path)
	if err != nil {
		return "", err
	}

	// change suffix
	result = strings.TrimSuffix(result, filepath.Ext(result)) + suffix

	// rebase relative path into destination folder
	result = filepath.Join(dstlib, result)

	return result, nil
}

func transformLibrary(srclib string, dstlib string) []TransformJob {
	jobs := make([]TransformJob, 0)

	filepath.WalkDir(srclib, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".flac" {
			/* This is a valid track, let's make
			   a job out of it. */
			rebasedPath, err := rebasePathWithSuffix(path, srclib, dstlib, ".mp3")
			if err != nil {
				panic(err)
			}

			jobs = append(jobs, TransformJob{
				action: ConvertAction,
				source: path,
				dest:   rebasedPath,
			})
		}

		// TODO: handle errors
		return nil
	})

	return jobs
}

func runFFmpeg(source string, dest string) error {
	cmd := exec.Command("ffmpeg",
		"-y",

		/* input file */
		"-i", source,

		/* mp3 v0 */
		"-q:a", "0",

		/* output file */
		dest,
	)

	err := cmd.Run()
	return err
}

func runWorker(inbox chan TransformJob, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[%d] spawned job worker\n", workerID)

	for job := range inbox {
		fmt.Printf("[%d] got job: %s\n", workerID, job.source)

		/* TODO: respect job.action */
		err := runFFmpeg(job.source, job.dest)
		if err != nil {
			fmt.Printf("[%d] ffmpeg error: %v\n", workerID, err)
		}
	}
}

func main() {
	var srclib string
	var dstlib string
	var workers int

	flag.StringVar(&srclib, "i", "",
		"Source folder (aka. your lossless music library)")
	flag.StringVar(&dstlib, "o", "",
		"Destination folder (aka. the lossy copy)")
	flag.IntVar(&workers, "j", 4,
		"Number of workers for parallel processing")
	flag.Parse()

	if srclib == "" {
		fmt.Println("Missing: -i")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	if dstlib == "" {
		fmt.Println("Missing: -o")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	jobs := transformLibrary(srclib, dstlib)

	jobChannel := make(chan TransformJob, len(jobs))

	for _, job := range jobs {
		/* Make sure the target directory exists */
		os.MkdirAll(filepath.Dir(job.dest), 0755)

		jobChannel <- job
	}

	close(jobChannel)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go runWorker(jobChannel, i, &wg)
	}

	wg.Wait()
}