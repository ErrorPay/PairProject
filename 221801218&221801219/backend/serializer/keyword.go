package serializer

import "backend/model"

type KeywordStat struct {
	KeywordStatUnits []KeywordStatUnit `json:"keyword_stat_units"`
}

type KeywordStatUnit struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Year    uint16 `json:"year"`
	Freq    uint64 `json:"freq"`
	Meeting string `json:"meeting"`
}

func BuildKeywordStat(keywords []model.Keyword) KeywordStat {
	var keywordStat KeywordStat
	cnt := 1
	for _, keyword := range keywords {
		keywordStat.KeywordStatUnits = append(keywordStat.KeywordStatUnits, KeywordStatUnit{
			ID:      cnt,
			Content: keyword.Content,
			Year:    keyword.Year,
			Freq:    keyword.Freq,
			Meeting: keyword.Meeting,
		})
		cnt = cnt + 1
	}

	return keywordStat
}

func BuildKeywordStatResponse(keywords []model.Keyword) Response {
	return Response{
		Code: 0,
		Data: BuildKeywordStat(keywords),
		Msg:  "Success",
	}
}
