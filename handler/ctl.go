/*
 * @Author: your name
 * @Date: 2020-09-24 10:35:54
 * @LastEditTime: 2020-09-30 11:44:24
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
	UniqueID        int       `json:"id" gorm:"column:UniqueID"`
	GroupID         string    `json:"groupid" gorm:"column:GroupID"`
	PGroupID        string    `json:"pgroupid" gorm:"column:PGroupID"`
	DeptID          string    `json:"deptid" gorm:"column:DeptID"`
	FromID          string    `json:"fromid" gorm:"column:FromID"`
	FromName        string    `json:"fromname" gorm:"column:FromName"`
	ToID            string    `json:"toid" gorm:"column:ToID"`
	ToName          string    `json:"toname" gorm:"column:ToName"`
	ToSavedValue    string    `json:"tosavedvalue" gorm:"column:ToSavedValue"`
	ToForceReplyTAG int       `json:"toforcereplytag" gorm:"column:ToForceReplyTAG"`
	CCID            string    `json:"ccid" gorm:"column:CCID"`
	CCName          string    `json:"ccname" gorm:"column:CCName"`
	CCSavedValue    string    `json:"ccsavedvalue" gorm:"column:CCSavedValue"`
	CCForceReply    int       `json:"ccforcereplytag" gorm:"column:CCForceReplyTAG"`
	BCCID           string    `json:"bccid" gorm:"column:BCCID"`
	BCCName         string    `json:"bccname" gorm:"column:BCCName"`
	BCCSavedValue   string    `json:"bccsavedvalue" gorm:"column:BCCSavedValue"`
	SendDate        time.Time `json:"senddate" gorm:"column:SendDate"`

	Subject   string `json:"subject" gorm:"column:Subject"`
	Content   string `json:"content" gorm:"column:Content"`
	Html      string `json:"htmltag" gorm:"column:HtmlTAG"`
	Draft     string `json:"drafttag" gorm:"column:DraftTAG"`
	Del       string `json:"deltag" gorm:"column:DelTAG"`
	CompleteT int    `json:"completetag" gorm:"column:CompleteTAG"`
	ReplyID   int    `json:"replyid" gorm:"column:ReplyID"`
	Read      int    `json:"readtag" gorm:"column:ReadTAG"`

	//...
}

func (m EipMsg) TableName() string {
	return "EIP_MessageMaster"
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
		Order: "SendDate desc",
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

type Control interface {
	OpenOrm(...string) error
	Query(cmd Cmd) (interface{}, error)
	Update(cmd Cmd) error
}
