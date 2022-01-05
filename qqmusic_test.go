package qqmusic

import (
	"os"
	"testing"
)

func getCookie() string {
	return os.Getenv("QQMUSIC_COOKIE")
}

func TestSearchSinger(t *testing.T) {
	cookie := getCookie()
	if cookie == "" {
		t.Skip()
	}
	singerInfo, err := New(cookie).SearchSinger("美波")
	if err != nil || singerInfo == nil {
		t.Error(err, singerInfo)
	}
}

func TestGetSingerFansNum(t *testing.T) {
	cookie := getCookie()
	if cookie == "" {
		t.Skip()
	}
	count, err := New(cookie).GetSingerFansCount("003fA5G40k6hKc")
	if err != nil || count == 0 {
		t.Error(err, count)
	}
}

func TestGetSonglistBySinger(t *testing.T) {
	cookie := getCookie()
	if cookie == "" {
		t.Skip()
	}
	totalNum, songInfos, err := New(cookie).GetSonglistBySinger("003fA5G40k6hKc", 1, 2)
	if err != nil || len(songInfos) == 0 || totalNum == 0 {
		t.Error(err, songInfos, totalNum)
	}
}

func TestGetSongCommentCount(t *testing.T) {
	cookie := getCookie()
	if cookie == "" {
		t.Skip()
	}
	countMap, err := New(cookie).GetSongCommentCount([]int{311084300, 203650936})
	if err != nil || len(countMap) == 0 {
		t.Error(err, countMap)
	}
}

func TestGetSongFavorCount(t *testing.T) {
	cookie := getCookie()
	if cookie == "" {
		t.Skip()
	}
	countMap, err := New(cookie).GetSongFavorCount([]string{"004Q9vXu2FuqDs", "0045YdtG4FSRLN"})
	if err != nil || len(countMap) == 0 {
		t.Error(err, countMap)
	}
}
