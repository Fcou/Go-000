package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	// 用于通知主goroutine，server 已全部关闭
	shutdown := make(chan struct{})
	// 创建带有cancel的父context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 创建errgroup
	g, _ := errgroup.WithContext(ctx)

	// 并发启动server1
	server1 := http.Server{
		Addr: ":8081",
	}
	g.Go(func() error {
		if err := server1.ListenAndServe(); err != nil {
			cancel()
			return err
		}
		return nil
	})
	// 并发启动server2
	server2 := http.Server{
		Addr: ":8082",
	}
	g.Go(func() error {
		if err := server2.ListenAndServe(); err != nil {
			cancel()
			return err // 简单错误直接返回即可，不用pkg/errors处理
		}
		return nil
	})

	// 并发执行监听signal信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			select {
			case s := <-c:
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					cancel()
				default:
				}
			}
		}
	}()

	// context cancel后，关闭全部并发的 http server，全部关闭完成后通知主goroutine
	go func() {
		<-ctx.Done()
		go func() {
			log.Println(ctx.Err())
			if err := server1.Shutdown(context.Background()); err != nil {
				log.Println("server1 shutdown failed, err: %v\n", err)
			}
			if err := server2.Shutdown(context.Background()); err != nil {
				log.Println("server2 shutdown failed, err: %v\n", err)
			}
			log.Println("all servers graceful shutdown")
			close(shutdown)
			return
		}()
		// 超过3分钟没关闭完成，则强制退出
		<-time.After(time.Minute * 3)
		log.Println("all servers no-graceful shutdown")
		close(shutdown)
		return
	}()

	if err := g.Wait(); err != nil {
		// ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
		// var ErrServerClosed = errors.New("http: Server closed")
		log.Println(err) //打印第一个错误
	}
	<-shutdown
}
