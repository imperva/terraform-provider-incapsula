package incapsula

// Seeders the Condition tests need on top of the shared fake.

// abpConditionSeedManaged stores a managed (Imperva-owned) literal Condition.
// Such a Condition carries no account, and the API withholds its expression,
// its template and its timestamps, reporting the Unix epoch instead.
func abpConditionSeedManaged(f *abpFake, name, description string, recommendedDirectives ...string) string {
	recommended := make([]any, 0, len(recommendedDirectives))
	for _, directive := range recommendedDirectives {
		recommended = append(recommended, directive)
	}
	return f.seed(abpFakeConditions, map[string]any{
		"account_id":             "",
		"name":                   name,
		"description":            description,
		"code":                   "",
		"template":               "",
		"last_change_by":         "",
		"recommended_directives": recommended,
		"created_at":             "1970-01-01T00:00:00Z",
		"modified_at":            "1970-01-01T00:00:00Z",
	})
}
