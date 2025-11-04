package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type ResultPOW struct {
	Hash  string
	Nonce int
}

func SimularProofOfWork(blockData string, dificultad int, canal chan ResultPOW, cancelar chan struct{}) {
	targetPrefix := strings.Repeat("0", dificultad)
	nonce := 0
	for {
		select {
		case <-cancelar:
			return
		default:
			data := fmt.Sprintf("%s%d", blockData, nonce)
			hashBytes := sha256.Sum256([]byte(data))
			hashString := fmt.Sprintf("%x", hashBytes)

			if strings.HasPrefix(hashString, targetPrefix) {
				canal <- ResultPOW{
					Hash:  hashString,
					Nonce: nonce,
				}
				return
			}
			nonce++
		}
	}
}

func EncontrarPrimos(max int, canal chan []int, cancelar chan struct{}) {
	var primes []int
	for i := 2; i < max; i++ {
		select {
		case <-cancelar:
			return
		default:
			isPrime := true
			for j := 2; j*j <= i; j++ {
				if i%j == 0 {
					isPrime = false
					break
				}
			}
			if isPrime {
				primes = append(primes, i)
			}
		}
	}
	canal <- primes
}
func EncontrarPrimosSecuencial(max int) []int {
	var primes []int
	for i := 2; i < max; i++ {
		isPrime := true
		for j := 2; j*j <= i; j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primes = append(primes, i)
		}
	}
	return primes
}

func SimularProofOfWorkSecuencial(blockData string, dificultad int) ResultPOW {
	targetPrefix := strings.Repeat("0", dificultad)
	nonce := 0

	for {
		data := fmt.Sprintf("%s%d", blockData, nonce)
		hashBytes := sha256.Sum256([]byte(data))
		hashString := fmt.Sprintf("%x", hashBytes)

		if strings.HasPrefix(hashString, targetPrefix) {
			// Retorna directamente el resultado sin usar canales
			return ResultPOW{
				Hash:  hashString,
				Nonce: nonce,
			}
		}
		nonce++
	}
}

func CalcularTrazaDeProductoDeMatrices(n int) int {
	// Se crean dos matrices con valores aleatorios para la multiplicación.
	m1 := make([][]int, n)
	m2 := make([][]int, n)
	for i := 0; i < n; i++ {
		m1[i] = make([]int, n)
		m2[i] = make([]int, n)
		for j := 0; j < n; j++ {
			m1[i][j] = rand.Intn(10)
			m2[i][j] = rand.Intn(10)
		}
	}
	// Se realiza la multiplicación y se calcula la traza en el proceso.
	trace := 0
	for i := 0; i < n; i++ {
		sum := 0
		for k := 0; k < n; k++ {
			sum += m1[i][k] * m2[k][i]
		}
		trace += sum
	}
	return trace
}

func EscribirArchivo(nombreArchivo *os.File, contenido string) {
	_, err := fmt.Fprintln(nombreArchivo, contenido)

	if err != nil {
		panic(err)
	}
}

func calcularPromedio(tiempos []time.Duration) time.Duration {
	var suma int64
	for _, t := range tiempos {
		suma += t.Milliseconds()
	}
	return time.Duration(suma/int64(len(tiempos))) * time.Millisecond
}

