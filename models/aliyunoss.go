package models

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func UploadAliyunOss(videoname string, filename string) {
	Endpoint := "https://oss-cn-shenzhen.aliyuncs.com"
	AccessKeyID := "LTAI5tR71zCFyCwUtGdrV8ax"
	AccessKeyIDSecret := "EkxUpOu6Oa0dPVOROUJ69cVuGf79Uo"
	client, err := oss.New(Endpoint, AccessKeyID, AccessKeyIDSecret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	bucket, err := client.Bucket("ling-hu")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	path := "douyin/" + videoname
	err = bucket.PutObjectFromFile(path, filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
