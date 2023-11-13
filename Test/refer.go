package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tidwall/gjson"
)

const forInterval = 20
const compareInterval = 70

var fileArg = flag.String("f", "mock.json", "模拟数据文件名")
var onceCyclePublished []int64 = make([]int64, 0)
var repeatCyclePublished []int64 = make([]int64, 0)
var resetRepeatCyclePublished bool = false

type logger struct {
	prefix string
}

func (l logger) Println(v ...interface{}) {
	fmt.Println(append([]interface{}{l.prefix + ":"}, v...)...)
}

func (l logger) Printf(format string, v ...interface{}) {
	if len(format) > 0 && format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	fmt.Printf(l.prefix+":"+format, v...)
}

func ContainsGeneric[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func replaceTimestamp(value string) string {
	timestamp := fmt.Sprintf("%v", time.Now().UnixMilli())
	result := strings.ReplaceAll(value, "[timestamp]", timestamp)
	result = strings.ReplaceAll(result, `"[timestampINT]"`, timestamp)
	return result
}

func getRepeatMessageCycle(repeatArray []gjson.Result) (int64, int64) {
	var maxValue int64 = 0
	var minValue int64 = math.MaxInt64
	for _, v := range repeatArray {
		wait := v.Get("wait").Int()
		if wait > maxValue {
			maxValue = wait
		}
		if wait < minValue {
			minValue = wait
		}
	}
	return maxValue, minValue
}

func uploadFiles(minioClient *minio.Client, uploads []gjson.Result) {
	for _, v := range uploads {
		srcPath := v.Get("srcPath").String()
		destPath := replaceTimestamp(v.Get("destPath").String())
		contentType := v.Get("contentType").String()
		bucket := v.Get("bucket").String()
		upload(minioClient, bucket, srcPath, destPath, contentType)
	}
}

func getMessage(minioClient *minio.Client, startTimestamp int64, onceArray, repeatArray []gjson.Result, repeatCycleMax, repeatCycleMin int64) []*paho.Publish {
	result := make([]*paho.Publish, 0)
	nowTimestamp := time.Now().UnixMilli()
	interval := nowTimestamp - startTimestamp
	for _, v := range onceArray {
		wait := v.Get("wait").Int()
		if wait < interval && interval < wait+compareInterval {
			if !ContainsGeneric(onceCyclePublished, wait) {
				fmt.Printf("interval is %v, wait is %v\n", interval, wait)
				uploads := v.Get("upload").Array()
				if len(uploads) > 0 {
					go uploadFiles(minioClient, uploads)
				}
				topic := v.Get("topic").String()
				if len(topic) > 0 {
					publish := &paho.Publish{
						QoS:     byte(v.Get("qos").Int()),
						Retain:  v.Get("retain").Bool(),
						Topic:   v.Get("topic").String(),
						Payload: []byte(replaceTimestamp(v.Get("payload").String())),
					}
					result = append(result, publish)
					onceCyclePublished = append(onceCyclePublished, wait)
				}
			}
		}
	}
	intervalx := interval % (repeatCycleMax + compareInterval)
	for _, v := range repeatArray {
		wait := v.Get("wait").Int()

		if wait < intervalx && intervalx < wait+compareInterval {
			if wait == repeatCycleMin && resetRepeatCyclePublished {
				repeatCyclePublished = repeatCyclePublished[:0]
				resetRepeatCyclePublished = false
			}
			if !ContainsGeneric(repeatCyclePublished, wait) {
				fmt.Printf("intervalx is %v, wait is %v\n", intervalx, wait)
				uploads := v.Get("upload").Array()
				if len(uploads) > 0 {
					go uploadFiles(minioClient, uploads)
				}
				publish := &paho.Publish{
					QoS:     byte(v.Get("qos").Int()),
					Retain:  v.Get("retain").Bool(),
					Topic:   v.Get("topic").String(),
					Payload: []byte(replaceTimestamp(v.Get("payload").String())),
				}
				result = append(result, publish)
				repeatCyclePublished = append(repeatCyclePublished, wait)
				if wait == repeatCycleMax {
					resetRepeatCyclePublished = true
				}
			}
		}
	}
	return result
}

func initMinIO(endpoint, accessKeyID, secretAccessKey, bucket string) *minio.Client {
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	minioClient, _ := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
		Region: "cn-north-1",
	})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
	// if errBucketExists == nil && exists {
	// 	fmt.Printf("we already own %s\n", bucket)
	// } else {
	// 	fmt.Printf("we already own %s\n", bucket)
	// }
	return minioClient
}

