package graph

func ComplexityConfig() ComplexityRoot {
	var c ComplexityRoot

	c.Query.User = func(userComplexity int) int {
		return 56
	}

	c.Query.Achievements = func(userComplexity int) int {
		return 14
	}

	c.Query.AchievementCategories = func(userComplexity int) int {
		return 9
	}

	c.User.Achievements = func(userComplexity int) int {
		return 13
	}

	c.User.Avatar = func(userComplexity int) int {
		return 14
	}

	c.User.Chills = func(userComplexity int) int {
		return 26
	}

	c.Achievement.Category = func(userComplexity int) int {
		return 13
	}

	c.AchievementCategory.Achievements = func(userComplexity int) int {
		return 9
	}

	return c
}
