/*
 * @Author: your name
 * @Date: 2020-09-22 11:18:00
 * @LastEditTime: 2020-09-28 15:29:12
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\db.go
 */
package msg

import "time"

type Msg struct {
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

type CustomWhere struct {
	Query interface{}
	Args  []interface{}
}

type Messages interface {
	//get msgs for read,sort by senddate
	GetUnread(where CustomWhere, start int, count int) ([]Msg, error)
	//get msg for uniqueid
	GetIndex(id int) (*Msg, error)
	//get all msgs,sort by senddate
	GetAll(where CustomWhere, start int, count int) ([]Msg, error)
	//set read
	MarkRead(idx int) error

	GetUnradCount(where CustomWhere) (int, error)
	GetCount(where CustomWhere) (int, error)

	//...
}