func upload(minioClient *minio.Client, bucketName, srcPath, destPath, contentType string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	info, err := minioClient.FPutObject(ctx, bucketName, destPath, srcPath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Printf("upload %s failure, %v\n", srcPath, err)
	} else {
		fmt.Printf("successfully uploaded %s of size %d\n", srcPath, info.Size)
	}
}

func main() {
	flag.Parse()
	fileName := *fileArg
	jsonFile, _ := os.Open(fileName)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	jsonStr := string(byteValue)
	valueMap := gjson.Get(jsonStr, "config")
	if valueMap.Exists() {
		minioEndpoint := valueMap.Get("minioEndpoint").String()
		minioAccessKeyId := valueMap.Get("minioAccessKeyId").String()
		minioSecretAccessKey := valueMap.Get("minioSecretAccessKey").String()
		minioBucket := valueMap.Get("minioBucket").String()
		minioClient := initMinIO(minioEndpoint, minioAccessKeyId, minioSecretAccessKey, minioBucket)

		host := valueMap.Get("host").String()
		port := valueMap.Get("port").Int()
		username := valueMap.Get("username").String()
		password := valueMap.Get("password").String()
		keeplive := valueMap.Get("keeplive").Int()
		client := replaceTimestamp(valueMap.Get("client").String())
		var serverStr string
		if len(username) > 0 {
			serverStr = fmt.Sprintf("mqtt://%v:%v@%v:%v", username, password, host, port)
		} else {
			serverStr = fmt.Sprintf("mqtt://%v:%v", host, port)
		}
		serverURL, _ := url.Parse(serverStr)
		cliCfg := autopaho.ClientConfig{
			BrokerUrls:     []*url.URL{serverURL},
			KeepAlive:      uint16(keeplive),
			OnConnectionUp: func(*autopaho.ConnectionManager, *paho.Connack) { fmt.Println("mqtt connection up") },
			OnConnectError: func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) },
			Debug:          paho.NOOPLogger{},
			ClientConfig: paho.ClientConfig{
				ClientID:      client,
				OnClientError: func(err error) { fmt.Printf("server requested disconnect: %s\n", err) },
				OnServerDisconnect: func(d *paho.Disconnect) {
					if d.Properties != nil {
						fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
					} else {
						fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
					}
				},
			},
		}
		cliCfg.Debug = logger{prefix: "autoPaho"}
		cliCfg.PahoDebug = logger{prefix: "paho"}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cm, err := autopaho.NewConnection(ctx, cliCfg)
		if err != nil {
			panic(err)
		}

		onceArray := gjson.Get(jsonStr, "data.once").Array()
		repeatArray := gjson.Get(jsonStr, "data.repeat").Array()
		cycleMax, cycleMin := getRepeatMessageCycle(repeatArray)
		fmt.Printf("RepeatMessageCycle Max %v, Min %v\n", cycleMax, cycleMin)
		startTimestamp := time.Now().UnixMilli()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				err = cm.AwaitConnection(ctx)
				if err != nil {
					fmt.Printf("publisher done (AwaitConnection: %s)\n", err)
					return
				}

				publishes := getMessage(minioClient, startTimestamp, onceArray, repeatArray, cycleMax, cycleMin)
				if len(publishes) > 0 {
					for _, v := range publishes {
						pr, err := cm.Publish(ctx, v)
						if err != nil {
							fmt.Printf("error publishing: %s\n", err)
						} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 { // 16 = Server received message but there are no subscribers
							fmt.Printf("reason code %d received\n", pr.ReasonCode)
						} else {
							fmt.Printf("send message: %s\n", v.Payload)
						}
					}
				}

				select {
				case <-time.After(time.Duration(forInterval) * time.Millisecond):
				case <-ctx.Done():
					fmt.Println("publisher done")
					return
				}
			}
		}()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		signal.Notify(sig, syscall.SIGTERM)

		<-sig
		fmt.Println("signal caught - exiting")
		cancel()

		wg.Wait()
		fmt.Println("shutdown complete")
	}
}
