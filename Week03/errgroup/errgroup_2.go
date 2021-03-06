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
		return server1.ListenAndServe()
	})
	// 并发启动server2
	server2 := http.Server{
		Addr: ":8082",
	}
	g.Go(func() error {
		return server2.ListenAndServe()
	})

	// 并发执行监听signal信号，接受到信号则开始关闭全部server流程
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		<-c
		cancel()
	}()

	// context cancel后，关闭全部并发的 http server，全部关闭完成后通知主goroutine
	go func() {
		<-ctx.Done()
		// 开始优雅关闭
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //5秒超时控制
		defer cancel()

		go func() {
			log.Println(ctx.Err())
			if err := server1.Shutdown(ctx); err != nil {
				log.Println("server1 shutdown failed, err: %v\n", err)
			}
			if err := server2.Shutdown(ctx); err != nil {
				log.Println("server2 shutdown failed, err: %v\n", err)
			}
			log.Println("all servers graceful shutdown")
			close(shutdown)
			return
		}()
		<-ctx.Done()
		log.Println("all servers no-graceful shutdown")
		close(shutdown)
		return
	}()

	if err := g.Wait(); err != nil {
		cancel() // 收到第一个错误后，开始关闭全部server流程
		log.Println(err)
	}
	<-shutdown
}
