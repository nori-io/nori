package nori

type Config struct {
	Nori struct {
		Storage string
	}
	Plugins struct {
		Dir []interface{}
	}
}
