package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

const (
	keyPrefix = "bench:"
	keyLength = 16
	zsetSize  = 500
	listSize  = 100
)

var (
	mode     = flag.String("mode", "ring", "cluster | ring")
	nRequest = flag.Int("requests", 10000, "num of requests per client")
	nClient  = flag.Int("clients", 20, "num of clients")
	nKey     = flag.Int("keys", 100000, "random key nums")
	batchKey = flag.Int("batchKeys", 10, "key num in pipeline operation")
)

type RedisClient interface {
	redis.Cmdable
}

func init() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func newClusterClient(config *viper.Viper) RedisClient {
	rc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: config.GetStringSlice("addrs"),
	})
	if err := rc.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
	return rc
}

func newRingClient(config *viper.Viper) RedisClient {
	rc := redis.NewRing(&redis.RingOptions{
		Addrs: config.GetStringMapString("addrs"),
	})
	err := rc.ForEachShard(context.Background(), func(ctx context.Context, client *redis.Client) error {
		if err := client.Ping(ctx).Err(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return rc
}

func newRedisClient(mode string, config *viper.Viper) RedisClient {
	switch mode {
	case "ring":
		return newRingClient(config)
	case "cluster":
		return newClusterClient(config)
	}
	return nil
}

func randomString(length int) string {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(b)
}

type benchConfig struct {
	mode    string
	nReq    int
	nClient int
	nKey    int
	nBtach  int

	redisConfig *viper.Viper
}

type host struct {
	config *benchConfig

	ctx      context.Context
	cancel   context.CancelFunc
	running  sync.WaitGroup
	signalC  chan os.Signal
	counterC chan int64

	keys    []string
	clients []RedisClient
}

func newHost(config *benchConfig) *host {
	ctx, cancel := context.WithCancel(context.Background())
	h := &host{
		ctx:      ctx,
		cancel:   cancel,
		config:   config,
		counterC: make(chan int64, 1000),
		signalC:  make(chan os.Signal, 1),
		keys:     make([]string, 0, config.nKey),
		clients:  make([]RedisClient, 0, config.nClient),
	}
	for i := 0; i < config.nKey; i++ {
		h.keys = append(h.keys, keyPrefix+randomString(keyLength))
	}
	for i := 0; i < config.nClient; i++ {
		c := newRedisClient(config.mode, config.redisConfig)
		h.clients = append(h.clients, c)
	}
	h.watch()
	return h
}

func (h *host) watch() {
	signal.Notify(h.signalC,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGQUIT,
	)
	go func() {
		select {
		case <-h.signalC:
			h.cancel()
			return
		}
	}()
}

func (h *host) RandomKey(poolSize int) string {
	if poolSize <= 0 {
		poolSize = len(h.keys)
	} else if poolSize > len(h.keys) {
		poolSize = len(h.keys)
	}
	return h.keys[rand.Intn(poolSize)]
}

func (h *host) BenchmarkCmd(cmd *redisCmd, delete bool) bool {
	done := make(chan struct{})

	var totalReq int64
	startTime := time.Now()

	for _, c := range h.clients {
		h.running.Add(1)
		go h.startClient(c, cmd)
	}
	go func() {
		h.running.Wait()
		close(done)
	}()

	showAndClear := func(exit bool) {
		elapse := time.Since(startTime).Seconds()
		fmt.Printf("%-10s: %d total reqs, %f ops/sec\n", cmd.Name(), totalReq, float64(totalReq)/elapse)
		if delete || exit {
			if exit {
				fmt.Println("cleaning ...")
			}
			h.deleteBenchKeys(h.clients[0])
		}
	}
	for {
		select {
		case n := <-h.counterC:
			totalReq += n
		case <-h.ctx.Done():
			showAndClear(true)
			h.running.Wait()
			return false
		case <-done:
			showAndClear(false)
			return true
		}
	}
}

func (h *host) startClient(c RedisClient, cmd *redisCmd) {
	defer func() {
		h.running.Done()
	}()
	reportStep := 100
	lastReport := 0
	for i := 1; i <= h.config.nReq; i++ {
		cmd.Do(c, h)

		if i%reportStep == 0 || i == h.config.nReq {
			h.counterC <- int64(i - lastReport)
			lastReport = i
		}

		select {
		case <-h.ctx.Done():
			return
		default:
		}
	}
}

func (h *host) deleteBenchKeys(c RedisClient) {
	ctx := context.Background()
	for _, key := range h.keys {
		c.Del(ctx, key)
	}
}

type cmdFunction func(c RedisClient, h *host) error

func pingFunc(c RedisClient, h *host) error {
	return c.Ping(h.ctx).Err()
}

func setFunc(c RedisClient, h *host) error {
	val := randomString(rand.Intn(64) + 2)
	return c.Set(h.ctx, h.RandomKey(0), val, time.Hour*12).Err()
}

func getFunc(c RedisClient, h *host) error {
	return c.Get(h.ctx, h.RandomKey(0)).Err()
}

func ttlFunc(c RedisClient, h *host) error {
	return c.TTL(h.ctx, h.RandomKey(0)).Err()
}

func delFunc(c RedisClient, h *host) error {
	return c.Del(h.ctx, h.RandomKey(0)).Err()
}

func zaddFunc(c RedisClient, h *host) error {
	poolSize := h.config.nReq * h.config.nClient / zsetSize
	member := redis.Z{
		Score:  rand.Float64() * 10000,
		Member: rand.Intn(zsetSize/2) + 1,
	}
	return c.ZAdd(h.ctx, h.RandomKey(poolSize), &member).Err()
}

func zrankFunc(c RedisClient, h *host) error {
	poolSize := h.config.nReq * h.config.nClient / zsetSize
	member := strconv.Itoa(rand.Intn(zsetSize/2) + 1)
	return c.ZRank(h.ctx, h.RandomKey(poolSize), member).Err()
}

func zrangeFunc(c RedisClient, h *host) error {
	poolSize := h.config.nReq * h.config.nClient / zsetSize
	return c.ZRange(h.ctx, h.RandomKey(poolSize), 0, 49).Err()
}

func lpushFunc(c RedisClient, h *host) error {
	poolSize := h.config.nReq * h.config.nClient / listSize
	return c.LPush(h.ctx, h.RandomKey(poolSize), 0).Err()
}

func lpopFunc(c RedisClient, h *host) error {
	poolSize := h.config.nReq * h.config.nClient / listSize
	return c.LPop(h.ctx, h.RandomKey(poolSize)).Err()
}

func lrangeFunc(c RedisClient, h *host) error {
	poolSize := h.config.nReq * h.config.nClient / listSize
	return c.LRange(h.ctx, h.RandomKey(poolSize), 0, 49).Err()
}

func batchSetFunc(c RedisClient, h *host) error {
	val := randomString(rand.Intn(64) + 2)
	_, err := c.Pipelined(h.ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < h.config.nBtach; i++ {
			pipe.Set(h.ctx, h.RandomKey(0), val, time.Hour*12)
		}
		return nil
	})
	return err
}

func batchGetFunc(c RedisClient, h *host) error {
	_, err := c.Pipelined(h.ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < h.config.nBtach; i++ {
			pipe.Get(h.ctx, h.RandomKey(0))
		}
		return nil
	})
	return err
}

func batchDelFunc(c RedisClient, h *host) error {
	_, err := c.Pipelined(h.ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < h.config.nBtach; i++ {
			pipe.Del(h.ctx, h.RandomKey(0))
		}
		return nil
	})
	return err
}

