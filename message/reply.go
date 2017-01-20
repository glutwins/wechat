package message

type Reply interface {
	SetToUserName(string)
	SetFromUserName(string)
	SetCreateTime(int64)
}
