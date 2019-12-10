package db

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (d *DB) CleanUp(t *time.Time) {
	rows := d.Exec(`
		DELETE FROM users
		WHERE id IN(
			SELECT id from (
				select u.id id, u.created_at, u.role role, count(r.user_id) items
				from users u 
				left join routes r on r.user_id = u.id
				group by u.id
			) t1 
			WHERE 
				role = "guest" 
				AND items = 0 
				AND created_at < NOW() - INTERVAL 1 week
		)
	`).RowsAffected

	if rows > 0 {
		logrus.Info(fmt.Sprintf("DB: Cleaned %v empty guests", rows))
	}
}
