package config

// yaml
// log:
// both: all #all ,file console,console is dufault
// level: info # info ,error ...
// format: json #json,text default text
// mode: date #size,date ,default size
// file :
// path: "./" # file path
// size : 60 #M
// max_age: 80  # 日志文件存储最大天数

type Config struct {
	Type   string // log type std zap logrus,zap is default
	Both   string // all ,file console,console is default
	Level  string // info ,error ...
	Format string // json text
	File   FileConf
}
type FileConf struct {
	Mode   string // size data
	Path   string // file path
	MaxAge int    // file maxAge
	Size   int    // file size （M）
}
