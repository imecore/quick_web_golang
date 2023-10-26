package provider

var (
	Database       *Mysql
	Cache          *Redis
	SessionManager *Session
)

func Init() {
	Database = (&Mysql{}).New()
	Cache = (&Redis{}).New()
	SessionManager = (&Session{}).New()
}

type Provider interface {
	Start()
	Close()
}
