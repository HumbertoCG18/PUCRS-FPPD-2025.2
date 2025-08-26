# Análise de Código Go: Operações em Árvore Binária (Sequencial vs. Concorrente)

Este código Go é um prato cheio para quem quer entender operações em árvores binárias. Ele não só faz o básico, como percorrer, somar e buscar valores, mas coloca lado a lado duas filosofias: a abordagem sequencial, mais tradicional, e a concorrente, usando goroutines e canais — o superpoder do Go.

## A Estrutura da Árvore: `Nodo`

A peça fundamental da nossa árvore é a struct `Nodo`. É bem simples: cada nó tem um valor e dois ponteiros, um para o filho da esquerda e outro para o da direita.

```go
type Nodo struct {
    v int
    e *Nodo // Ponteiro para o nó da esquerda
    d *Nodo // Ponteiro para o nó da direita
}
```

## Funções e Operações

Vamos dar uma olhada nas principais funções e o que cada uma faz.

### 1. `caminhaERD` - Percorrendo a Árvore

Uma função clássica para passear pela árvore. A sigla `ERD` entrega o caminho: **Esquerda -> Raiz -> Direita**. A função mergulha recursivamente pela esquerda, imprime o valor do nó em que está e depois explora a direita. Na prática, isso faz com que os valores sejam impressos em ordem crescente, como mágica.

```go
func caminhaERD(r *Nodo) {
    if r != nil {
        caminhaERD(r.e)
        fmt.Print(r.v, ", ") // Imprime o valor do nó atual
        caminhaERD(r.d)
    }
}
```

### 2. `soma` vs. `somaConc` - Somando os Valores

Para somar os valores, temos duas versões que mostram bem a diferença de pensamento.

#### `soma` (Sequencial)

A função `soma` é bem direta. Se o nó não for nulo, ela retorna o valor do nó atual somado ao resultado da chamada recursiva para a esquerda e para a direita. Simples e eficaz, mas executa tudo em uma única thread.

```go
func soma(r *Nodo) int {
    if r != nil {
        return r.v + soma(r.e) + soma(r.d)
    }
    return 0
}
```

#### `somaConc` (Concorrente)

Já a versão concorrente, `somaConc`, bota o processador para trabalhar. Em vez de fazer uma conta de cada vez, ela "quebra" a tarefa: lança duas goroutines (pense nelas como threads leves) para calcular a soma da esquerda e da direita ao mesmo tempo. Cada lado faz seu trabalho e envia o resultado por um canal. A função original só precisa esperar esses dois valores chegarem para somar tudo. Em árvores gigantes e máquinas com vários núcleos, a diferença de velocidade pode ser brutal.

```go
// Funcao "wraper" que inicia a concorrência
func somaConc(r *Nodo) int {
    s := make(chan int)
    go somaConcCh(r, s) // Dispara a função concorrente
    return <-s           // Espera e retorna o resultado final
}

func somaConcCh(r *Nodo, s chan int) {
    if r != nil {
        s1 := make(chan int)
        go somaConcCh(r.e, s1) // Soma a esquerda em paralelo
        go somaConcCh(r.d, s1) // Soma a direita em paralelo
        // Espera os dois resultados, soma com o valor atual e envia
        s <- (r.v + <-s1 + <-s1)
    } else {
        s <- 0 // Se o nó é nulo, a soma é zero
    }
}
```

### 3. `busca` vs. `buscaC` - Procurando um Valor

A mesma briga between sequencial e concorrente aparece na hora de buscar um valor.

#### `busca` (Sequencial)

A busca sequencial verifica o nó atual. Se não for o valor, ela chama a si mesma para a sub-árvore esquerda. Se não encontrar, tenta a sub-árvore direita. O operador `||` (OU) garante que a busca pare assim que o valor for encontrado.

```go
func busca(r *Nodo, val int) bool {
    if r == nil {
        return false // Chegou ao fim de um galho e não achou
    }
    if r.v == val {
        return true // Achou!
    }
    // Procura na esquerda OU na direita
    return busca(r.e, val) || busca(r.d, val)
}
```

#### `buscaC` (Concorrente)

A busca concorrente é mais esperta. Ela também manda "espiões" (goroutines) para os dois lados da árvore ao mesmo tempo. O truque aqui é o `select`, que funciona como um vigia. Ele fica de olho nos canais de resposta dos dois lados e, no momento em que o primeiro "espião" grita "Achei!", ele encerra a operação e avisa que o valor foi encontrado, sem nem precisar esperar o outro terminar a busca. É eficiência pura.

```go
func buscaC(r *Nodo, val int) bool {
    resultado := make(chan bool, 1)
    go buscaConc(r, val, resultado)
    return <-resultado
}

func buscaConc(r *Nodo, val int, ret chan bool) {
    // ... (verificações iniciais se r é nulo ou se o valor foi encontrado)

    retE := make(chan bool, 1)
    retD := make(chan bool, 1)

    go buscaConc(r.e, val, retE) // Busca na esquerda
    go buscaConc(r.d, val, retD) // Busca na direita

    // Espera pelo primeiro resultado positivo que chegar
    for i := 0; i < 2; i++ {
        select {
        case resE := <-retE:
            if resE {
                ret <- true
                return
            }
        case resD := <-retD:
            if resD {
                ret <- true
                return
            }
        }
    }
    ret <- false // Se ninguém achou, retorna false
}
```

### 4. `retornaParImpar` - Separando Valores

Essa função é diferente. Em vez de devolver um resultado único, ela funciona como uma central de triagem. Conforme percorre a árvore, ela analisa cada número: se for par, joga no canal `saidaP`; se for ímpar, vai para o `saidaI`. É uma ótima forma de processar e categorizar dados em paralelo.

```go
func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
    if r != nil {
        if r.v%2 == 0 {
            saidaP <- r.v // Envia o valor para o canal de pares
        } else {
            saidaI <- r.v // Envia o valor para o canal de ímpares
        }
        retornaParImpar(r.e, saidaP, saidaI, fin)
        retornaParImpar(r.d, saidaP, saidaI, fin)
    } else {
        fin <- struct{}{} // Sinaliza que chegou ao fim de um galho
    }
}
```

### 5. `main` - Juntando Tudo

É na função `main` que o show acontece. Ela primeiro monta a árvore na mão. Depois, bota a `retornaParImpar` para rodar em segundo plano e fica "pescando" os números pares e ímpares que chegam nos canais, imprimindo-os na tela assim que aparecem. Por fim, ela chama todas as outras funções de soma e busca para a gente poder ver a diferença na prática. Aquele loop com `select` é um exemplo perfeito do poder do Go para lidar com várias coisas acontecendo ao mesmo tempo.

```go
func main() {
    // ... (criação da árvore)

    saidaP := make(chan int)
    saidaI := make(chan int)
    fin := make(chan struct{})

    go retornaParImpar(root, saidaP, saidaI, fin)

    fim := false
    // Loop para consumir os dados dos canais de pares e ímpares
    for count := 0; count < 20 && !fim; {
        select {
        case par := <-saidaP:
            fmt.Println("Par:", par)
            count++
        case impar := <-saidaI:
            fmt.Println("Impar:", impar)
            count++
        case <-fin:
            // Lógica para saber quando parar de verdade
            // (a lógica atual com 'fim' é simplificada)
            fim = true
        }
    }

    // ... (chamadas para as outras funções)
}
```