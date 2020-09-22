/*
 * @Author: your name
 * @Date: 2020-09-22 11:18:00
 * @LastEditTime: 2020-09-22 12:30:36
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\db.go
 */
package msg

type Message struct {
	Id      int    `json:"id"`
	Title   string `json:title`
	Content string `json:content`
	//...
}

type Messages interface {
	GetUnread() ([]Message, error)
	GetIndex(id int) (*Message, error)
	//...
}
