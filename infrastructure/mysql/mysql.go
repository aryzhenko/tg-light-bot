package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"tg-light-bot/common"
	"tg-light-bot/domain"
)

//https://zetcode.com/golang/mysql/

type Connection struct {
	conn *sql.DB
}

func NewConnection(config common.DbConfig) *Connection {
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=%s", config.User, config.Password, config.Host, config.DbName, url.QueryEscape("Europe/Kiev"))
	sqlC, err := sql.Open("mysql", connectStr)
	if err != nil {
		panic(err)
	}
	return &Connection{conn: sqlC}
}

func (c *Connection) SaveUser(user domain.User) error {
	if user.IsNew {
		return c.insertUser(user)
	} else {
		return c.updateUser(user)
	}
}

func (c *Connection) insertUser(user domain.User) error {
	_, err := c.conn.Query("insert into users "+
		"(id,first_name,last_name,user_name,is_bot,language_code,can_join_groups,can_read_all_group_messages,supports_inline_queries,created_at,last_activity) "+
		"values (?,?,?,?,?,?,?,?,?,?,?)",
		user.ID,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.IsBot,
		user.LanguageCode,
		user.CanJoinGroups,
		user.CanReadAllGroupMessages,
		user.SupportsInlineQueries,
		user.CreatedAt,
		user.LastActivity,
	)
	return err
}

func (c *Connection) updateUser(user domain.User) error {
	_, err := c.conn.Query("update users set "+
		"first_name=?,"+
		"last_name=?,"+
		"user_name=?,"+
		"is_bot=?,"+
		"language_code=?,"+
		"can_join_groups=?,"+
		"can_read_all_group_messages=?,"+
		"supports_inline_queries=?,"+
		"created_at=?,"+
		"last_activity=? "+
		"where id=?",
		user.FirstName,
		user.LastName,
		user.UserName,
		user.IsBot,
		user.LanguageCode,
		user.CanJoinGroups,
		user.CanReadAllGroupMessages,
		user.SupportsInlineQueries,
		user.CreatedAt,
		user.LastActivity,
		user.ID,
	)
	return err
}

func (c *Connection) GetUser(id int64) (*domain.User, error) {
	user := &domain.User{IsNew: false}
	query := "select id,first_name,last_name,user_name,is_bot,language_code,can_join_groups,can_read_all_group_messages,supports_inline_queries,created_at,last_activity " +
		"from users where id=?"
	row := c.conn.QueryRow(query, id)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.IsBot,
		&user.LanguageCode,
		&user.CanJoinGroups,
		&user.CanReadAllGroupMessages,
		&user.SupportsInlineQueries,
		&user.CreatedAt,
		&user.LastActivity,
	)

	return user, err
}

func (c *Connection) GetUserIds() ([]int64, error) {
	var result []int64

	query := "select id from users"
	res, err := c.conn.Query(query)
	if err != nil {
		return result, err
	}
	defer res.Close()

	for res.Next() {
		var item int64
		err := res.Scan(&item)
		if err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (c *Connection) GetLightState() (*domain.LightState, error) {
	result := domain.LightState{}
	query := "select is_on bool, changed_at from light_status order by id desc limit 1"
	row := c.conn.QueryRow(query)
	err := row.Scan(
		&result.IsOn,
		&result.ChangedAt,
	)

	return &result, err
}

func (c *Connection) SaveLightState(state *domain.LightState) error {
	query := "insert into light_status (is_on, changed_at) values (?, ?)"
	log.Println("Saving state: ", state)
	_, err := c.conn.Query(query, state.IsOn, state.ChangedAt)
	return err
}
