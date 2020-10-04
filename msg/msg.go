package msg

type Messages interface {
	//get msgs for read,sort by senddate
	GetUnread(int, int) (interface{}, error)
	//get msg for uniqueid
	GetIndex(int) (interface{}, error)
	//get all msgs,sort by senddate
	GetAll(int, int) (interface{}, error)
	//set read
	MarkRead(int) error
	//insert new
	New(interface{}) error

	GetUnreadCount() (int64, error)
	GetCount() (int64, error)

	//...
}
