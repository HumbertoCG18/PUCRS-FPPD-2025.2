/*   1. Crie um diretório para seus fontes, por exemplo goprogs.
     2. Dentro desta pasta voce ira criar seus programas em Go.
     3. Mude para o diretório goprogs (cd gopros)
     4. Entre: "go mod init goprogs".    Este comando criará um arquivo go.mod em goprogs
     5. Entre: "go get github.com/nsf/termbox-go"
     6. Agora voce pode usar este pacote.  veja exemplo abaixo.
	 Instrucoes:
	 Use as setas para andar sobre o tabuleiro.
	 Digite uma letra para gravar a mesma por onde passa.
	 Digite ESC para sair.
*/

package main

import (
	"fmt"

	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync() // cosmestic purpose
}

type GameState struct {
	xSize int
	ySize int
	mesa  [][]rune
}
type Player struct {
	x  int
	y  int
	ch rune
}

func printState(g GameState) {
	reset()
	//fmt.Println(g)
	for i := 0; i < g.xSize; i++ {
		for j := 0; j < g.ySize; j++ {
			fmt.Print(string(g.mesa[i][j]))
		}
		fmt.Println("")
	}
}

func main() {
	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	fmt.Println("Enter any key to see their ASCII code or press ESC button to quit")

	const Xs = 25
	const Ys = 50

	mesa1 := make([][]rune, Xs)
	for i := 0; i < Xs; i++ {
		mesa1[i] = make([]rune, Ys)
	}
	for i := 0; i < Xs; i++ {
		for j := 0; j < Ys; j++ {
			mesa1[i][j] = '.'
		}
	}
	g := GameState{xSize: Xs, ySize: Ys, mesa: mesa1}
	p := Player{x: 0, y: 0, ch: ' '}
	pOld := Player{x: -1, y: -1, ch: ' '}

	printState(g)

keyPressListenerLoop:
	for {
		if !(p == pOld) {
			g.mesa[p.x][p.y] = p.ch
			printState(g)
			pOld = p
		}
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {

			// -CASOS NAO USADOS AQUI ----------------------------------------------
			case term.KeyF1:
				reset()
				fmt.Println("F1 pressed")
			case term.KeyF2:
				reset()
				fmt.Println("F2 pressed")
			case term.KeyF3:
				reset()
				fmt.Println("F3 pressed")
			case term.KeyF4:
				reset()
				fmt.Println("F4 pressed")
			case term.KeyF5:
				reset()
				fmt.Println("F5 pressed")
			case term.KeyF6:
				reset()
				fmt.Println("F6 pressed")
			case term.KeyF7:
				reset()
				fmt.Println("F7 pressed")
			case term.KeyF8:
				reset()
				fmt.Println("F8 pressed")
			case term.KeyF9:
				reset()
				fmt.Println("F9 pressed")
			case term.KeyF10:
				reset()
				fmt.Println("F10 pressed")
			case term.KeyF11:
				reset()
				fmt.Println("F11 pressed")
			case term.KeyF12:
				reset()
				fmt.Println("F12 pressed")
			case term.KeyInsert:
				reset()
				fmt.Println("Insert pressed")
			case term.KeyDelete:
				reset()
				fmt.Println("Delete pressed")
			case term.KeyHome:
				reset()
				fmt.Println("Home pressed")
			case term.KeyEnd:
				reset()
				fmt.Println("End pressed")
			case term.KeyPgup:
				reset()
				fmt.Println("Page Up pressed")
			case term.KeyPgdn:
				reset()
				fmt.Println("Page Down pressed")
			case term.KeySpace:
				reset()
				fmt.Println("Space pressed")
			case term.KeyBackspace:
				reset()
				fmt.Println("Backspace pressed")
			case term.KeyEnter:
				reset()
				fmt.Println("Enter pressed")
			case term.KeyTab:
				reset()
				fmt.Println("Tab pressed")
				// -ATE AQUI NAO USA AINDA ----------------------------------------------

				// -USAMOS DAQUI PARA BAIXO----------------------------------------------
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyArrowUp:
				// reset()
				fmt.Println("Arrow Up pressed")
				p.x = (p.x + Xs - 1) % Xs
			case term.KeyArrowDown:
				// reset()
				fmt.Println("Arrow Down pressed")
				p.x = (p.x + 1) % Xs
			case term.KeyArrowLeft:
				// reset()
				fmt.Println("Arrow Left pressed")
				p.y = (p.y + Ys - 1) % Ys
			case term.KeyArrowRight:
				// reset()
				fmt.Println("Arrow Right pressed")
				p.y = (p.y + 1) % Ys
			default:
				// we only want to read a single character or one key pressed event
				// reset()
				p.ch = ev.Ch
				fmt.Println("ASCII : ", ev.Ch)
				// --------------------------------------------------------------------

			}

		case term.EventError:
			panic(ev.Err)
		}
	}
}
