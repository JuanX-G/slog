package database

var (
    AllowedTables = map[string][]string{
        "users": {"id", "user_name", "email", "date_created", "password", "description"},
        "posts": {"id", "author_id", "content", "date_created", "title", "tags", "likes"},
        "post_likes": {"user_id", "post_id", "created_at"},
    }
)

func isTableAllowed(table string) bool {
    _, ok := AllowedTables[table]
    return ok
}

func areColumnsAllowed(table string, cols []string) bool {
    allowed, ok := AllowedTables[table]
    if !ok {
        return false
    }

    allowedSet := make(map[string]struct{}, len(allowed))
    for _, c := range allowed {
        allowedSet[c] = struct{}{}
    }

    for _, c := range cols {
        if _, exists := allowedSet[c]; !exists {
            return false
        }
    }
    return true
}

