package request

func RedirectResponse(redirect string) map[string]string {
	return map[string]string{
		"redirect": redirect,
	}
}
