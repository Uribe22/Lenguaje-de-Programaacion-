# Control II: Ejecucion Especulativa

## Integrantes
- Lorena Uribe
- Andres Gonzalez

## Profesor
Alonso Inostrosa Psijas

## Repositorio
https://github.com/Uribe22/Lenguaje-de-Programaacion-.git

## Descripcion del Proyecto

Este proyecto implementa el patron de ejecucion especulativa utilizando el lenguaje Go. Se aprovechan las goroutines y canales para manejar la concurrencia entre dos ramas de computo que se ejecutan en paralelo: una rama realiza Proof of Work y la otra busca numeros primos.

Mientras ambas ramas se ejecutan, se calcula una funcion de decision basada en la multiplicacion de matrices y el calculo de su traza. Una vez obtenido el resultado de esta funcion, se selecciona la rama correspondiente y se cancela la otra mediante el cierre de un canal de cancelacion.

## Compilacion

Para compilar el programa:

```bash
go build -o control2.exe main.go
```

## Ejecucion

El programa recibe los siguientes parametros:

```bash
go run main.go <blockData> <dificultad> <maxPrimos> <n> <umbral> <nombreArchivo> [--bmmode]
```

**Descripcion de parametros:**
- `blockData`: Datos del bloque para Proof-of-Work
- `dificultad`: Numero de ceros iniciales requeridos en el hash
- `maxPrimos`: Limite superior para la busqueda de numeros primos
- `n`: Dimension de las matrices para calcular la traza
- `umbral`: Valor que determina que rama ejecutar (si traza > umbral → Rama A, sino → Rama B)
- `nombreArchivo`: Archivo donde se registran los resultados
- `--bmmode`: (Opcional) Activa el modo benchmark con 30 ejecuciones

**Ejemplos de uso:**

Ejecucion normal:
```bash
go run main.go blockchain_data 4 5000 15 200 salida.txt
```

Modo benchmark:
```bash
go run main.go bloque_genesis 5 5000 15 200 benchmark.txt --bmmode
```

## Analisis de Rendimiento

Se realizaron tres pruebas con diferentes configuraciones de parametros. Cada prueba ejecuta 30 iteraciones en modo especulativo y 30 en modo secuencial para calcular tiempos promedio y el speedup resultante.

### Prueba 1: Carga Moderada
```bash
go run main.go blockchain_data 4 5000 15 200 salida1.txt --bmmode
```

**Resultados:**
- Tiempo Secuencial: 17ms
- Tiempo Especulativo: 17ms
- Speedup: 1.00

### Prueba 2: Carga Alta
```bash
go run main.go bloque_genesis 5 5000 15 200 salida2.txt --bmmode
```

**Resultados:**
- Tiempo Secuencial: 1431ms
- Tiempo Especulativo: 1410ms
- Speedup: 1.01

### Prueba 3: Carga Baja (Rama B)
```bash
go run main.go datos_prueba 3 100000 15 5000 salida3.txt --bmmode
```

**Resultados:**
- Tiempo Secuencial: 3ms
- Tiempo Especulativo: 4ms
- Speedup: 0.75

## Conclusiones

Los resultados obtenidos demuestran que la ejecucion especulativa no siempre ofrece mejoras de rendimiento. El beneficio depende de la relacion entre el tiempo de computo de las ramas y el tiempo de la funcion de decision.

En la Prueba 2, donde el tiempo de ejecucion es considerable (>1s), se observa una ligera mejora con speedup de 1.01. Esto se debe a que mientras se calcula la funcion de decision, las ramas ya estan ejecutandose en paralelo.

En contraste, la Prueba 3 muestra un speedup de 0.75, indicando que el modo secuencial fue mas eficiente. Esto ocurre cuando las tareas son muy rapidas (<5ms) y el overhead de crear goroutines y sincronizar canales supera el beneficio de la paralelizacion.

## Estructura del Codigo

El archivo `main.go` contiene las siguientes funciones principales:

- **SimularProofOfWork**: Calcula un nonce que cumple con la dificultad especificada usando SHA-256
- **EncontrarPrimos**: Busca numeros primos hasta el valor maximo indicado
- **CalcularTrazaDeProductoDeMatrices**: Calcula la traza del producto de dos matrices (funcion de decision)
- **EjecutarEspeculativo**: Ejecuta ambas ramas en paralelo y cancela la rama no seleccionada
- **EjecutarSecuencial**: Ejecuta unicamente la rama correcta segun el resultado de la decision
- **GuardarBenchmark**: Realiza 30 ejecuciones de cada modo y calcula estadisticas de rendimiento
