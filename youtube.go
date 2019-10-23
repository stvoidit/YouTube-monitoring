package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type MetadataRenderer struct {
	ID    string `json:"videoId"`
	Title struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"title"`
	Views struct {
		VCR struct {
			ViewCount struct {
				Count string `json:"simpleText"`
			} `json:"viewCount"`
		} `json:"videoViewCountRenderer"`
	} `json:"viewCount"`
	Likes struct {
		LBR struct {
			LikeCount    int64 `json:"likeCount"`
			DislikeCount int64 `json:"dislikeCount"`
		} `json:"likeButtonRenderer"`
	} `json:"likeButton"`
}

type YouTube struct {
	AutoPlay int64  `json:"autoplay_count"`
	RVC      string `json:"rvs"`
	RawWNR   string `json:"watch_next_response"`
	WNR      struct {
		RContext struct {
			TCWR struct {
				Res1 struct {
					Res2 struct {
						Contents []struct {
							SectionRenderer struct {
								Contents []struct {
									MetadataRenderer *MetadataRenderer `json:"videoMetadataRenderer"`
								} `json:"contents"`
							} `json:"itemSectionRenderer"`
						} `json:"contents"`
					} `json:"results"`
				} `json:"results"`
			} `json:"twoColumnWatchNextResults"`
		} `json:"contents"`
	}
}

type CorrentVideo struct {
	ID       string
	Views    int64
	Likes    int64
	Dislikes int64
	Title    string
}

func ParseYoutube(html []byte) (*CorrentVideo, error) {
	pattern := regexp.MustCompile(`'RELATED_PLAYER_ARGS': (.*),\n`)
	data := pattern.FindSubmatch(html)
	replacer := strings.NewReplacer(`//`, ``, `""`, `"`)
	cv := new(CorrentVideo)
	if len(data) > 1 {
		youtube := new(YouTube)
		js := replacer.Replace(string(data[1]))
		json.Unmarshal([]byte(js), &youtube)
		json.Unmarshal([]byte(youtube.RawWNR), &youtube.WNR)

		clearPatterns := regexp.MustCompile(`[^\d]+`)
		clearViews := clearPatterns.ReplaceAll([]byte(youtube.WNR.RContext.TCWR.Res1.Res2.Contents[0].SectionRenderer.Contents[0].MetadataRenderer.Views.VCR.ViewCount.Count), []byte(""))
		views, _ := strconv.ParseInt(string(clearViews), 10, 64)

		cv.ID = youtube.WNR.RContext.TCWR.Res1.Res2.Contents[0].SectionRenderer.Contents[0].MetadataRenderer.ID
		cv.Views = views
		cv.Likes = youtube.WNR.RContext.TCWR.Res1.Res2.Contents[0].SectionRenderer.Contents[0].MetadataRenderer.Likes.LBR.LikeCount
		cv.Dislikes = youtube.WNR.RContext.TCWR.Res1.Res2.Contents[0].SectionRenderer.Contents[0].MetadataRenderer.Likes.LBR.DislikeCount
		cv.Title = youtube.WNR.RContext.TCWR.Res1.Res2.Contents[0].SectionRenderer.Contents[0].MetadataRenderer.Title.Runs[0].Text
		return cv, nil
	}
	return cv, errors.New("can't parse")
}

func main() {
	url := "..."
	client := &http.Client{}
	r, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	youtube, err := ParseYoutube(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(youtube)

}
