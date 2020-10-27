package batalhanaval

import "fmt"

//TamanhoTabuleiro tamanho do tabuleiro
const TamanhoTabuleiro = 10

//Tabuleiro struct para definir as operações de um tabuleiro
type Tabuleiro struct {
	tabuleiro [][]byte
}

//NovoTabDefesa construtor de um tabuleiro de defesa
func NovoTabDefesa(tabuleiro [][]byte) *Tabuleiro {
	tab := &Tabuleiro{}
	tab.tabuleiro = tabuleiro
	return tab
}

//NovoTabAtaque construtor de um tabuleiro de ataque
func NovoTabAtaque() *Tabuleiro {
	tab := &Tabuleiro{}
	tab.tabuleiro = make([][]byte, TamanhoTabuleiro)
	for i := 0; i < TamanhoTabuleiro; i++ {
		tab.tabuleiro[i] = make([]byte, TamanhoTabuleiro)
		for j := 0; j < TamanhoTabuleiro; j++ {
			tab.tabuleiro[i][j] = '-'
		}
	}
	return tab
}

//Imprimir função que imprime o tabuleiro corrente
func (t *Tabuleiro) Imprimir() {
	//Primeira linha -> índice
	fmt.Print("  ")
	for i := 1; i <= TamanhoTabuleiro; i++ {
		fmt.Print(" ", i)
	}
	fmt.Println()
	//Linhas seguintes
	for i := 0; i < TamanhoTabuleiro; i++ {
		fmt.Printf(" %c", ('A' + i))
		for j := 0; j < TamanhoTabuleiro; j++ {
			fmt.Printf(" %c", t.tabuleiro[i][j])
		}
		fmt.Println()
	}
}

//ReceberTiro função que do tabuleiro de uleiro
//que recebe um tiro e verifica se acertou e se
//afundou um navio
func (t *Tabuleiro) ReceberTiro(i, j int) bool {
	acertou := t.tabuleiro[i][j] == 'N'
	t.RegistrarTiro(i, j, acertou)
	return acertou
}

//RegistrarTiro função que registra um tiro realizado já com o resultado se acertou ou não
func (t *Tabuleiro) RegistrarTiro(i, j int, acertou bool) {
	if acertou {
		t.tabuleiro[i][j] = 'V'
	} else {
		t.tabuleiro[i][j] = 'X'
	}
}

//AfundouTodos função que verifica se afundou todos os navios
//Obs.: esta função só é implementada deste jeito porque
//o tamanho do tabuleiro é fixo e pequeno
func (t *Tabuleiro) AfundouTodos() bool {
	for i := 0; i < TamanhoTabuleiro; i++ {
		for j := 0; j < TamanhoTabuleiro; j++ {
			if t.tabuleiro[i][j] == 'N' {
				return false
			}
		}
	}
	return true
}
