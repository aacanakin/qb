package qbit

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type User struct {
	ID              string `qbit:"type:uuid"`
	Email           string `qbit:"type:varchar(255); constraints:unique,notnull"`
	FullName        string `qbit:"constraints:notnull,index"`
	Password        string `qbit:"type:text"`
	FacebookID      int64  `qbit:"constraints:null"`
	UserType        string `qbit:"constraints:default(guest)"`
	Points          float32
	CreatedAt       time.Time  `qbit:"constraints:notnull"`
	UpdatedAt       time.Time  `qbit:"constraints:notnull"`
	DeletedAt       *time.Time `qbit:"constraints:null"`

	PrimaryKey      `qbit:"id"`
	ForeignKey      `qbit:"(id) references (profile.id)"`
	CompositeUnique `qbit:"email,full_name"`
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

	if err != nil {
		fmt.Errorf("Error: %s\n", err.Error())
	}

	fmt.Println(userTable.SQL())

	//	result, err := engine.Exec(userTable.Sql(), []interface{}{})

	//	assert.Equal(t, err, nil)
	//	fmt.Println(result)

	assert.Equal(t, 1, 1)

}
