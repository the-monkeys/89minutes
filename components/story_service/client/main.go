package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/89minutes/89minutes/pb"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pool "gopkg.in/cheggaaa/pb.v1"
)

const (
	address = "localhost:8080"

	// Adjust the size for which a large file will be broken
	// down into multiple parts during a stream request
	chunkSize = 5 * MB
)

const MB = 1 << 20

func main() {
	app := cli.NewApp()
	app.Name = "89 minutes client"
	app.Usage = "Upload file"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		uploadCommand(),
	}
	app.Run(os.Args)
}

type uploader struct {
	id          string
	dir         string
	client      pb.StoryServiceClient
	ctx         context.Context
	wg          sync.WaitGroup
	requests    chan string // each request is a filepath on client accessible to client
	pool        *pool.Pool
	DoneRequest chan string
	FailRequest chan string
}

// /NewUploader creates a object of type uploader and creates fixed worker goroutines/threads
func NewUploader(ctx context.Context, client pb.StoryServiceClient, dir string, id string) *uploader {
	d := &uploader{
		id:          id,
		ctx:         ctx,
		client:      client,
		dir:         dir,
		requests:    make(chan string),
		DoneRequest: make(chan string),
		FailRequest: make(chan string),
	}
	for i := 0; i < 5; i++ {
		d.wg.Add(1)
		go d.worker(i + 1)
	}
	d.pool, _ = pool.StartPool()
	return d
}

func (d *uploader) Stop() {
	close(d.requests)
	d.wg.Wait()
	d.pool.RefreshRate = 500 * time.Millisecond
	d.pool.Stop()
}

func (d *uploader) worker(workerID int) {
	logrus.Infof("The d with worker id %d: %+v", workerID, d)
	defer d.wg.Done()
	var (
		buf        []byte
		firstChunk bool
	)
	for request := range d.requests {

		file, err := os.Open(request)
		if err != nil {
			logrus.Errorf("cannot open the file %s: error: %+v", request, err)
			return
		}

		defer file.Close()

		//start uploader
		streamUploader, err := d.client.UploadStoryAndFiles(d.ctx)
		if err != nil {
			logrus.Errorf("failed to create upload stream for file, error: %v", err)
			return
		}

		defer streamUploader.CloseSend()
		stat, err := file.Stat()
		if err != nil {
			logrus.Errorf("Unable to get file size, error %v", err)
			return
		}

		//start progress bar
		bar := pool.New64(stat.Size()).Postfix(" " + filepath.Base(request))
		bar.Units = pool.U_BYTES
		d.pool.Add(bar)

		//create a buffer of chunkSize to be streamed
		buf = make([]byte, chunkSize)
		firstChunk = true
		for {
			n, errRead := file.Read(buf)
			if errRead != nil {
				if errRead == io.EOF {
					errRead = nil
					break
				}

				errRead = errors.Wrapf(errRead,
					"errored while copying from file to buf")
				return
			}
			if firstChunk {
				err = streamUploader.Send(&pb.UploadStoryAndFilesReq{
					Id:       d.id,
					Content:  buf[:n],
					Filename: request,
				})
				firstChunk = false
			} else {
				err = streamUploader.Send(&pb.UploadStoryAndFilesReq{
					Content: buf[:n],
				})
			}
			if err != nil {
				bar.FinishPrint("failed to send chunk via stream file : " + request)
				break
			}

			bar.Add64(int64(n))
		}
		status, err := streamUploader.CloseAndRecv()

		if err != nil { //retry needed

			fmt.Println("failed to receive upstream status response")
			bar.FinishPrint("Error uploading file : " + request + " Error :" + err.Error())
			bar.Reset(0)
			d.FailRequest <- request
			return
		}

		if status.Code != pb.UploadStatusCode_Ok { //retry needed

			bar.FinishPrint("Error uploading file : " + request + " :" + status.Message)
			bar.Reset(0)
			d.FailRequest <- request
			return
		}
		//fmt.Println("writing done for : " + request + " by " + strconv.Itoa(workerID))
		d.DoneRequest <- request
		bar.Finish()
	}

}

func (d *uploader) Do(filepath string) {
	d.requests <- filepath
}

func UploadFiles(ctx context.Context, client pb.StoryServiceClient, filepathlist []string, dir string) error {
	uuid := uuid.New().String()
	logrus.Infof("The id: %v", uuid)
	d := NewUploader(ctx, client, dir, uuid)

	var errorUploadBulk error

	if dir != "" {

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		defer d.Stop()

		go func() {
			for _, file := range files {

				if !file.IsDir() {

					d.Do(dir + "/" + file.Name())

				}
			}
		}()

		for _, file := range files {
			if !file.IsDir() {
				select {

				case <-d.DoneRequest:

					//fmt.Println("sucessfully sent :" + req)

				case req := <-d.FailRequest:

					fmt.Println("failed to  send " + req)
					errorUploadBulk = errors.Wrapf(errorUploadBulk, " Failed to send %s", req)

				}
			}
		}
		fmt.Println("All done ")
	} else {

		go func() {
			for _, file := range filepathlist {
				d.Do(file)
			}
		}()

		defer d.Stop()

		for i := 0; i < len(filepathlist); i++ {
			select {

			case <-d.DoneRequest:
			//	fmt.Println("sucessfully sent " + req)
			case req := <-d.FailRequest:
				fmt.Println("failed to  send " + req)
				errorUploadBulk = errors.Wrapf(errorUploadBulk, " Failed to send %s", req)
			}
		}

	}

	return errorUploadBulk
}

func uploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "Uploads files from server in parallel",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "a",
				Value: "localhost:port",
				Usage: "server address",
			},
			cli.StringFlag{
				Name:  "d",
				Value: ".",
				Usage: "base directory",
			},
			cli.StringFlag{
				Name:  "tls-path",
				Value: "",
				Usage: "directory to the TLS server.crt file",
			},
		},
		Action: func(c *cli.Context) error {
			options := []grpc.DialOption{}
			if p := c.String("tls-path"); p != "" {
				creds, err := credentials.NewClientTLSFromFile(
					filepath.Join(p, "server.crt"),
					"")
				if err != nil {
					log.Println(err)
					return err
				}
				options = append(options, grpc.WithTransportCredentials(creds))
			} else {
				options = append(options, grpc.WithInsecure())
			}
			addr := c.String("a")

			conn, err := grpc.Dial(addr, options...)
			if err != nil {
				log.Fatalf("cannot connect: %v", err)
			}
			defer conn.Close()

			return UploadFiles(context.Background(), pb.NewStoryServiceClient(conn), []string{}, c.String("d"))
		},
	}
}
