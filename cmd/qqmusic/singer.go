package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/xwjdsh/qqmusic"
)

func singerAction(client *qqmusic.Client, c *cli.Context) error {
	singerInfo, err := client.SearchSinger(strings.Join(c.Args().Slice(), " "))
	if err != nil {
		return err
	}

	if singerInfo == nil {
		return errors.New("singer not found!")
	}

	fansCount, err := client.GetSingerFansCount(singerInfo.CustomInfo.Singermid)
	if err != nil {
		return err
	}

	songList := []*qqmusic.Songinfo{}
	for page := 1; ; page += 1 {
		_, songs, err := client.GetSonglistBySinger(singerInfo.CustomInfo.Singermid, page, 50)
		if err != nil {
			return err
		}
		if len(songs) == 0 {
			break
		}

		songMids := []string{}
		songIds := []int{}
		for _, song := range songs {
			info := song.Info
			songMids = append(songMids, info.Mid)
			songIds = append(songIds, info.ID)
		}

		commentCountMap, err := client.GetSongCommentCount(songIds)
		if err != nil {
			return err
		}

		favorCountMap, err := client.GetSongFavorCount(songMids)
		if err != nil {
			return err
		}

		for _, song := range songs {
			song.CommnetCount = commentCountMap[strconv.Itoa(song.Info.ID)]
			song.FavorCount = favorCountMap[song.Info.Mid]
		}

		if c.Bool("solo") {
			for _, song := range songs {
				if len(song.Info.Singer) == 1 {
					songList = append(songList, song)
				}
			}
		} else {
			songList = append(songList, songs...)
		}
	}

	data := [][]string{
		{
			fmt.Sprintf("%s https://y.qq.com/n/ryqq/singer/%s", singerInfo.Title, singerInfo.CustomInfo.Singermid),
			strconv.Itoa(fansCount),
			singerInfo.CustomInfo.AlbumNum,
			singerInfo.CustomInfo.SongNum,
			singerInfo.CustomInfo.MvNum,
		},
	}

	singerTable := tablewriter.NewWriter(os.Stdout)
	singerTable.SetHeader([]string{"Name", "Fans", "Album", "Song", "MV"})
	singerTable.AppendBulk(data)

	sort.Slice(songList, func(i, j int) bool {
		if c.String("order") == "comment" {
			return songList[i].CommnetCount > songList[j].CommnetCount
		}
		return songList[i].FavorCount > songList[j].FavorCount
	})

	songData := [][]string{}
	for i := 0; i < c.Int("count"); i++ {
		song := songList[i]
		info := song.Info
		songData = append(songData,
			[]string{
				fmt.Sprintf("<%s> https://y.qq.com/n/ryqq/songDetail/%s", info.Title, info.Mid),
				fmt.Sprintf("<%s> https://y.qq.com/n/ryqq/albumDetail/%s", info.Album.Title, info.Album.Mid),
				strconv.Itoa(song.CommnetCount),
				strconv.Itoa(song.FavorCount),
			},
		)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Album", "Comment", "Favor"})
	table.SetRowLine(true) // Enable row line

	table.SetCenterSeparator("*")
	table.SetColumnSeparator("╪")
	table.SetRowSeparator("-")

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(songData)

	singerTable.Render()
	fmt.Println()
	table.Render()
	return nil
}
