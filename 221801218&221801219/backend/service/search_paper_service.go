package service

import (
	"backend/model"
	"backend/serializer"
	"backend/util"
	"regexp"
	"strconv"
)

type SearchPaperService struct {
}

func (service *SearchPaperService) Search(title, keyword string, page int64, meeting string) serializer.Response {
	if title == "" && keyword == "" {
		return serializer.ParamErr("参数错误", nil)
	}

	if _, err := model.Engine.Exec("CREATE TEMPORARY TABLE `temp_table` (`id` INTEGER PRIMARY KEY NOT NULL, `title` TEXT NOT NULL, `abstract` TEXT NOT NULL, `meeting` TEXT NOT NULL, `year` INTEGER NOT NULL, `origin_link` TEXT NOT NULL, `click` INTEGER DEFAULT 0 NOT NULL);"); err != nil {
		return serializer.DBErr("Creating temporary table", err)
	}

	re := regexp.MustCompile(`\s`)
	t := re.Split(title, -1)
	k := re.Split(keyword, -1)

	var sql string
	for _, i := range k {
		sql = "INSERT OR IGNORE INTO temp_table SELECT * FROM paper WHERE paper.id IN (SELECT paper_keyword.paper_id FROM paper_keyword WHERE paper_keyword.keyword_id IN (SELECT keyword.id FROM keyword WHERE keyword.content LIKE '%" + i + "%'));"

		if i == "" {
			continue
		}
		_, err := model.Engine.Exec(sql)
		if err != nil {
			util.Log().Error(err.Error())
			return serializer.ParamErr("0", err)
		}
	}

	if title != "" {
		for _, tt := range t {
			sql = "INSERT OR IGNORE INTO temp_table SELECT * FROM paper WHERE title LIKE '%" + tt + "%'"
			sqlDelete := "DELETE FROM temp_table WHERE title NOT LIKE '%" + tt + "%'"
			if tt == "" {
				continue
			}

			if keyword == "" {
				_, err := model.Engine.Exec(sql)
				if err != nil {
					util.Log().Error(err.Error())
					return serializer.ParamErr("2", err)
				}
			}

			_, err := model.Engine.Exec(sqlDelete)
			if err != nil {
				util.Log().Error(err.Error())
				return serializer.DBErr("2.5", err)
			}
		}
	}

	if meeting != "" {
		sqlDelete := "DELETE * FROM temp_table WHERE meeting != ?"
		_, err := model.Engine.Exec(sqlDelete, meeting)
		if err != nil {
			util.Log().Error(err.Error())
			return serializer.ParamErr("meeting", err)
		}
	}

	result, _ := model.Engine.Query("SELECT DISTINCT Count(id) FROM temp_table")
	total, _ := strconv.Atoi(string(result[0]["Count(id)"]))
	pageCount := util.TotalPages(int64(total))
	if util.PageOverFlow(int64(total), page) {
		return serializer.ParamErr("没有结果", nil)
	}

	var papers []model.Paper
	err := model.Engine.Table("temp_table").Distinct().Limit(util.PaperPageMaxSize, int(util.PaperPageMaxSize*(page-1))).Asc("id").Find(&papers)
	if err != nil {
		util.Log().Error(err.Error())
		return serializer.ParamErr("3", err)
	}

	_, _ = model.Engine.Exec("drop table temp.temp_table")

	return serializer.BuildPaperListResponse(papers, pageCount, page)
}
