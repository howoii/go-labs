package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	runner = flag.Int("runner", 1, "concurrent threads num")
	delay  = flag.Int("delay", 100, "delay between calls (ms)")
	stdout = flag.Bool("stdout", false, "output log to console")

	logPath = "./press.log"

	wg     sync.WaitGroup
	errc   chan error
	sigc   chan os.Signal
	client *http.Client
	logger *zap.Logger
)

func main() {
	flag.Parse()
	initClient()
	initLogger()
	initSignal()
	errc = make(chan error)
	//run threads
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < *runner; i++ {
		wg.Add(1)
		go thread(ctx, i)
	}

	select {
	case <-ctx.Done():
	case <-sigc:
		cancel()
	case <-errc:
		cancel()
	}
	wg.Wait()
}

func initClient() {
	client = &http.Client{
		Timeout: time.Second * 10,
	}
}

func initLogger() {
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	var w io.Writer
	if *stdout {
		w = os.Stdout
	} else {
		w, _ = os.Create(logPath)
	}
	syncer := zapcore.AddSync(w)
	core := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)
	logger = zap.New(core)
}

func initSignal() {
	sigc = make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGQUIT)
}

func thread(ctx context.Context, threadT int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(*delay)))
	timer := time.NewTicker(time.Millisecond * time.Duration(*delay))
	defer func() {
		wg.Done()
		timer.Stop()
	}()
	var iteration int64
	tLogger := logger.With(zap.Int("thread", threadT))
	for {
		iteration++
		t := time.Now()

		res, err := client.Get("http://localhost:5000/search")
		if err != nil {
			tLogger.Error("Get Failed", zap.Error(err))
		} else {
			tLogger.Info("Get Success", zap.Int64("iteration", iteration), zap.String("status", res.Status), zap.Int64("duration", time.Since(t).Milliseconds()))
			res.Body.Close()
		}

		select {
		case <-ctx.Done():
			fmt.Printf("thread %d exit\n", threadT)
			return
		case <-timer.C:
		}
	}
}
