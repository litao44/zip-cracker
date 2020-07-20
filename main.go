package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"

	"github.com/yeka/zip"

	"github.com/litao44/zip-cracker/generator"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	in := flag.String("in", "", "--in ./test.zip")
	co := flag.Int("co", runtime.NumCPU(), "--co 4")
	max := flag.Int("max", 6, "--max 6")
	min := flag.Int("min", 4, "--min 4")
	pool := flag.String("pool", generator.DefaultPool, `--pool "123456"`)

	flag.Parse()

	inFile := *in
	if inFile == "" {
		printUsage()
		os.Exit(0)
	}

	err := run(*in, *co, *max, *min, *pool)
	if err != nil {
		log.Fatal(err)
	}
}

func run(zipFile string, concurrent int, maxPasswordLen, minPasswordLen int, pool string) error {
	log.Printf("crack %s, %d goroutines\n", zipFile, concurrent)

	gen, err := generator.NewPasswordGeneratorWithPool(maxPasswordLen, minPasswordLen, pool)
	if err != nil {
		return err
	}

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	var file *zip.File
	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		file = f
		break
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	onSignal(cancel, os.Interrupt)

	result := crack(ctx, concurrent, file, gen)
	if result.ok {
		log.Printf("crack success, password is %s\n", result.password)
	} else {
		log.Printf("crack done, no password found\n")
	}

	return nil
}

func onSignal(callback func(), signals ...os.Signal) {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, signals...)
		select {
		case s := <-sig:
			log.Printf("catch signal: %s\n", s.String())
			callback()
		}
	}()
}

type crackResult struct {
	password string
	ok       bool
}

func crack(ctx context.Context, concurrent int, file *zip.File, gen generator.PasswordGeneratorInterface) *crackResult {
	done := make(chan *crackResult, concurrent)

	wg := sync.WaitGroup{}
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			crackWorker(ctx, done, file, gen)
		}()
	}

	result := <-done
	wg.Wait()
	close(done)
	return result
}

func crackWorker(ctx context.Context, done chan<- *crackResult, file *zip.File, gen generator.PasswordGeneratorInterface) {
	failResult := &crackResult{
		ok: false,
	}

	for {
		select {
		case <-ctx.Done():
			done <- failResult
			return

		default:
			password, last := gen.Generate()
			if correctPassword(file, password) {
				result := &crackResult{
					password: password,
					ok:       true,
				}
				done <- result
				return
			}

			if last {
				done <- failResult
				return
			}
		}
	}
}

func correctPassword(file *zip.File, password string) bool {
	log.Printf("check password: %s\n", password)

	file.SetPassword(password)

	r, err := file.Open()
	if err != nil {
		return false
	}
	defer r.Close()

	_, err = io.Copy(ioutil.Discard, r)
	if err != nil {
		return false
	}

	log.Printf("password %s ok!\n", password)

	return true
}

func printUsage() {
	fmt.Printf("%s --in test.zip\n", os.Args[0])
}
