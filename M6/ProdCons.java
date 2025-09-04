/* 
   Exercício 2
   Exemplo de produtor consumidor com buffer finito.
   Disciplina:  Sistemas Operacionais
   PUCRS - Escola Politecnica
   Prof: Fernando Dotti
*/
import java.util.concurrent.Semaphore;

class FiniteBuffer {
	
  private int size;    
  private int in = 0;
  private int out = 0;
  private int[] buffer; 
  
  private Semaphore naoCheio;  
  private Semaphore naoVazio;  
  private Semaphore mutex;

  private void incrIn() {  in = (in+1)%size; }
  private void incrOut() { out = (out+1)%size; }

  public FiniteBuffer(int _size){ 	
    size = _size;
    buffer = new int[size];         // armazena os itens

    mutex = new Semaphore(1);	    // para exclusao mutua  (sc)
    naoCheio = new Semaphore(size); // controle de espaco disponivel
    naoVazio = new Semaphore(0);    // controle de itens
  } 
   
  public void insert(int v){
    try { naoCheio.acquire();  //  espera ter espaco  
          mutex.acquire();     //  entra sc
	    } catch (InterruptedException ie) {}  
        buffer[in]=v;          //  sc: insere elemento
        incrIn();              //  sc: insere elemento
      mutex.release();         //  sai sc 
      naoVazio.release();      //  avisa que nao esta vazio
  } 
  
  public int delete(){
    int val;
    try {  naoVazio.acquire(); //  espera nao estar vazio
           mutex.acquire();    //  entra na sc
		}  catch (InterruptedException ie) {}  
        val = buffer[out];     //  sc: retira elemento
        incrOut();             //  sc: retira elemento
      mutex.release();         //  sai da sc 
      naoCheio.release();      //  avisa que tem espaco
      return val;
  } 
}

class ProducerThread extends Thread { 
  private int id;
  private int limit;
  private FiniteBuffer fb;

    public ProducerThread(int _id, FiniteBuffer _fb, int _limit){
	   id = _id;	 fb = _fb;     limit = _limit;
    }

    public void run() {
       for (int i = 0; i < limit; i++) {
         fb.insert(i); 
         System.out.println("Prod "+id+"  val "+i);   
       }
    }
}

class ConsumerThread extends Thread {
  private int id;
  private int limit;
  private FiniteBuffer fb;

    public ConsumerThread(int _id, FiniteBuffer _fb, int _limit){
	   id = _id;	 fb = _fb;   limit = _limit;
    }

    public void run() {
       int v;
       for (int i = 0; i < limit; i++) {
         v = fb.delete();    
         System.out.println("                                  Cons "+id+"  val "+v);
       }
    }
}

class TesteProdCons {
	public static void main(String[] args) {
		  
	  FiniteBuffer fb = new FiniteBuffer(10);
	   
      ProducerThread p = new ProducerThread(1,fb,10);
      ProducerThread q = new ProducerThread(2,fb,10);
      ProducerThread r = new ProducerThread(3,fb,10);

      ConsumerThread s = new ConsumerThread(3,fb,15);
      ConsumerThread t = new ConsumerThread(4,fb,15);
	  
      p.start();  q.start();  r.start();  s.start();  t.start();
      try { p.join(); q.join(); r.join(); s.join(); t.join(); }
      catch (InterruptedException e) { }
      System.out.println("Fim ");
    }	
}
