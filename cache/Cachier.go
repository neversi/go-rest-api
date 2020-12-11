package cache

import "gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"

// PostCache ...
type PostCache interface{
	Set(key string, value []*models.Task) error
	Get(key string) ([]*models.Task, error)
}