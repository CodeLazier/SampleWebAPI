/*
 * @Author: your name
 * @Date: 2020-09-24 10:35:54
 * @LastEditTime: 2020-10-01 21:41:46
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\msg\orm.go
 */
package handler

import "time"

type UpdateInfo struct {
	Field string
	Value interface{}
}

type Cmd struct {
	Model interface{} // query struct or model
	//Select
	Start int    // start offset
	Count int    // get count records
	Order string // order
	//add others if business needs
	Query interface{}
	Args  []interface{}

	//update
	Update UpdateInfo

	//if calc count
	CalcCount bool
}

type EipMsg struct {
	//gorm.Model        // orm map
	Id       int       `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	Title    string    `json:"title" gorm:"column:title"`
	Content  string    `json:"content" gorm:"column:content"`
	CreateAt time.Time `json:"createAt" gorm:"column:createAt;autoCreateTime"`

	//...
}

func (m EipMsg) TableName() string {
	return "messages"
}

func NewRecord(query interface{}, args ...interface{}) Cmd {
	a := make([]interface{}, 0)
	a = append(a, args...)
	return Cmd{
		Model: &EipMsg{},
		Query: query,
		Args:  a,
	}
}

func NewMultiRecords(start, count int) Cmd {
	var msgs []EipMsg
	return Cmd{
		Model: &msgs,
		Start: start,
		Count: count,
		Order: "\"id\" desc", // ID为自增字段,createAt为创建时间,所以createAt可能会有重复(大并发量),但Id不会
	}
}

func NewUpdateRecord(field string, value interface{}, query interface{}, args ...interface{}) Cmd {
	a := make([]interface{}, 0)
	a = append(a, args...)
	return Cmd{
		Model:  &EipMsg{},
		Query:  query,
		Args:   a,
		Update: UpdateInfo{Field: field, Value: value},
	}
}

func NewUpdateMsg(msg EipMsg) Cmd {
	return Cmd{
		Model:  &EipMsg{},
		Update: UpdateInfo{Field: "^All", Value: msg},
	}
}

func NewInsertMsg(msg interface{}) Cmd {
	return Cmd{
		Model: &EipMsg{},
		Update: UpdateInfo{
			Field: "^Insert",
			Value: msg,
		},
	}
}

type Control interface {
	OpenOrm(...string) error
	Query(cmd Cmd) (interface{}, error)
	Update(cmd Cmd) error
	Insert(cmd Cmd) (interface{}, error)
}
