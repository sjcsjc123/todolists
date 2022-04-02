package errorCode

const (
	MysqlError     = 50001 //mysql数据库错误
	ParseJsonError = 50002 //json解析错误
	JwtToken       = 50003 //jwtToken错误
	RedisError     = 50004 //redis error
	InvalidToken   = 50005
	UsernameRepeat = 50006
	PhoneRepeat    = 50007
	PhoneError     = 50008
	PasswordError  = 50009
)
