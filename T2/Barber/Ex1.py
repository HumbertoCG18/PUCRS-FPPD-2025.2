import threading
import random
import time

mtx = threading.Lock()
cv = threading.Condition(mtx)
waiting = []
NUM_CHAIRS = 3

def barber():
    while True:
        with cv:
            while not waiting:
                cv.wait()
            cust = waiting.pop(0)
        print(f'Barber cutting hair of customer {cust}')
        time.sleep(random.randint(1, 3))
        print(f'Barber finished with customer {cust}')
        with cv:
            cv.notify()

def customer(i):
    time.sleep(random.randint(1, 4))
    with cv:
        if len(waiting) < NUM_CHAIRS:
            waiting.append(i)
            print(f'Customer {i} waiting')
            cv.notify()
            while waiting[0] != i:
                cv.wait()
            print(f'Customer {i} got haircut')
        else:
            print(f'Customer {i} left (no seat)')

if __name__ == '__main__':
    t1 = threading.Thread(target=barber)
    t1.start()
    for i in range(5):
        threading.Thread(target=customer, args=(i,)).start()
    t1.join()