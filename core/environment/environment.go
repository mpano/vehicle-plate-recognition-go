package environment

type Environment struct {
	Port  string
	DBURL string
}

func New(
	port string,
	DBURL string,

) *Environment {
	return &Environment{
		Port:  port,
		DBURL: DBURL,
	}
}
