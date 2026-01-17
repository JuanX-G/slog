package database

type NotFoundError struct {
	queryParam any
}

func (e NotFoundError) Error() string {
	return e.queryParam.(string)
}


