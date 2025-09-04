/* Copyright (C) 2006 M. Ben-Ari. See copyright.txt */
class Count extends Thread {
    static volatile int n = 0;

    public void run() {
      for (int i = 0; i < 10000; i++) {
				n++;  
      }
    }
    public static void main(String[] args) {
      Count p = new Count();
      Count q = new Count();
      p.start();
      q.start();
      try { p.join(); q.join(); }
      catch (InterruptedException e) { }
      System.out.println("The value of n is " + n);
    }
}
