package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
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
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			server1.Shutdown(context.Background())
			server2.Shutdown(context.Background())
			shutdown <- struct{}{}
			return
		}

	}()

	if err := g.Wait(); err != nil {
		log.Println(err) //打印第一个错误
	}
	<-shutdown
	log.Println("all servers shutdown success")
}
