package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func Print(path string) error {
	cmd := exec.Command("powershell", "start", "-verb", "printto", path)
	_, err := cmd.Output()

	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	_ = out.Close()
	_ = resp.Body.Close()

	return err
}

func recreateDir() {
	config := GetConfig()
	err := os.Mkdir(config.Main.UploadsDir, 0777)
	if err != nil && !os.IsExist(err) {
		panic(err)
	} else if err != nil {
		err = os.RemoveAll(config.Main.UploadsDir)
		if err != nil {
			panic(err)
		}
		recreateDir()
	}
}

func main() {
	config := GetConfig()

	recreateDir()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.IP, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.Database,
	})

	ctx := context.Background()
	pubsub := rdb.Subscribe(ctx, config.Redis.Channel)

	fmt.Println("Connected, listening for print requests")

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := pubsub.Channel()
	for msg := range ch {
		log.Printf("Printing %s\n", msg.Payload)
		path := fmt.Sprintf("%s/%x.pdf", config.Main.UploadsDir, md5.Sum([]byte(msg.Payload)))

		err := DownloadFile(path, msg.Payload)
		if err != nil {
			log.Println(err)
			continue
		}

		err = Print(path)
		log.Println(err)
	}
}
