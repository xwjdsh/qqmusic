package qqmusic

import (
	"testing"
)

func TestSearchSinger(t *testing.T) {
	singerInfo, err := New().SearchSinger("美波")
	if err != nil || singerInfo == nil {
		t.Error(err, singerInfo)
	}
}

func TestGetSingerFansNum(t *testing.T) {
	count, err := New().GetSingerFansCount("003fA5G40k6hKc")
	if err != nil || count == 0 {
		t.Error(err, count)
	}
}

func TestGetSonglistBySinger(t *testing.T) {
	totalNum, songInfos, err := New().GetSonglistBySinger("003fA5G40k6hKc", 1, 2)
	if err != nil || len(songInfos) == 0 || totalNum == 0 {
		t.Error(err, songInfos, totalNum)
	}
}

func TestGetSongCommentCount(t *testing.T) {
	countMap, err := New().GetSongCommentCount([]int{311084300, 203650936})
	if err != nil || len(countMap) == 0 {
		t.Error(err, countMap)
	}
}

func TestGetSongFavorCount(t *testing.T) {
	countMap, err := New().GetSongFavorCount([]string{"004Q9vXu2FuqDs", "0045YdtG4FSRLN"})
	if err != nil || len(countMap) == 0 {
		t.Error(err, countMap)
	}
}
