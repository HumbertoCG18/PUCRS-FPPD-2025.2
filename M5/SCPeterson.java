/* 
   Algoritmo de exclusao mutua por SW para dois processos,
   de Peterson, em Java, exemplificado com contador compartilhado. 
   Disciplina:  MCC ou SO
   PUCRS - Escola Politecnica
   Prof: Fernando Dotti

   Note o uso do modificador volatile.
   Retire as protecoes da sc e veja o resultado.
   Retire volatile das definicoes de Peterson e veja o resultado.
*/

//  Leia o código, veja como threads sao criadas em Java.
//  Note que CounterThread herda de Thread - esta sendo parte de Java.
//  Note a passagem de referencias no construtor de CounterThread, 
//  depois utilizadas no corpo de run() ,  que define o que a thread executa
//  quando o metodo start() e invocado.

// 1) Note o uso do modificador volatile na variavel last.
//    Retire o modificador e teste.    Pesquise e responda:
//    Qual a função do modificador volatile ?


// ============= Peterson ==============
// algoritmo de exclusao mutua para dois processos, de Peterson.
// metodos implementam protocolos de entrada e saida da SC
class Peterson {
	private volatile boolean[] flag = new boolean[2];
	private volatile int last; // retire volatile daqui e veja o resultado... busque a explicacao
	
	public Peterson(){
  		flag[0]=false;  
  		flag[1]=false;
  		last=1;	
	}	
		
    public void lock(int id) {
     int j = 1 - id;
     flag[id] = true;
     last = id;
     while (!(!flag[j] || last == j)) {};
    }

    public void unlock(int id) {
     flag[id] = false;
    }		
}

// =================================================
// aqui temos uma aplicacao de Peterson.
// um simples contador que pode ser incrementado por multiplas threads
// cada thread invoca incr
// incr usa os protocolos do algoritmo de Peterson para proteger a variavel compartilhada
// se voce retirar as protecoes vera que deixa de ficar consistente
class CounterSC {
	private  int n;     
	private Peterson sc;
		
	public CounterSC(){ 	
		n = 0;
	    sc = new Peterson();
	} 
		
	public void incr(int id){
		sc.lock(id);   
        n++;   
		sc.unlock(id);  
	} 	
	
	public int value() { return n; }
}

// thread que usa o contador compartilhado.
// faz um nro de incrementos cfe sua instanciacao
class CounterThread extends Thread {

	private int id;
    private CounterSC c_sc;
	private int limit;

    public CounterThread(int _id, CounterSC _c_sc, int _limit){
		id = _id;	
    	c_sc = _c_sc;
		limit = _limit;
    }

    public void run() {
       for (int i = 0; i < limit; i++) {
		     c_sc.incr(id);    
      }
    }
}

// teste de peterson cria varias threads que vao acessar o mesmo contador,
// que usa o alg de peterson para proteger a 
class TestePeterson {
	public static void main(String[] args) {
	
	  int nrIncr = 100000;	

	  CounterSC c = new CounterSC();
      CounterThread p = new CounterThread(0,c,nrIncr);
      CounterThread q = new CounterThread(1,c,nrIncr);
	  
      p.start();
      q.start();
      try { p.join(); q.join(); }
	  
      catch (InterruptedException e) { }
      System.out.println("The value of n is " + c.value());
    }
	
}