func EjecutarEspeculativoDetallado(blockData string, dificultad int, maxPrimos int, dimMatriz int, umbral int, nombreArchivo string) {
	inicioEjecucion := time.Now()

	arch, err := os.Create(nombreArchivo)
	if err != nil {
		panic(err)
	}
	defer arch.Close()

	resultadoA := make(chan ResultPOW, 1)
	resultadoB := make(chan []int, 1)
	cancelarA := make(chan struct{})
	cancelarB := make(chan struct{})

	// Rama A
	inicioRamaA := time.Now()
	go SimularProofOfWork(blockData, dificultad, resultadoA, cancelarA)

	// Rama B
	inicioRamaB := time.Now()
	go EncontrarPrimos(maxPrimos, resultadoB, cancelarB)

	resultTraza := CalcularTrazaDeProductoDeMatrices(dimMatriz)

	var seleccionada string
	var finalRamaA, finalRamaB time.Time
	var res any

	// Se selecciona la rama A
	if resultTraza > umbral {
		seleccionada = "A"
		close(cancelarB)
		finalRamaB = time.Now()
		res = <-resultadoA
		finalRamaA = time.Now()
	} else {
		seleccionada = "B"
		close(cancelarA)
		finalRamaA = time.Now()
		res = <-resultadoB
		finalRamaB = time.Now()
	}

	totalEjecucion := time.Since(inicioEjecucion)

	durA := finalRamaA.Sub(inicioRamaA)
	durB := finalRamaB.Sub(inicioRamaB)

	EscribirArchivo(arch, fmt.Sprintf("Tiempo total de ejecución: %d ms", totalEjecucion.Milliseconds()))
	EscribirArchivo(arch, fmt.Sprintf("\nRama A\n- Inicio: %d ms\n- Final: %d ms\n- Duración: %d ms", inicioRamaA.UnixMilli(), finalRamaA.UnixMilli(), durA.Milliseconds()))
	EscribirArchivo(arch, fmt.Sprintf("\nRama B\n- Inicio: %d ms\n- Final: %d ms\n- Duración: %d ms", inicioRamaB.UnixMilli(), finalRamaB.UnixMilli(), durB.Milliseconds()))

	EscribirArchivo(arch, fmt.Sprintf("\nTraza calculada: %d", resultTraza))
	EscribirArchivo(arch, fmt.Sprintf("Umbral: %d", umbral))
	EscribirArchivo(arch, fmt.Sprintf("\nRama seleccionada: %s", seleccionada))

	if seleccionada == "A" {
		EscribirArchivo(arch, "\nResultado:")
		EscribirArchivo(arch, fmt.Sprintf("- Hash: %v", res.(ResultPOW).Hash))
		EscribirArchivo(arch, fmt.Sprintf("- Nonce: %v", res.(ResultPOW).Nonce))
	} else {
		EscribirArchivo(arch, fmt.Sprintf("\nResultado: %v", res))
	}

	fmt.Printf("Resultados escritos en %s", nombreArchivo)
}

func EjecutarSecuencial(blockData string, dificultad int, maxPrimos int, dimMatriz int, umbral int) time.Duration {
	inicioEjecucion := time.Now()
	traza := CalcularTrazaDeProductoDeMatrices(dimMatriz)

	if traza > umbral {
		_ = SimularProofOfWorkSecuencial(blockData, dificultad)
		totalEjecucion := time.Since(inicioEjecucion)

		return totalEjecucion
	} else {
		_ = EncontrarPrimosSecuencial(maxPrimos)
		totalEjecucion := time.Since(inicioEjecucion)

		return totalEjecucion
	}
}

func EjecutarEspeculativo(blockData string, dificultad int, maxPrimos int, dimMatriz int, umbral int) time.Duration {
	inicioEjecucion := time.Now()

	// Se establecen los canales para resultados y señal de cancelación.
	resultadoA := make(chan ResultPOW, 1)
	resultadoB := make(chan []int, 1)
	cancelarA := make(chan struct{})
	cancelarB := make(chan struct{})

	// Lanza ambas goroutines.
	go SimularProofOfWork(blockData, dificultad, resultadoA, cancelarA)
	go EncontrarPrimos(maxPrimos, resultadoB, cancelarB)

	resultTraza := CalcularTrazaDeProductoDeMatrices(dimMatriz)

	// Selección de rama según resultado de la traza y umbral
	if resultTraza > umbral {
		close(cancelarB)
		_ = <-resultadoA
	} else {
		close(cancelarA)
		_ = <-resultadoB
	}

	return time.Since(inicioEjecucion)
}

