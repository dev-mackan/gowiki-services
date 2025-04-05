package reposervice

import "fmt"

type NotFoundError struct {
	err string
}

func (e *NotFoundError) Error() string {
	if e.err == "" {
		return fmt.Sprintf("Resource not found")
	}
	return fmt.Sprintf(e.err)
}

type DatabaseError struct {
	err string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("%s", e.err)
}
