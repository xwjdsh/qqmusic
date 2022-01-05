package qqmusic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const baseUrl = "https://u.y.qq.com/cgi-bin/musicu.fcg"

type SingerInfo struct {
	Title      string `json:"title"`
	Desciption string `json:"desciption"`
	CustomInfo struct {
		AlbumNum  string `json:"album_num"`
		ExtraDesc string `json:"extra_desc"`
		IconType  string `json:"icon_type"`
		Ifpicurl  string `json:"ifpicurl"`
		IsFollow  string `json:"is_follow"`
		Singermid string `json:"mid"`
		MvNum     string `json:"mv_num"`
		SongNum   string `json:"song_num"`
	} `json:"custom_info"`
	TrackList struct {
		Items []struct {
			ID   int    `json:"id"`
			Mid  string `json:"mid"`
			Name string `json:"name"`
		} `json:"items"`
	} `json:"track_list"`
}

type SingerInfoResult struct {
	Code    int    `json:"code"`
	Ts      int64  `json:"ts"`
	StartTs int64  `json:"start_ts"`
	Traceid string `json:"traceid"`
	Result  struct {
		Code int `json:"code"`
		Data struct {
			Body struct {
				Zhida struct {
					List []*SingerInfo `json:"list"`
				} `json:"zhida"`
			} `json:"body"`
			Code int `json:"code"`
			Ver  int `json:"ver"`
		} `json:"data"`
	} `json:"result"`
}

type Songinfo struct {
	CommnetCount int `json:"-"`
	FavorCount   int `json:"-"`
	Info         struct {
		ID       int    `json:"id"`
		Type     int    `json:"type"`
		Mid      string `json:"mid"`
		Name     string `json:"name"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Singer   []struct {
			ID    int    `json:"id"`
			Mid   string `json:"mid"`
			Name  string `json:"name"`
			Title string `json:"title"`
			Type  int    `json:"type"`
			Uin   int    `json:"uin"`
			Pmid  string `json:"pmid"`
		} `json:"singer"`
		Album struct {
			ID         int    `json:"id"`
			Mid        string `json:"mid"`
			Name       string `json:"name"`
			Title      string `json:"title"`
			Subtitle   string `json:"subtitle"`
			TimePublic string `json:"time_public"`
			Pmid       string `json:"pmid"`
		} `json:"album"`
		Mv struct {
			ID    int    `json:"id"`
			Vid   string `json:"vid"`
			Name  string `json:"name"`
			Title string `json:"title"`
			Vt    int    `json:"vt"`
		} `json:"mv"`
		Interval    int      `json:"interval"`
		Isonly      int      `json:"isonly"`
		Language    int      `json:"language"`
		Genre       int      `json:"genre"`
		IndexCd     int      `json:"index_cd"`
		IndexAlbum  int      `json:"index_album"`
		TimePublic  string   `json:"time_public"`
		Status      int      `json:"status"`
		Fnote       int      `json:"fnote"`
		Label       string   `json:"label"`
		URL         string   `json:"url"`
		Bpm         int      `json:"bpm"`
		Version     int      `json:"version"`
		Trace       string   `json:"trace"`
		DataType    int      `json:"data_type"`
		ModifyStamp int      `json:"modify_stamp"`
		Pingpong    string   `json:"pingpong"`
		Aid         int      `json:"aid"`
		Ppurl       string   `json:"ppurl"`
		Tid         int      `json:"tid"`
		Ov          int      `json:"ov"`
		Sa          int      `json:"sa"`
		Es          string   `json:"es"`
		Vs          []string `json:"vs"`
	} `json:"songInfo"`
}

type Client struct {
	cookie string
}

func New(cookie string) *Client {
	return &Client{
		cookie: cookie,
	}
}

func (c *Client) GetSingerFansCount(singerMid string) (int, error) {
	module := "Concern.ConcernSystemServer"
	method := "cgi_qry_concern_num"
	param := map[string]interface{}{
		"vec_userinfo": []map[string]interface{}{
			{
				"usertype": 1,
				"userid":   singerMid,
			},
		},
	}
	data, err := c.get(module, method, param)
	if err != nil {
		return 0, err
	}
	var result struct {
		Code    int   `json:"code"`
		Ts      int64 `json:"ts"`
		StartTs int64 `json:"start_ts"`
		Result  struct {
			Code int `json:"code"`
			Data struct {
				MapSingerNum map[string]struct {
					SingerFollownum int `json:"singer_follownum"`
					UserFansnum     int `json:"user_fansnum"`
					UserFollownum   int `json:"user_follownum"`
				} `json:"map_singer_num"`
				MapUserNum struct {
				} `json:"map_user_num"`
			} `json:"data"`
		} `json:"result"`
	}

	fansCount := 0
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, err
	}
	if result.Code != 0 || result.Result.Code != 0 {
		return 0, fmt.Errorf("GetSingerFansCount error: %s", string(data))
	}

	if v, ok := result.Result.Data.MapSingerNum[singerMid]; ok {
		fansCount = v.UserFansnum
	}

	return fansCount, nil
}

func (c *Client) SearchSinger(name string) (*SingerInfo, error) {
	data, err := c.keywordSearch(name)
	if err != nil {
		return nil, err
	}

	var result SingerInfoResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 || result.Result.Data.Code != 0 {
		return nil, fmt.Errorf("SearchSinger error: %s", string(data))
	}

	if len(result.Result.Data.Body.Zhida.List) == 0 {
		return nil, fmt.Errorf("SearchSinger not found: %s", string(data))
	}
	info := result.Result.Data.Body.Zhida.List[0]
	desc := result.Result.Data.Body.Zhida.List[0].Desciption
	r := strings.Split(desc, " ")
	if len(r) == 3 {
		info.CustomInfo.SongNum = strings.Split(r[0], ":")[1]
		info.CustomInfo.AlbumNum = strings.Split(r[1], ":")[1]
		info.CustomInfo.MvNum = strings.Split(r[2], ":")[1]
	}

	return info, nil
}

func (c *Client) keywordSearch(keyword string) ([]byte, error) {
	module := "music.search.SearchCgiService"
	method := "DoSearchForQQMusicDesktop"
	param := map[string]interface{}{
		"remoteplace":  "txt.yqq.center",
		"search_type":  0,
		"query":        keyword,
		"page_num":     1,
		"num_per_page": 10,
	}

	return c.post(module, method, param)
}

func (c *Client) GetSonglistBySinger(singerMid string, page, pageSize int) (int, []*Songinfo, error) {
	module := "musichall.song_list_server"
	method := "GetSingerSongList"
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 100
	}
	param := map[string]interface{}{
		"singerMid": singerMid,
		"begin":     (page - 1) * pageSize,
		"num":       pageSize,
		"order":     1,
	}
	data, err := c.get(module, method, param)
	if err != nil {
		return 0, nil, err
	}
	var result struct {
		Code    int   `json:"code"`
		Ts      int64 `json:"ts"`
		StartTs int64 `json:"start_ts"`
		Result  struct {
			Code int `json:"code"`
			Data struct {
				Singermid string      `json:"singerMid"`
				Totalnum  int         `json:"totalNum"`
				Songlist  []*Songinfo `json:"songList"`
			} `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, nil, err
	}
	if result.Code != 0 || result.Result.Code != 0 {
		return 0, nil, fmt.Errorf("GetSonglistBySinger error: %s", string(data))
	}

	return result.Result.Data.Totalnum, result.Result.Data.Songlist, nil
}