func GuardarBenchmark(nombreArchivo, blockData string, dificultad, maxPrimos, dimMatriz, umbral int,
	tiemposEsp, tiemposSeq []time.Duration, promedioEsp, promedioSeq time.Duration, speedup float64) {

	arch, err := os.Create(nombreArchivo)
	if err != nil {
		panic(err)
	}
	defer arch.Close()

	fmt.Fprintf(arch, "-- BENCHMARK: 30 EJECUCIONES --\n")
	fmt.Fprintf(arch, "Parámetros:\n")
	fmt.Fprintf(arch, " - Block Data:        %s\n", blockData)
	fmt.Fprintf(arch, " - Dificultad PoW:    %d\n", dificultad)
	fmt.Fprintf(arch, " - Max Primos:        %d\n", maxPrimos)
	fmt.Fprintf(arch, " - Dimensión Matriz:  %d\n", dimMatriz)
	fmt.Fprintf(arch, " - Umbral:            %d\n\n", umbral)

	fmt.Fprintf(arch, "-- RESULTADOS --\n")
	fmt.Fprintf(arch, " Promedio Especulativo: %d ms\n", promedioEsp.Milliseconds())
	fmt.Fprintf(arch, " Promedio Secuencial:   %d ms\n", promedioSeq.Milliseconds())
	fmt.Fprintf(arch, " Speedup:               x%.2f\n\n", speedup)

	fmt.Fprintf(arch, "Tiempos Especulativos:\n")
	for i, t := range tiemposEsp {
		fmt.Fprintf(arch, "  Ejecución %2d: %d ms\n", i+1, t.Milliseconds())
	}

	fmt.Fprintf(arch, "\nTiempos Secuenciales:\n")
	for i, t := range tiemposSeq {
		fmt.Fprintf(arch, "  Ejecución %2d: %d ms\n", i+1, t.Milliseconds())
	}

	fmt.Printf("Resultados escritos en %s", nombreArchivo)
}

func main() {
	rand.Seed(1)

	if len(os.Args) != 7 && len(os.Args) != 8 {
		fmt.Println("Error: se requieren 6 argumentos: <blockData> <dificultad> <maxPrimos> <n> <umbral> <nombreArchivo> [--bmmode: modo benchmark opcional]")
		return
	}

	modoBenchmark := false
	if len(os.Args) == 8 {
		if os.Args[7] == "--bmmode" {
			modoBenchmark = true
		} else {
			fmt.Println("Error: el séptimo argumento debe ser '--bmmode' para modo benchmark.")
			return
		}
	}

	blockData := os.Args[1]
	nombreArchivo := os.Args[6]

	// Parseo de argumentos numéricos.
	dificultad, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error: la dificultad debe ser un número entero.")
		return
	}

	maxPrimos, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Error: el número de primos debe ser un número entero.")
		return
	}

	dimMatriz, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println("Error: la dimensión de la matriz debe ser un número entero.")
		return
	}

	umbral, err := strconv.Atoi(os.Args[5])
	if err != nil {
		fmt.Println("Error: el umbral debe ser un número entero.")
		return
	}

	if modoBenchmark {
		// Modo benchmark: ejecutar 30 veces cada modo y guardar resultados.
		tiemposEsp := make([]time.Duration, 30)
		tiemposSeq := make([]time.Duration, 30)

		// Ejecutar especulativo
		fmt.Println("Ejecutando modo especulativo...")
		for i := 0; i < 30; i++ {
			fmt.Printf("  Progreso: %d/%d\r", i+1, 30)
			tiemposEsp[i] = EjecutarEspeculativo(blockData, dificultad, maxPrimos, dimMatriz, umbral)
		}

		// Ejecutar secuencial
		fmt.Println("\nEjecutando modo secuencial...")
		for i := 0; i < 30; i++ {
			fmt.Printf("  Progreso: %d/%d\r", i+1, 30)
			tiemposSeq[i] = EjecutarSecuencial(blockData, dificultad, maxPrimos, dimMatriz, umbral)
		}

		// Calcular estadísticas
		promedioEsp := calcularPromedio(tiemposEsp)
		promedioSeq := calcularPromedio(tiemposSeq)
		speedup := float64(promedioSeq.Milliseconds()) / float64(promedioEsp.Milliseconds())

		GuardarBenchmark(nombreArchivo, blockData, dificultad, maxPrimos, dimMatriz, umbral, tiemposEsp, tiemposSeq, promedioEsp, promedioSeq, speedup)

	} else {
		// Modo normal: ejecuta el modo especulativo.
		EjecutarEspeculativoDetallado(blockData, dificultad, maxPrimos, dimMatriz, umbral, nombreArchivo)
	}
}
