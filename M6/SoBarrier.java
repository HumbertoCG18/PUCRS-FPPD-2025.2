/* 
   Implementacao de Barreira com Semaforo, em Java
   PUCRS - Escola Politecnica
   Prof: Fernando Dotti
*/
import java.util.concurrent.Semaphore;

// ============= a nossa Barreira ==============
// Esta ee uma barreira reutilizavel, usando semaforos.
// Baseado em "Little book of semaphores".
// Veja material passo a passo com explicacao do funcionamento

class Barrier {
	
	private Semaphore mutex;  
	private Semaphore catraca1;  
	private Semaphore catraca2;  
	private int count;         // nro que chegou
	private int n;             // nro de threads a chegar na barreira
	
	public Barrier(int nmax){ 	
		n = nmax;
		count = 0;
        mutex = new Semaphore(1);
		catraca1 = new Semaphore(0);
		catraca2 = new Semaphore(1);	
	} 
		
	public void arrive() throws InterruptedException {
		mutex.acquire();        // acesso exclusivo aa variaveis da barreira
		count++;
		if (count == n) {        // chegou ao nro de threads
			catraca2.acquire();   //  bloqueia a segunda catraca
			catraca1.release();   //  solta a primeira
		}
		mutex.release();       // solta bloqueio de variaveis
		catraca1.acquire();    // chega na catraca - aqui param todos antes de count == n
		catraca1.release();    // quando liberado, libera proximo da catraca, um a um
	}

	public void leave() throws InterruptedException {      // funcionamento parecido ao de acima, mas para sair do ponto critico
		mutex.acquire();
		count--;
		if (count == 0) {
			catraca1.acquire();
			catraca2.release();
		}
		mutex.release();
		catraca2.acquire();
		catraca2.release();
	}
}

