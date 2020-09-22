/*
 * @Author: your name
 * @Date: 2020-09-22 11:18:00
 * @LastEditTime: 2020-09-22 14:31:43
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\db.go
 */
package msg

type Msg struct {
	Id      int    `json:"id"`
	Title   string `json:title`
	Content string `json:content`
	//...
}

type Messages interface {
	GetUnread() ([]Msg, error)
	GetIndex(id int) (*Msg, error)
	//...
}
