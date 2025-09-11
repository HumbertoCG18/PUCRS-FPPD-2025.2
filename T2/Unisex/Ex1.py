import threading
import time
import random

class Lightswitch:
    def __init__(self):
        self.counter = 0
        self.mutex = threading.Lock()

    def lock(self, semaphore):
        self.mutex.acquire()
        self.counter += 1
        if self.counter == 1:
            semaphore.acquire()
        self.mutex.release()

    def unlock(self, semaphore):
        self.mutex.acquire()
        self.counter -= 1
        if self.counter == 0:
            semaphore.release()
        self.mutex.release()


empty = threading.Semaphore(1)         # garante exclusão entre homens e mulheres
maleSwitch = Lightswitch()
femaleSwitch = Lightswitch()
maleMultiplex = threading.Semaphore(3)   # máximo de 3 homens
femaleMultiplex = threading.Semaphore(3) # máximo de 3 mulheres

def man(id):
    while True:
        time.sleep(random.uniform(0.5, 2))  # tempo fora do banheiro
        print(f"👨 Homem {id} quer entrar no banheiro.")

        maleSwitch.lock(empty)              # primeiro homem bloqueia mulheres
        maleMultiplex.acquire()             # só até 3 homens
        print(f"👨 Homem {id} entrou no banheiro.")
        time.sleep(random.uniform(1, 3))    # usando banheiro
        print(f"👨 Homem {id} saiu do banheiro.")
        maleMultiplex.release()
        maleSwitch.unlock(empty)            # último homem libera mulheres


def woman(id):
    while True:
        time.sleep(random.uniform(0.5, 2))  # tempo fora do banheiro
        print(f"👩 Mulher {id} quer entrar no banheiro.")

        femaleSwitch.lock(empty)            # primeira mulher bloqueia homens
        femaleMultiplex.acquire()           # só até 3 mulheres
        print(f"👩 Mulher {id} entrou no banheiro.")
        time.sleep(random.uniform(1, 3))    # usando banheiro
        print(f"👩 Mulher {id} saiu do banheiro.")
        femaleMultiplex.release()
        femaleSwitch.unlock(empty)          # última mulher libera homens


if __name__ == "__main__":
    num_homens = 5
    num_mulheres = 5

    for i in range(num_homens):
        threading.Thread(target=man, args=(i,), daemon=True).start()

    for i in range(num_mulheres):
        threading.Thread(target=woman, args=(i,), daemon=True).start()

    # deixa as threads rodando
    while True:
        time.sleep(1)
