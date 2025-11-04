# Control II: Ejecucion Especulativa - Lenguajes de Programacion II

## Integrantes
- Lorena Uribe
- Andres Gonzalez

## Profesor
Alonso Inostrosa Psijas

## Repositorio
https://github.com/Uribe22/Lenguaje-de-Programaacion-.git

## Descripcion del Proyecto

Este proyecto implementa el patron de ejecucion especulativa en el lenguaje de programacion Go, utilizando Goroutines y canales para manejar la concurrencia y sincronizacion. El objetivo es simular la ejecucion en paralelo de dos ramas de computo intensivo mientras se evalua una condicion costosa en terminos de tiempo. El programa selecciona la rama correcta segun el resultado de la evaluacion y cancela la rama incorrecta, optimizando el uso de recursos.

## Requisitos de Implementacion

1. Simulacion de Tareas:
   - El programa ejecuta dos funciones de computo intensivo: SimularProofOfWork (para la rama A) y EncontrarPrimos (para la rama B).
   - La evaluacion de la condicion de seleccion de la rama correcta se basa en la funcion CalcularTrazaDeProductoDeMatrices.

2. Logica de Sincronizacion:
   - Se utilizan canales para la comunicacion entre las goroutines y el hilo principal.
   - La rama incorrecta se cancela una vez se selecciona la rama correcta.

3. Parametros de Entrada:
   - n: Dimension de las matrices para la multiplicacion.
   - umbral: Valor utilizado para decidir cual de las dos ramas se ejecutara.
   - nombre_archivo: Nombre del archivo donde se registran las metricas de ejecucion.

4. Salida del Programa:
   - Los tiempos de inicio y fin de cada computo, el resultado de la rama seleccionada y el tiempo total de ejecucion.

## Instrucciones de Compilacion y Ejecucion

### Compilacion
```bash
go build -o control2.exe main.go
```

### Ejecucion
Para ejecutar el programa, use el siguiente comando con los parametros adecuados:

```bash
go run main.go <blockData> <dificultad> <maxPrimos> <n> <umbral> <nombreArchivo> [--bmmode]
```

Donde:
- blockData: Datos del bloque para Proof-of-Work
- dificultad: Numero de ceros iniciales requeridos en el hash
- maxPrimos: Limite superior para busqueda de numeros primos
- n: Dimension de las matrices para calcular la traza (funcion de decision)
- umbral: Valor usado para decidir que rama ejecutar
- nombreArchivo: Archivo donde se registran las metricas
- --bmmode: (Opcional) Activa modo benchmark con 30 ejecuciones

### Ejemplos de Uso

Modo especulativo (ejecucion unica):
```bash
go run main.go blockchain_data 4 5000 15 200 salida.txt
```

Modo benchmark (30 ejecuciones):
```bash
go run main.go bloque_genesis 5 5000 15 200 benchmark.txt --bmmode
```

## Analisis de Rendimiento

### Metodologia
Se ejecutaron 30 simulaciones especulativas y 30 simulaciones secuenciales con parametros fijos para cada conjunto de pruebas.

### Mediciones Realizadas

#### Test 1: Parametros balanceados (Rama A)
Parametros:
- Block Data: blockchain_data
- Dificultad PoW: 4
- Max Primos: 5000
- Dimension Matriz: 15
- Umbral: 200

| Estrategia | Tiempo Promedio (ms) |
|------------|---------------------|
| Especulativo | 17 |
| Secuencial | 17 |

**Speedup = 17 / 17 = 1.00**

#### Test 2: Carga pesada (Rama A)
Parametros:
- Block Data: bloque_genesis
- Dificultad PoW: 5
- Max Primos: 5000
- Dimension Matriz: 15
- Umbral: 200

| Estrategia | Tiempo Promedio (ms) |
|------------|---------------------|
| Especulativo | 1410 |
| Secuencial | 1431 |

**Speedup = 1431 / 1410 = 1.01**

#### Test 3: Rama B (Primos)
Parametros:
- Block Data: datos_prueba
- Dificultad PoW: 3
- Max Primos: 100000
- Dimension Matriz: 15
- Umbral: 5000

| Estrategia | Tiempo Promedio (ms) |
|------------|---------------------|
| Especulativo | 4 |
| Secuencial | 3 |

**Speedup = 3 / 4 = 0.75**

### Conclusiones

1. En tareas de carga moderada a alta (Test 1 y 2), la ejecucion especulativa muestra overhead concurrente despreciable o incluso ventaja de rendimiento (speedup 1.00-1.01).

2. El Test 2 demuestra que con tareas pesadas (mas de 1 segundo), el modo especulativo puede superar ligeramente al secuencial (speedup 1.01), ya que mientras se calcula la funcion de decision, ambas ramas ya estan ejecutandose.

3. En tareas muy rapidas (menos de 5ms, Test 3), el overhead de crear goroutines y sincronizar canales se vuelve significativo (speedup 0.75), resultando en una penalizacion de rendimiento del 25%.

4. La ejecucion especulativa es beneficiosa cuando el tiempo de la funcion de decision es comparable o menor al tiempo de ejecucion de las ramas, permitiendo paralelizar el computo.

5. El patron implementado cumple correctamente con la cancelacion de la rama perdedora mediante canales, evitando desperdicio de recursos computacionales.

## Estructura del Proyecto

El proyecto contiene los siguientes archivos:
- main.go: Codigo principal que implementa la logica de ejecucion especulativa y secuencial
- README.md: Documentacion sobre la ejecucion, parametros de entrada y analisis de rendimiento
- Archivos de salida: Archivos generados durante la ejecucion que contienen los tiempos de ejecucion y resultados
