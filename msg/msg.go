/*
 * @Author: your name
 * @Date: 2020-09-29 15:26:41
 * @LastEditTime: 2020-09-29 15:30:34
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\msg\msg.go
 */
package msg

type Messages interface {
	//get msgs for read,sort by senddate
	GetUnread(start int, count int) (interface{}, error)
	//get msg for uniqueid
	GetIndex(id int) (interface{}, error)
	//get all msgs,sort by senddate
	GetAll(start int, count int) (interface{}, error)
	//set read
	MarkRead(idx int) error

	GetUnreadCount() (int64, error)
	GetCount() (int64, error)

	//...
}
