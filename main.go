package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

/*
Написать класс "Очереди задач", которая будет использоваться в следующем контексте:
1. В очередь могут складывать задачи несколько горутин, количество которых известно в момент создания очереди.
2. Задачи из очереди исполняются ровно в одной горутине
3. Задачи содержат в себе все необходимые параметры для исполнения.
4. Задачи от одной и той же горутины должны быть исполнены в порядке FIFO
Для проверки решения реализуйте программу, принимающую в аргументах командной строки три параметра
-writers    - количество пишущих горутин
-arr-size   - размер массива (читайте далее)
-iter-count - количество итераций
Программа должна запустить горутины в количестве "-writers",
каждая из которых должна  "-iter-count"  раз сделать следующее:
1. Сгенерировать массив длины "-arr-size" случайных целых чисел.
2. Положить в очередь задач задачу на сортировку этого массива и вывода на экран информации в следующем формате:
"{writerGoroutineIdenitifier} {queueInsertionTime} {min} {median} {max}\n"
Таким образом у каждой пишущей горутины должен быть свой идентификатор,
у каждой задачи должен быть проставлен timestamp момента добавления в очередь.
min = array[0], median=array[size/2], max=array[size - 1].

main горутина выступает в роли исполнителя
*/
const DEFAULT_WRITERS = 2
const DEFAULT_ARR_SISE = 10
const DEFAULT_ITER_COUNT = 5
const ARRAY_VALUE_MIN = -99
const ARRAY_VALUE_MAX = 99

var wg sync.WaitGroup

//--------------- to_do Задачи от одной и той же горутины должны быть исполнены в порядке FIFO
/*type Queue interface {
	Add(id int, task Task)
	Get()
	Init()
}
type MyQueue struct {
	TaskMap map[int][]Task
}

func (q MyQueue) Get() {
	panic("implement me")
}


func (q MyQueue) Add(id int, task Task)  {
	q.TaskMap[id] = append(w.TaskMap[id], task)
	q.queue = append(q.queue, task)
}
func (q MyQueue) Init()  {
	w.TaskMap = make(map[int][]Task)
}*/

//--------------- to_do перенести в task/task.go
type Task interface {
	Do()
}
type Worker struct {
	WriterGoroutineIdenitifier int
	QueueInsertionTime         time.Time
	Arr                        []int
	IterCount                  int
}

func (w Worker) Do() {
	sort.Ints(w.Arr[:])
	fmt.Println(fmt.Sprintf(
		"{%d} {%s} {%d} {%d} {%d}\n",
		w.WriterGoroutineIdenitifier,
		w.QueueInsertionTime.Format(time.RFC3339Nano),
		w.Arr[0],
		w.Arr[len(w.Arr)/2],
		w.Arr[len(w.Arr)-1],
	))
}

type params struct {
	writers   int
	arrSize   int
	iterCount int
}

func readParams() *params {
	writers := flag.Int("writers", DEFAULT_WRITERS, "an int")
	arrSize := flag.Int("arr-size", DEFAULT_ARR_SISE, "an int")
	iterCount := flag.Int("iter-count", DEFAULT_ITER_COUNT, "an int")
	flag.Parse()
	p := params{
		writers:   *writers,
		arrSize:   *arrSize,
		iterCount: *iterCount,
	}
	return &p
}

func writer(id int, p params, tasks chan Task) {
	defer wg.Done()
	for i := 0; i < p.iterCount; i++ {
		var arr = arrayGen(p.arrSize)
		tasks <- Worker{
			WriterGoroutineIdenitifier: id,
			IterCount:                  i,
			Arr:                        arr,
			QueueInsertionTime:         time.Now(),
		}
	}

}

func arrayGen(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = ARRAY_VALUE_MIN + rand.Intn(ARRAY_VALUE_MAX-ARRAY_VALUE_MIN)
	}
	return arr
}

func worker(tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		task.Do()
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var p = readParams()
	var queue = make(chan Task, p.writers*p.iterCount) //to_do заменить на queue
	wg.Add(p.writers + 1)
	for i := 0; i < p.writers; i++ {
		go writer(i, *p, queue)
	}
	go worker(queue, &wg)
	wg.Wait()
}
