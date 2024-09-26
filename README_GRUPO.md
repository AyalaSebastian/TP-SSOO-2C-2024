# Formato para que sea mas legible el codigo 
1. Variables en camelCase (si la variable es global PascalCase)
2. Funciones separadas con _ 


# Preguntas para los soportes
1. Cuando llega una systemcall se hacen dos conexiones (memoria,kernel) o solo memoria y ella le avisa a kernel? 
2. "Al llegar un nuevo proceso a esta cola y la misma esté vacía se enviará un pedido a Memoria para inicializar el mismo, si la respuesta es positiva se crea el TID 0 de ese proceso y se lo pasa al estado READY y se sigue la misma lógica con el proceso que sigue"
    Si se pasa el hilo, que hacemos con el PCB?
3. Como es el tema de las colas? Cuando estan en New va el PCB y en Ready va el TCB? O son distintas colas para procesos e hilos? no hace mucho sentido ya que hilos de kernel se toman como procesos
4. "PROCESS_EXIT, esta syscall finalizará el PCB correspondiente al TCB que ejecutó la instrucción, enviando todos sus TCBs asociados a la cola de EXIT. Esta instrucción sólo será llamada por el TID 0 del proceso y le deberá indicar a la memoria la finalización de dicho proceso."
    Que se refiere con "finalizará"? eliminar? mover a exit?
    

# Orden de levantamiento de los módulos
1. Filesystem
2. Memoria
3. CPU
4. Kernel