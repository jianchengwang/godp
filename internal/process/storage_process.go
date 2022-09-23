package process

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

func OssCreateBudget(endpoint string, accessKeyId string, accessKeySecret string, bucketName string) error {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	err = client.CreateBucket(bucketName)
	if err != nil {
		return err
	}

	return nil
}

func OssDoGetObject(endpoint string, accessKeyId string, accessKeySecret string, bucketName string, objectKey string, localFilePath string) error {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = bucket.GetObjectToFile(objectKey, localFilePath)

	return err
}

func OssDoPutObject(endpoint string, accessKeyId string, accessKeySecret string, bucketName string, objectKey string, localFilePath string) error {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	// 设置跨域资源共享规则
	rule1 := oss.CORSRule{
		AllowedOrigin: []string{"*"},
		AllowedMethod: []string{"GET"},
		AllowedHeader: []string{},
		ExposeHeader:  []string{},
		MaxAgeSeconds: 200,
	}
	err = client.SetBucketCORS(bucketName, []oss.CORSRule{rule1})
	if err != nil {
		fmt.Println("Error:", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(objectKey, localFilePath)
	if err != nil {
		return err
	}

	return nil
}

func OssSetCORSRule(endpoint string, accessKeyId string, accessKeySecret string, bucketName string) error {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	// 设置跨域资源共享规则
	rule1 := oss.CORSRule{
		AllowedOrigin: []string{"*"},
		AllowedMethod: []string{"GET"},
		AllowedHeader: []string{},
		ExposeHeader:  []string{},
		MaxAgeSeconds: 200,
	}
	err = client.SetBucketCORS(bucketName, []oss.CORSRule{rule1})
	if err != nil {
		return err
	}
	return nil
}

func CosCreateBudget(regionId string, secretId string, secretKey string, bucketName string) error {
	// 将 examplebucket-1250000000 和 COS_REGION 修改为真实的信息
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, _ := url.Parse("https://" + bucketName + ".cos." + regionId + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: secretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})

	_, err := c.Bucket.Put(context.Background(), nil)
	if err != nil {
		return err
	}

	return nil
}

func CosDoGetObject(regionId string, secretId string, secretKey string, bucketName string, objectKey string, localFilePath string) error {
	// 将 examplebucket-1250000000 和 COS_REGION 修改为真实的信息
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, _ := url.Parse("https://" + bucketName + ".cos." + regionId + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: secretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	name := objectKey
	// 1.通过字符串上传对象
	//f := strings.NewReader("test")
	//
	//
	//_, err := c.Object.Put(context.Background(), name, f, nil)
	//if err != nil {
	//	panic(err)
	//}
	// 2.通过本地文件上传对象
	_, err := c.Object.PutFromFile(context.Background(), name, localFilePath, nil)
	if err != nil {
		return err
	}
	// 3.通过文件流上传对象
	//fd, err := os.Open("./test")
	//if err != nil {
	//	panic(err)
	//}
	//defer fd.Close()
	//_, err = c.Object.Put(context.Background(), name, fd, nil)
	//if err != nil {
	//	panic(err)
	//}

	return nil
}

func CosDoPutObject(regionId string, secretId string, secretKey string, bucketName string, objectKey string, localFilePath string) error {
	// 将 examplebucket-1250000000 和 COS_REGION 修改为真实的信息
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, _ := url.Parse("https://" + bucketName + ".cos." + regionId + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: secretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	name := objectKey
	// Case1 通过resp.Body下载对象，Body需要关闭
	//name := "test/example"
	//resp, err := c.Object.Get(context.Background(), name, nil)
	//log_status(err)
	//
	//bs, _ := ioutil.ReadAll(resp.Body)
	//resp.Body.Close()
	//fmt.Printf("%s\n", string(bs))

	// Case2 下载对象到文件. Body需要关闭
	//fd, err := os.OpenFile("test", os.O_WRONLY|os.O_CREATE, 0660)
	//log_status(err)
	//
	//defer fd.Close()
	//resp, err = c.Object.Get(context.Background(), name, nil)
	//log_status(err)
	//
	//io.Copy(fd, resp.Body)
	//resp.Body.Close()

	// Case3 下载对象到文件
	_, err := c.Object.GetToFile(context.Background(), name, localFilePath, nil)
	if err != nil {
		return err
	}

	// Case4 range下载对象，可以根据range实现并发下载
	//opt := &cos.ObjectGetOptions{
	//	ResponseContentType: "text/html",
	//	Range:               "bytes=0-3",
	//}
	//resp, err = c.Object.Get(context.Background(), name, opt)
	//log_status(err)
	//bs, _ = ioutil.ReadAll(resp.Body)
	//resp.Body.Close()
	//fmt.Printf("%s\n", string(bs))

	// Case5 下载对象到文件，查看下载进度
	//opt = &cos.ObjectGetOptions{
	//	Listener: &cos.DefaultProgressListener{},
	//}
	//_, err = c.Object.GetToFile(context.Background(), name, "test", opt)
	//log_status(err)

	return nil
}
