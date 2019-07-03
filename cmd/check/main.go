package main

import (
	"github.com/terminatingcode/exhibit-resource/resource"
	"github.com/cloudboss/ofcourse/ofcourse"
)

func main() {
	ofcourse.Check(&resource.Resource{})
}
