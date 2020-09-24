/*
 * @Author: your name
 * @Date: 2020-09-22 11:18:00
 * @LastEditTime: 2020-09-24 14:59:38
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\db.go
 */
package msg

type Msg struct {
	//gorm.Model        // orm map
	Id      int    `json:"id"`
	Title   string `json:"title" gorm:"->"`
	Content string `json:"content" gorm:"->"`
	Read    bool   `json:"read"`
	//...
}

type Messages interface {
	GetUnread() ([]Msg, error)
	GetIndex(id int) (*Msg, error)
	GetAll() ([]Msg, error)
	MarkRead(idx int) error

	//...
}
