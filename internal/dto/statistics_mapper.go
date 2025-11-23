package dto

import "Git_PR_plugin_avito_tech/internal/model"

func ToUserStatistics(data *model.UserStatisticsData) UserStatistics {
	return UserStatistics{
		UserID:               data.UserID,
		Username:             data.Username,
		TeamName:             data.TeamName,
		TotalAssignments:     data.TotalAssignments,
		OpenAssignments:      data.OpenAssignments,
		CompletedAssignments: data.CompletedAssignments,
	}
}

func ToUserStatisticsList(dataList []*model.UserStatisticsData) []UserStatistics {
	result := make([]UserStatistics, len(dataList))
	for i, data := range dataList {
		result[i] = ToUserStatistics(data)
	}
	return result
}

func ToPullRequestStatistics(data *model.PRStatisticsData) PullRequestStatistics {
	return PullRequestStatistics{
		PullRequestID:   data.PullRequestID,
		PullRequestName: data.PullRequestName,
		AuthorUsername:  data.AuthorUsername,
		Status:          data.Status,
		ReviewersCount:  data.ReviewersCount,
	}
}

func ToPullRequestStatisticsList(dataList []*model.PRStatisticsData) []PullRequestStatistics {
	result := make([]PullRequestStatistics, len(dataList))
	for i, data := range dataList {
		result[i] = ToPullRequestStatistics(data)
	}
	return result
}
