# Formato para que sea mas legible el codigo 
1. Variables en camelCase (si la variable es global PascalCase)
2. Funciones separadas con _ 


# Preguntas para los soportes
1. Cuando llega una systemcall se hacen dos conexiones (memoria,kernel) o solo memoria y ella le avisa a kernel? 
2. "Al llegar un nuevo proceso a esta cola y la misma esté vacía se enviará un pedido a Memoria para inicializar el mismo, si la respuesta es positiva se crea el TID 0 de ese proceso y se lo pasa al estado READY y se sigue la misma lógica con el proceso que sigue"
    Se pasa el PCB a la cola de ready o el hilo?
