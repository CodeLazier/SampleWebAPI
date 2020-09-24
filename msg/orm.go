/*
 * @Author: your name
 * @Date: 2020-09-24 10:35:54
 * @LastEditTime: 2020-09-24 15:39:24
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\msg\orm.go
 */
package msg

type OrmMsg interface {
	OpenOrm(...string) error
	Select(interface{}, ...interface{}) (interface{}, error)
	Update(idx int, field string, value interface{}) error
}
