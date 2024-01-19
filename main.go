package main

func main() {
	c := NewAppConfig()
	s := NewServer(c)
	s.Start()
}
