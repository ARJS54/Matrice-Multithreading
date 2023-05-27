/*
Grade:              A1
Jakub Pažej         18260179
Aleksandr Jakusevs  18257038
*/

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Matrix [][]int //better than typing out [][]int every time

func printMat(inM Matrix) { //convenient function to print out the matrix
	fmt.Println()
	for _, i := range inM {
		for _, j := range i {
			fmt.Print(" ", j)
		}
		fmt.Println()
	}
}

func printMat2(inM *Matrix) { //convenient function to print out the matrix
	fmt.Println()
	for _, i := range *inM {
		for _, j := range i {
			fmt.Print(" ", j)
		}
		fmt.Println()
	}
}

func rowCount(inM Matrix) int { //counts number of rows in Matrix
	return len(inM)
}

func colCount(inM Matrix) int { //counts number of columns in Matrix
	return len(inM[0])
}

func rowCount2(inM *Matrix) int { //counts number of rows in Matrix
	return len(*inM)
}

func colCount2(inM *Matrix) int { //counts number of columns in Matrix
	return len((*inM)[0])
}

func newMatrix(r, c int) [][]int { //creates a new matrix and passes it by value
	a := make([]int, c*r)
	m := make([][]int, r)
	lo, hi := 0, c
	for i := range m {
		m[i] = a[lo:hi:hi]
		lo, hi = hi, hi+c
	}
	return m
}

func newMatrix2(r, c int) Matrix { //creates a new matrix and passes it by value
	a := make([]int, c*r)
	m := make([][]int, r)
	lo, hi := 0, c
	for i := range m {
		m[i] = a[lo:hi:hi]
		lo, hi = hi, hi+c
	}
	return m
}

func workerCalc(l int, m1, m2, m3 map[int]int, inA, inB Matrix, nM [][]int, lock *sync.Mutex) { //Calculates the invidual loop of a for(for(for)) type configuration using maps
	total := inA[m1[l]][m2[l]] * inB[m2[l]][m3[l]]
	lock.Lock()
	nM[m1[l]][m3[l]] = nM[m1[l]][m3[l]] + total
	lock.Unlock()
}

func doCalc(inA Matrix, inB Matrix) [][]int {
	var i int            //counter variables for loops
	lock := sync.Mutex{} //Creates a lock for the gouroutines
	m := rowCount(inA)   //number of rows of the first matrix
	//n := colCount(inA)     	//number of columns of the first matrix
	p := rowCount(inB) //number of rows of the second matrix
	q := colCount(inB) //number of columns of the second matrix

	map1 := make(map[int]int)
	map2 := make(map[int]int)
	map3 := make(map[int]int) //Creates the maps used to figure which part of the matrix has to be calculated

	nM := newMatrix(m, q) //create new matrix (to return at the end)

	start := time.Now() //starttimers

	for i = 0; i < m*p*q; i++ { //Fills in the maps used to calculate the code
		map1[i] = i % m
		map2[i] = int((i % (m * p)) / m)
		map3[i] = int(i / (m * p))
	}

	/*
		for i = 0; i < m*p*q; i++ {
			fmt.Println("Loop:", i, "Map1:", map1[i], "Map2:", map2[i], "Map3:", map3[i])
		}
	*/

	for i = 0; i < m*p*q; i++ { //passes the i value to signify which loop it is in for each goroutine.
		go workerCalc(i, map1, map2, map3, inA, inB, nM, &lock)
	}

	time.Sleep(1 * time.Second)
	var elapsed = time.Since(start) - 1*time.Second
	fmt.Println("Time taken to calculate: ", elapsed)
	return nM
}

func mainFunc1() {
	// Use slices
	// Unlike arrays they are passed by reference,not by value
	//a := Matrix{{2, 3, 3}, {5, 6, 5}, {9, 6, 3}}
	//b := Matrix{{8, 18, 28}, {38, 48, 58}, {69, 87, 1}}

	a := Matrix{{2, 3}, {5, 6}, {9, 6}}
	b := Matrix{{8, 18, 28}, {38, 48, 58}}

	fmt.Println("\nMatrix A")
	fmt.Println("Number of cols in A ", colCount(a))
	printMat(a)

	fmt.Println("\nMatrix B")
	fmt.Println("Number of rows in B ", rowCount(b))
	printMat(b)

	fmt.Println("\nThe Go Result of Matrix Multiplication:")
	c := doCalc(a, b)
	printMat(c)
}

