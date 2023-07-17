package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func DemoSignal() {
	// Tạo một kênh để nhận tín hiệu từ hệ thống
	sigChan := make(chan os.Signal, 1)

	// Đăng ký tín hiệu SIGINT và SIGTERM với kênh
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Chạy một goroutine để lắng nghe và xử lý tín hiệu
	go func() {
		// Chờ nhận tín hiệu từ kênh
		sig := <-sigChan
		fmt.Println("\nReceived signal:", sig)

		// Thực hiện các tác vụ hoặc giải phóng tài nguyên trước khi thoát
		// ...

		// Thoát chương trình
		os.Exit(0)
	}()

	// Tiếp tục thực hiện các công việc khác
	// ...

	// Chờ cho đến khi nhận tín hiệu để thoát
	<-sigChan
}

func Demo2() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {

		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
