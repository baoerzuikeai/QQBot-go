package pixiv

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/NateScarlet/pixiv/pkg/artwork"
	"github.com/NateScarlet/pixiv/pkg/client"
)

func InitClient() *client.Client {
	// 使用 PHPSESSID Cookie 登录 (推荐)。
	c := &client.Client{}
	c.SetDefaultHeader("User-Agent", client.DefaultUserAgent)
	c.SetPHPSESSID("110414188_lYQ1F3cdqrydiSeZ031mi2TeDQznoJM3")
	c.CloseIdleConnections()
	return c
}

func Getimage(c *client.Client) string {
	// 所有查询从 context 获取客户端设置, 如未设置将使用默认客户端。
	var ctx = context.Background()
	ctx = client.With(ctx, c)
	imagebase64, _ := GetRankimage(ctx)
	return imagebase64
}

func GetSerchimage(c *client.Client, name string) string {
	// 所有查询从 context 获取客户端设置, 如未设置将使用默认客户端。
	var ctx = context.Background()
	ctx = client.With(ctx, c)
	imagebase64, err := Searchimage(ctx, name)
	if err != nil {
		log.Println("请求出错")
	}
	return imagebase64
}

func Searchimage(ctx context.Context, name string) (string, error) {
	// 搜索画作
	result, err := artwork.Search(ctx, name)
	if err != nil {
		log.Println(err)
		return "", err
	}
	//fmt.Println(result.JSON)                                        // json return data.
	artwork := result.Artworks() // []artwork.Artwork，只有部分数据，通过 `Fetch` `FetchPages` 方法获取完整数据。\
	//

	if len(artwork)-1 <= 0 {
		return "", errors.New("请求参数为0")
	}
	randomNumber := rand.Intn(len(artwork) - 1)
	urls := artwork[randomNumber].Image
	original_url := match(urls.Thumb)
	if original_url == "" {
		return "", errors.New("匹配出错")
	}
	// 生成0到n-1之间的随机整数
	imagebase64 := RequestImage(original_url, artwork[randomNumber].Title)
	//artwork.Search(ctx, "パチュリー・ノーレッジ", artwork.SearchOptionPage(2)) // 获取第二页
	return imagebase64, nil

}

func GetRankimage(ctx context.Context) (string, error) {
	rank := &artwork.Rank{Mode: "weekly"}
	err := rank.Fetch(ctx)
	if err != nil {
		log.Println(err)
		return "", err
	}

	randomNumber := rand.Intn(len(rank.Items) - 1)
	urls := rank.Items[randomNumber].Image
	original_url := match(urls.Regular)
	if original_url == "" {
		return "", errors.New("匹配出错")
	}
	imagebase64 := RequestImage(original_url, rank.Items[randomNumber].ID)
	return imagebase64, nil
}

func RequestImage(url string, name string) string {
	var body []byte
	req, err := http.NewRequest("GET", url, bytes.NewReader(body))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Referer", "https://www.pixiv.net/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request:", err)
	}
	defer resp.Body.Close()
	// file, err := os.Create("./image/" + name + ".jpg")
	// if err != nil {
	// 	log.Fatalf("创建文件失败: %v", err)
	// }
	// defer file.Close()
	// _, err = io.Copy(file, resp.Body)
	// if err != nil {
	// 	log.Println("保存文件失败", err)
	// }
	imagebyte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取字节流出错", err)
	}
	imagebase64 := base64.StdEncoding.EncodeToString(imagebyte)
	return imagebase64
}

func match(url string) string {
	pattern := `img/(\d{4}/\d{2}/\d{2}/\d{2}/\d{2}/\d{2}/\d+)_p\d+`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(url)
	if match != nil {
		// 输出匹配的结果
		fmt.Println(match) // 输出: 2024/08/14/03/06/07/121468483
		return "https://i.pximg.net/img-original/" + match[0] + ".jpg"
	} else {
		fmt.Println("No match found")
		return ""
	}
}