func splitLoopCalc(p, q, i int, inA, inB, nM *Matrix, wg *sync.WaitGroup) { //function to calculate a row in a matrix
	var line []int  //line array that is added into the Matrix line at the end
	defer wg.Done() //wg-1 when function finishes

	var j, k, total int     //counter variables for loops
	for j = 0; j < q; j++ { //columns of the second matrix
		for k = 0; k < p; k++ { //rows of the second matrix
			total = total + (*inA)[i][k]*(*inB)[k][j] //calculates one column
			fmt.Print("(", (*inA)[i][k], " * ", (*inB)[k][j], ") + ")
		}
		fmt.Println("giving", total)
		line = append(line, total) //appends the column at the end of the line array
		//fmt.Println(line)
		total = 0
	}
	fmt.Println()
	(*nM)[i] = line //inserts the line array into the right position in the matrix
}

func splitLoopCalcQuietly(p, q, i int, inA, inB *Matrix, nM *Matrix, wg *sync.WaitGroup) { //basically just splitLoopCals without comments, used by the avgRunTime function
	var line []int
	defer wg.Done()
	var j, k, total int
	for j = 0; j < q; j++ {
		for k = 0; k < p; k++ {
			total = total + (*inA)[i][k]*(*inB)[k][j]
		}
		line = append(line, total)
		total = 0
	}
	(*nM)[i] = line
}

func doCalc2(inA, inB *Matrix) *Matrix {
	var i int //counter variables for loops

	m := rowCount2(inA) //number of rows of the first matrix
	//n := colCount(inA)     	//number of columns of the first matrix
	p := rowCount2(inB) //number of rows of the second matrix
	q := colCount2(inB) //number of columns of the second matrix

	nM := newMatrix2(m, q) //create new matrix (to return at the end)

	var wg sync.WaitGroup
	wg.Add(m)

	for i = 0; i < m; i++ { //rows of the first matrix
		go splitLoopCalc(p, q, i, inA, inB, &nM, &wg) //concurrency magic happens here
	}
	wg.Wait()
	return &nM
}

func doCalcQuietly(inA, inB *Matrix, wgr *sync.WaitGroup) *Matrix { //basically just doCalc without comments, used by the avgRunTime function
	defer wgr.Done()
	var i int
	m := rowCount2(inA)
	p := rowCount2(inB)
	q := colCount2(inB)
	nM := newMatrix2(m, q)
	var wg sync.WaitGroup
	wg.Add(m)
	for i = 0; i < m; i++ {
		go splitLoopCalcQuietly(p, q, i, inA, inB, &nM, &wg)
	}
	wg.Wait()
	return &nM
}

func avgRunTime(repeats int, inA, inB *Matrix) { // calculates avg run time of multiplications by running a lot of them then averaging time taken
	var i int
	var cumulative = time.Nanosecond
	var wgr sync.WaitGroup
	wgr.Add(repeats)
	start := time.Now()
	for i = 0; i < repeats; i++ {
		doCalcQuietly(inA, inB, &wgr)
	}
	wgr.Wait() // wait for all goroutines to finish
	var temp = time.Since(start)
	fmt.Println("Time taken to do ", repeats, " multiplications:", temp)
	cumulative = time.Duration(int64(temp) / int64(repeats))
	fmt.Println("Average time taken by one matrix multiplication:", cumulative)
}

func mainFunc2() {
	// Use slices
	// Unlike arrays they are passed by reference,not by value
	a := Matrix{{2, 3}, {5, 6}, {9, 6}}
	b := Matrix{{8, 18, 28}, {38, 48, 58}}

	fmt.Println("\nMatrix A")
	fmt.Println("Number of cols in A ", colCount2(&a))
	printMat2(&a)

	fmt.Println("\nMatrix B")
	fmt.Println("Number of rows in B ", rowCount2(&b))
	printMat2(&b)

	fmt.Println("\nThe Go Result of Matrix Multiplication:")
	c := doCalc2(&a, &b)
	printMat2(c)
	//avgRunTime(10000000, &a, &b)
}

func splitLoopCalc2(l, m, p, q int, inA, inB, nM *Matrix, wg *sync.WaitGroup) { //function to calculate a row in a matrix
	var line []int  //line array that is added into the Matrix line at the end
	defer wg.Done() //wg-1 when function finishes
	numbers := ""
	var i, j, k int //counter variables for loops
	total := 0
	for i = l; i < m; i = i + 2 {
		for j = 0; j < q; j++ { //columns of the second matrix
			for k = 0; k < p; k++ { //rows of the second matrix
				total = total + (*inA)[i][k]*(*inB)[k][j] //calculates one column
				numbers = numbers + "(" + strconv.Itoa((*inA)[i][k]) + " * " + strconv.Itoa((*inB)[k][j]) + ") + "
			}
			fmt.Println(numbers, "giving", total)
			numbers = ""
			line = append(line, total) //appends the column at the end of the line array
			//fmt.Println(line)
			total = 0
		}
		fmt.Println()
		(*nM)[i] = line //inserts the line array into the right position in the matrix
		line = nil
	}
}

