package happy

import "fmt"

func main() {
	RouteAny("/",Controller)
	Run()
}

func Controller(this *Context)  {
	fmt.Printf("进入Controller")
}