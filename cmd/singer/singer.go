package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/xwjdsh/qqmusic"
)

var (
	songCount = flag.Int("count", 10, "song count")
	orderBy   = flag.String("order", "favor", "order by [favor|comment]")
)

func main() {
	flag.Parse()
	args := flag.Args()
	exitIf(len(args) == 0, "singer name is required!")

	c := qqmusic.New()

	singerInfo, err := c.SearchSinger(strings.Join(args, " "))
	exitIfErr(err)
	exitIf(singerInfo == nil, "singer not found!")

	fansCount, err := c.GetSingerFansCount(singerInfo.Singermid)
	exitIfErr(err)

	songList := []*qqmusic.Songinfo{}
	for page := 1; ; page += 1 {
		total, songs, err := c.GetSonglistBySinger(singerInfo.Singermid, page, 50)
		exitIfErr(err)

		songMids := []string{}
		songIds := []int{}
		for _, song := range songs {
			info := song.Info
			songMids = append(songMids, info.Mid)
			songIds = append(songIds, info.ID)
		}

		commentCountMap, err := c.GetSongCommentCount(songIds)
		exitIfErr(err)

		favorCountMap, err := c.GetSongFavorCount(songMids)
		exitIfErr(err)

		for _, song := range songs {
			song.CommnetCount = commentCountMap[strconv.Itoa(song.Info.ID)]
			song.FavorCount = favorCountMap[song.Info.Mid]
		}

		songList = append(songList, songs...)
		if len(songList) >= total {
			break
		}
	}

	data := [][]string{
		{
			fmt.Sprintf("%s https://y.qq.com/n/ryqq/singer/%s", singerInfo.Singername, singerInfo.Singermid),
			strconv.Itoa(fansCount),
			strconv.Itoa(singerInfo.Albumnum),
			strconv.Itoa(singerInfo.Songnum),
			strconv.Itoa(singerInfo.Mvnum),
		},
	}

	singerTable := tablewriter.NewWriter(os.Stdout)
	singerTable.SetHeader([]string{"Name", "Fans", "Album", "Song", "MV"})
	singerTable.AppendBulk(data)

	sort.Slice(songList, func(i, j int) bool {
		if *orderBy == "comment" {
			return songList[i].CommnetCount > songList[j].CommnetCount
		}
		return songList[i].FavorCount > songList[j].FavorCount
	})

	songData := [][]string{}
	for i := 0; i < *songCount; i++ {
		song := songList[i]
		info := song.Info
		songData = append(songData,
			[]string{
				fmt.Sprintf("<%s> https://y.qq.com/n/ryqq/songDetail/%s", info.Name, info.Mid),
				fmt.Sprintf("<%s> https://y.qq.com/n/ryqq/albumDetail/%s", info.Album.Name, info.Mid),
				strconv.Itoa(song.CommnetCount),
				strconv.Itoa(song.FavorCount),
			},
		)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Album", "Commnet", "Favor"})
	table.SetRowLine(true) // Enable row line

	table.SetCenterSeparator("*")
	table.SetColumnSeparator("╪")
	table.SetRowSeparator("-")

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(songData)

	singerTable.Render()
	fmt.Println()
	table.Render()
}

func exitIf(b bool, msg interface{}) {
	if b {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func exitIfErr(err error) {
	exitIf(err != nil, err)
}
