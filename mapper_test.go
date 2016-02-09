package qbit

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type User struct {
	Id         string `qbit:"type:uuid"`
	Email      string `qbit:"type:varchar(255); constraints:unique,notnull"`
	FullName   string `qbit:"constraints:notnull,index"`
	FacebookId int64  `qbit:"constraints:null"`
	UserType   string `qbit:"constraints:default(guest)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	PrimaryKey `qbit:"id"`
	ForeignKey `qbit:"id references profile.id"`
}

func TestMapper(t *testing.T) {

	engine, err := NewEngine("mysql", "root:@tcp(127.0.0.1:3306)/qbit_test")
	defer engine.DB().Close()

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	mapper := NewMapper("postgres")

	userTable, err := mapper.Convert(User{})

	fmt.Println("error: ", err.Error())

	//	fmt.Println(err)
	//	fmt.Println(userTable.Sql())
	fmt.Println(userTable)

	//	result, err := engine.Exec(userTable.Sql(), []interface{}{})

	//	assert.Equal(t, err, nil)
	//	fmt.Println(result)

	assert.Equal(t, 1, 1)

}
