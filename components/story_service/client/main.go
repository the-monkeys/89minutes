package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/89minutes/89minutes/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"

	// Adjust the size for which a large file will be broken
	// down into multiple parts during a stream request
	chunkSize = 5 * MB
)

const MB = 1 << 20

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStoryServiceClient(conn)

	if len(os.Args) == 1 {
		log.Fatal("need filename argument!")
	}
	fileName := os.Args[1]

	UploadImage(c, uuid.New().String(), fileName)
}

func UploadImage(laptopClient pb.StoryServiceClient, storyId string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	fileInfo, err := os.Stat(file.Name())
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.GetBlob(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	// req := &pb.UploadImageRequest{
	// 	Data: &pb.UploadImageRequest_Info{
	// 		Info: &pb.ImageInfo{
	// 			LaptopId:  laptopID,
	// 			ImageType: filepath.Ext(imagePath),
	// 		},
	// 	},
	// }
	sizeFile := fileInfo.Size()
	req := &pb.StoryRequest{
		Data: &pb.StoryRequest_Info{
			Info: &pb.FileHeader{
				Name:     file.Name(),
				FileSize: &sizeFile,
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		// req := &pb.UploadImageRequest{
		// 	Data: &pb.UploadImageRequest_ChunkData{
		// 		ChunkData: buffer[:n],
		// 	},
		// }

		req := &pb.StoryRequest{
			Data: &pb.StoryRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}

// func UploadStream(client pb.StoryServiceClient, fileName string) time.Duration {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	fh, err := os.Open(fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer fh.Close()

// 	stat, err := fh.Stat()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	size := stat.Size()
// 	log.Printf("[UploadStream] Will stream file %q with size %d\n", fh.Name(), size)

// 	start := time.Now()

// 	stream, err := client.GetBlob(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	header := &pb.FileHeader{Name: fh.Name(), FileSize: &size}
// 	err = stream.Send(&pb.BlobReq{

// 		Contents: &pb.BlobReq_Header{
// 			Header: header,
// 		},
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Sent header. Now sending data chunks...")
// 	// time.Sleep(5 * time.Second)

// 	buf := make([]byte, chunkSize)
// 	chunkCount := 0
// 	for {
// 		n, err := fh.Read(buf)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		//fmt.Printf("Sending chunk #%d with size %d\n", i, n)
// 		err = stream.Send(&pb.BlobReq{
// 			Id:       uuid.New().String(),
// 			Title:    "My First Blog",
// 			Contents: &pb.BlobReq_Chunk{Chunk: buf[:n]}})
// 		chunkCount++
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	stream.CloseAndRecv()

// 	took := time.Since(start)
// 	log.Printf("  Sent %d chunk(s)\n", chunkCount)
// 	log.Printf("[UploadStream] Took: %s\n", took)
// 	return took
// }