func splitLoopCalcQuietly2(l, m, p, q int, inA, inB, nM *Matrix, wg *sync.WaitGroup) { //basically just splitLoopCals without printouts, used by the avgRunTime function
	var line []int  //line array that is added into the Matrix line at the end
	defer wg.Done() //wg-1 when function finishes
	var i, j, k int //counter variables for loops
	total := 0
	for i = l; i < m; i = i + 2 {
		for j = 0; j < q; j++ { //columns of the second matrix
			for k = 0; k < p; k++ { //rows of the second matrix
				total = total + (*inA)[i][k]*(*inB)[k][j] //calculates one column
			}
			line = append(line, total) //appends the column at the end of the line array
			//fmt.Println(line)
			total = 0
		}
		(*nM)[i] = line //inserts the line array into the right position in the matrix
		line = nil
	}
}

func doCalc3(inA, inB *Matrix) *Matrix {

	m := rowCount2(inA) //number of rows of the first matrix
	//n := colCount(inA)     	//number of columns of the first matrix
	p := rowCount2(inB) //number of rows of the second matrix
	q := colCount2(inB) //number of columns of the second matrix

	nM := newMatrix2(m, q) //create new matrix (to return at the end)

	var wg sync.WaitGroup
	wg.Add(2)

	go splitLoopCalc2(0, m, p, q, inA, inB, &nM, &wg) //concurrency magic happens here
	splitLoopCalc2(1, m, p, q, inA, inB, &nM, &wg)    //concurrency magic doesn't happen here but its still magical

	wg.Wait()
	return &nM
}

func doCalcQuietly2(inA, inB *Matrix, wgr *sync.WaitGroup) *Matrix { //basically just doCalc without comments, used by the avgRunTime function
	defer wgr.Done()
	m := rowCount2(inA)
	p := rowCount2(inB)
	q := colCount2(inB)
	nM := newMatrix2(m, q)
	var wg sync.WaitGroup
	wg.Add(2)
	go splitLoopCalcQuietly2(0, m, p, q, inA, inB, &nM, &wg)
	splitLoopCalcQuietly2(1, m, p, q, inA, inB, &nM, &wg)
	wg.Wait()
	return &nM
}

func avgRunTime2(repeats int, inA, inB *Matrix) { // calculates avg run time of multiplications by running a lot of them then averaging time taken
	var i int
	var cumulative = time.Nanosecond
	var wgr sync.WaitGroup
	wgr.Add(repeats)
	var idk = time.Nanosecond
	start := time.Now()
	for i = 0; i < repeats; i++ {
		var diff = time.Now()
		var a = generateRandomSlice(10, 0, 3, 2)
		var b = generateRandomSlice(10, 0, 2, 3)
		idk = time.Duration(int64(time.Since(diff)) + int64(idk))
		doCalcQuietly2(&a, &b, &wgr)
	}
	wgr.Wait() // wait for all goroutines to finish
	var temp = time.Since(start)
	fmt.Println("Time taken to do ", repeats, " multiplications:", temp)
	cumulative = time.Duration((int64(temp) - (int64(idk))) / int64(repeats))
	fmt.Println("Average time taken by one matrix multiplication:", cumulative)
}

func generateRandomSlice(max, min, rows, cols int) Matrix {
	rand.Seed(time.Now().UnixNano())
	slice := make([][]int, 0)
	for i := 0; i < rows; i++ {
		var line []int
		for j := 0; j < cols; j++ {
			line = append(line, rand.Intn(max-min)+min)
		}
		slice = append(slice, line)
	}
	return slice
}

func mainFunc3() {
	// Use slices
	// Unlike arrays they are passed by reference,not by value
	a := generateRandomSlice(10, 0, 3, 2)
	b := generateRandomSlice(10, 0, 2, 3)

	fmt.Println("\nMatrix A")
	fmt.Println("Number of cols in A ", colCount2(&a))
	printMat2(&a)

	fmt.Println("\nMatrix B")
	fmt.Println("Number of rows in B ", rowCount2(&b))
	printMat2(&b)

	fmt.Println("\nThe Go Result of Matrix Multiplication:")
	c := doCalc3(&a, &b)
	printMat2(c)
	//avgRunTime(1000000, &a, &b)
}

func main() {
	//mainFunc1()
	mainFunc2()
	//mainFunc3()
}
