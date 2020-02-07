package nori

type Config struct {
	Nori struct {
		Storage string
		Hooks   []string
	}
	Plugins struct {
		Dir []interface{}
	}
}
