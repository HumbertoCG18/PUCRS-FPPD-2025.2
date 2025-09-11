import threading
import time
import random

NUM_FILO = 5
ESQUERDA = lambda i: i
DIREITA = lambda i: (i + 1) % NUM_FILO

# Cada garfo é um semáforo
garfos = [threading.Semaphore(1) for _ in range(NUM_FILO)]


def pensar(i):
    tempo = random.randint(1, 3)
    print(f"Filósofo {i} está pensando por {tempo}s.")
    time.sleep(tempo)


def comer(i):
    tempo = random.randint(1, 3)
    print(f"Filósofo {i} está comendo por {tempo}s.")
    time.sleep(tempo)


def filosofo(i):
    while True:
        pensar(i)

        # Último filósofo pega na ordem inversa para evitar deadlock
        if i == NUM_FILO - 1:
            primeiro, segundo = DIREITA(i), ESQUERDA(i)
        else:
            primeiro, segundo = ESQUERDA(i), DIREITA(i)

        # Pega garfos
        garfos[primeiro].acquire()
        garfos[segundo].acquire()

        print(f"Filósofo {i} pegou os garfos {primeiro} e {segundo}.")
        comer(i)

        # Devolve garfos
        garfos[primeiro].release()
        garfos[segundo].release()
        print(f"Filósofo {i} largou os garfos {primeiro} e {segundo}.\n")


if __name__ == "__main__":
    filosofos = [threading.Thread(target=filosofo, args=(i,)) for i in range(NUM_FILO)]

    for f in filosofos:
        f.start()

    for f in filosofos:
        f.join()

