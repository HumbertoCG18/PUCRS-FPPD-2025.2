import threading
import time
import random
from collections import deque

# Configuração
NUM_BARBEIROS = 3
NUM_SOFA = 4
NUM_STANDING = 16

# Semáforos
mutex = threading.Semaphore(1)
standingRoom = deque(maxlen = NUM_STANDING)
sofa = deque(maxlen = NUM_SOFA)
chair = threading.Semaphore(NUM_BARBEIROS)
barber_ready = threading.Semaphore(0)
customer_ready = threading.Semaphore(0)
cash = threading.Semaphore(0)
receipt = threading.Semaphore(0)

def barbeiro(id):
    while True:
        customer_ready.acquire()     # Espera cliente pronto
        chair.acquire()              # Ocupa uma cadeira de barbeiro
        print(f"Barbeiro {id} começou a atender um cliente.")
        barber_ready.release()       # Sinaliza que está pronto
        # Cortando cabelo
        time.sleep(random.uniform(1, 3))
        print(f"Barbeiro {id} terminou o corte.")
        cash.acquire()               # Espera cliente pagar
        print(f"Barbeiro {id} recebeu pagamento.")
        receipt.release()            # Dá recibo
        chair.release()              # Libera cadeira

def cliente(id):
    # Chegada aleatória
    time.sleep(random.uniform(0.5, 2))
    with mutex:  # Região crítica
        if len(standingRoom) == NUM_STANDING:
            print(f"Cliente {id} foi embora (sem espaço em pé).")
            return
        standingRoom.append(id)
        print(f"Cliente {id} esperando em pé ({len(standingRoom)}/{NUM_STANDING}).")
    
    # Espera vaga no sofá
    while True:
        with mutex:
            if len(sofa) < NUM_SOFA:
                sofa.append(id)
                standingRoom.popleft()  # Sai da fila em pé
                print(f"Cliente {id} sentou no sofá ({len(sofa)}/{NUM_SOFA}).")
                break
    
    # Espera barbeiro
    customer_ready.release()
    barber_ready.acquire()
    print(f"Cliente {id} cortando cabelo.")
    # Pagamento
    cash.release()
    receipt.acquire()
    print(f"Cliente {id} saiu após pagar.")

if __name__ == "__main__":
    # Threads barbeiros
    for b in range(NUM_BARBEIROS):
        threading.Thread(target = barbeiro, args = (b,), daemon=True).start()
    
    # Threads clientes
    for c in range(20):
        threading.Thread(target = cliente, args = (c,)).start()
    
    time.sleep(30)