type redisCmd struct {
	name string
	fn   cmdFunction
}

func (cmd *redisCmd) Name() string {
	return cmd.name
}

func (cmd *redisCmd) Do(c RedisClient, h *host) error {
	return cmd.fn(c, h)
}

func newCmd(name string, fn cmdFunction) *redisCmd {
	return &redisCmd{
		name: name,
		fn:   fn,
	}
}

func main() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath("./")
	if err := config.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var rcConfig *viper.Viper
	switch *mode {
	case "ring":
		rcConfig = config.Sub("redis_ring")
	case "cluster":
		rcConfig = config.Sub("redis_cluster")
	default:
		flag.Usage()
		return
	}

	bConfig := &benchConfig{
		mode:        *mode,
		nReq:        *nRequest,
		nClient:     *nClient,
		nKey:        *nKey,
		nBtach:      *batchKey,
		redisConfig: rcConfig,
	}
	h := newHost(bConfig)

	benchCases := []struct {
		name  string
		fn    cmdFunction
		clean bool
	}{
		{"PING", pingFunc, false},
		{"SET", setFunc, false},
		{"GET", getFunc, false},
		{"TTL", ttlFunc, false},
		{"DEL", delFunc, true},
		{fmt.Sprintf("BATCH_SET(%d keys)", bConfig.nBtach), batchSetFunc, false},
		{fmt.Sprintf("BATCH_GET(%d keys)", bConfig.nBtach), batchGetFunc, false},
		{fmt.Sprintf("BATCH_DEL(%d keys)", bConfig.nBtach), batchDelFunc, true},
		{"ZADD", zaddFunc, false},
		{"ZRANK", zrankFunc, false},
		{"ZRANGE_50", zrangeFunc, true},
		{"LPUSH", lpushFunc, false},
		{"LRANGE_50", lrangeFunc, false},
		{"LPOP", lpopFunc, true},
	}
	for _, v := range benchCases {
		if !h.BenchmarkCmd(newCmd(v.name, v.fn), v.clean) {
			break
		}
	}
}
