package web

import (
	"token/web/controller"
	"token/web/route"
)

func Start() {
	controller.Init()
	route.Init()
}
