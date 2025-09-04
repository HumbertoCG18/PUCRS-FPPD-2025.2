/* 
   Exemplo de exclusao mutua com semaforo, em Java,
   exemplificado com contador compartilhado. 
   PUCRS - Escola Politecnica
   Prof: Fernando Dotti

   as operacoes        acquire e  release
   sao equivalentes a  wait    e  signal   da literatura
*/

// a classe abaixo especifica um contador onde cada incremento
// faz a protecao para exclusao mutua.   Ela usa um semaforo
// iniciado em 1 para isso.
class CounterSema {
	
	private  int n;   
	private Semaphore s;  
	
	public CounterSema(){ 	
		n = 0;
        s = new Semaphore(1);	
	} 
		
	public void incr(int id){
		try { 
			s.acquire();     //  WAIT
		} catch (InterruptedException ie) {}  
        
		n++;                 //   codigo da SC
		
		   s.release();      //   SIGNAL
	} 	
	
	public int value() { return n; }
}


// abaixo apenas implementamos threads que invocam o 
// incremento em um objeto CounterSema
class CounterThread extends Thread {

	private int id;
    private CounterSema c_s;
	private int limit;

    public CounterThread(int _id, CounterSema _c_s, int _limit){
		id = _id;	
    	c_s = _c_s;
		limit = _limit;
    }

    public void run() {
       for (int i = 0; i < limit; i++) {
		     c_s.incr(id);    
      }
    }
}


// o Teste cria um unico objeto CounterSema e passa para 5 threads
// que incrementam o mesmo nrIncr vezes (100000 como exemplo)
class TesteSemaphore {
	public static void main(String[] args) {
	
	  int nrIncr = 100000;	
	  
	  CounterSema c = new CounterSema();
	  
	  // criamos os objetos - ainda nao executam
      CounterThread p = new CounterThread(0,c,nrIncr);
      CounterThread q = new CounterThread(1,c,nrIncr);
      CounterThread r = new CounterThread(2,c,nrIncr);
      CounterThread s = new CounterThread(3,c,nrIncr);
      CounterThread t = new CounterThread(4,c,nrIncr);
	  
      // aqui dispara as threads!!  o start dispara o metodo run de cada thread
	  // como um processo concorrente aos demais
	  p.start();
      q.start();
      r.start();
      s.start();
      t.start();
	  // aqui temos main + 5 threads ativas
      try { p.join(); q.join(); r.join(); s.join(); t.join(); }
	  
	  // aqui temos somente main ativa
      catch (InterruptedException e) { }
      System.out.println("The value of n is " + c.value());
    }
	
}
