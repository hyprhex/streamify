package main

func main() {
	api := NewAPIServer(":3000", nil)
	api.Serve()
}
