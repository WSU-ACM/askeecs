package main

func main() {
	s := NewServer()

	//Load cookie jar secret from file
	s.Init(".secret")
	s.Serve()
}