func (c *Client) GetSongCommentCount(songIds []int) (map[string]int, error) {
	module := "GlobalComment.GlobalCommentReadServer"
	method := "GetCommentCount"
	requestList := []map[string]interface{}{}
	for _, songId := range songIds {
		requestList = append(requestList, map[string]interface{}{
			"biz_type": 1,
			"biz_id":   strconv.Itoa(songId),
		})
	}
	param := map[string]interface{}{"request_list": requestList}

	data, err := c.get(module, method, param)
	if err != nil {
		return nil, err
	}
	var result struct {
		Code    int   `json:"code"`
		Ts      int64 `json:"ts"`
		StartTs int64 `json:"start_ts"`
		Result  struct {
			Code int `json:"code"`
			Data struct {
				ResponseList []struct {
					BizID      string `json:"biz_id"`
					BizSubType int    `json:"biz_sub_type"`
					BizType    int    `json:"biz_type"`
					Count      int    `json:"count"`
				} `json:"response_list"`
			} `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 || result.Result.Code != 0 {
		return nil, fmt.Errorf("GetSongCommentCount error: %s", string(data))
	}

	m := map[string]int{}
	for _, r := range result.Result.Data.ResponseList {
		m[r.BizID] = r.Count
	}
	return m, nil
}

func (c *Client) GetSongFavorCount(songMids []string) (map[string]int, error) {
	module := "music.musicasset.SongFavRead"
	method := "GetSongFansNumberByMid"
	param := map[string]interface{}{"v_songMid": songMids}

	data, err := c.post(module, method, param)
	if err != nil {
		return nil, err
	}
	var result struct {
		Code    int   `json:"code"`
		Ts      int64 `json:"ts"`
		StartTs int64 `json:"start_ts"`
		Result  struct {
			Code int `json:"code"`
			Data struct {
				MNumbers map[string]int `json:"m_numbers"`
			} `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 || result.Result.Code != 0 {
		return nil, fmt.Errorf("GetSongFavorCount error: %s", string(data))
	}

	return result.Result.Data.MNumbers, nil
}

func (c *Client) get(module, method string, param map[string]interface{}) ([]byte, error) {
	queryUrl := ""
	if module == "url" {
		values := url.Values{}
		for k, v := range param {
			if vs, ok := v.(string); ok {
				values.Set(k, vs)
			}
		}
		queryUrl = method + "?" + values.Encode()
	} else {
		d := map[string]interface{}{
			"result": map[string]interface{}{
				"module": module,
				"method": method,
				"param":  param,
			},
		}
		data, err := json.Marshal(d)
		if err != nil {
			return nil, err
		}

		values := url.Values{}
		values.Set("data", string(data))
		values.Set("inCharset", "utf8")
		values.Set("outCharset", "utf-8")
		values.Set("format", "json")
		queryUrl = baseUrl + "?" + values.Encode()
	}

	resp, err := http.Get(queryUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c *Client) post(module, method string, param map[string]interface{}) ([]byte, error) {
	d := map[string]interface{}{
		"result": map[string]interface{}{
			"module": module,
			"method": method,
			"param":  param,
		},
	}
	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("cookie", c.cookie)
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
