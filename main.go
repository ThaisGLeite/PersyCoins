package main

import (
	"log"
)

// Checando por erros no codigo e loga o erro, usar ao inves de iferr
func Check(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {

}
