package process

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"godp/internal/global"
)

// ResponseContent返回体。
type ResponseContent struct {
	Data      map[string]string `json:"data"`
	RequestID string            `json:"requestId"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
}

func CreateRepo(repoName string) error {

	// registry.cn-hangzhou.aliyuncs.com/sc_nft
	// 需要替换AccessKey ID和AccessKey Secret。
	client, err := cr.NewClientWithAccessKey(
		global.RepoRegionId, global.RepoAccessKeyId, global.RepoAccessKeySecret)
	request := cr.CreateCreateRepoRequest()

	// 下述参数无需调整。
	request.Domain = "cr." + global.RepoRegionId + ".aliyuncs.com"
	request.SetContentType("JSON")

	content := fmt.Sprintf(
		`{
			"Repo":{
				"RepoNamespace": "%s",
				"RepoName":      "%s", 
				"Summary":       "%s",
				"Detail":        "%s", 
				"RepoType":      "PRIVATE", 
			}
		}`, global.RepoNamespace, repoName, repoName, repoName,
	)

	request.SetContent([]byte(content))

	response, err := client.CreateRepo(request)
	print(response)
	return err
}
