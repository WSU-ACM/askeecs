package main

func main() {
	s := NewServer()
	s.Init(".secret")
	s.Serve()
}
