package main

import (
	"socketsApplication/batalhaNaval/batalhanaval"
)

func main() {
	var tabuleiro [10][10]byte
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			tabuleiro[i][j] = '-'
		}
	}

	tab := batalhanaval.NovoTabDefesa(tabuleiro)
	tab.ReceberTiro(5, 4)
	tab.Imprimir()
}
