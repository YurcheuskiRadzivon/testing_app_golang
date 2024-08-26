package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type clock struct {
	hour   int
	minute int
}
type sysParametrs struct {
	timenow    clock
	stTime     clock
	enTime     clock
	hourP      int
	lockTable  int
	countTable int
}
type client struct {
	timeIn    clock
	timeStTab clock
	numTab    int
}
type table struct {
	nameClient string
	revenue    int
	sumTime    clock
}

// Обработчик ошибок
func errorHandling(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
func clockComp(time1, time2 clock, fl *int) {
	if time2.hour < time1.hour {
		if *fl == 0 {
			errorHandling(errors.New("Time_input_error....."))
		}
		*fl = -1

	} else if time2.hour == time1.hour && time2.minute < time1.minute {
		if *fl == 0 {
			errorHandling(errors.New("Time_input_error....."))
		}
		*fl = -1
	}

}
func lenCheker(size int, n1 int, n2 int) {
	if size != n1 && size != n2 {
		err := errors.New("Invalid string format")
		errorHandling(err)
	}
}
func formatChecker(s string) {
	for _, char := range s {
		if !(char > '0' && char < '9') && !(unicode.IsLower(char)) && char != '_' && char != '-' {
			err := errors.New("Input name format error")
			errorHandling(err)
		}

	}
}

func fInf(scanner *bufio.Scanner) (int, [2]clock, int, [2]string) {
	scanner.Scan()
	num, err := strconv.Atoi(scanner.Text())
	errorHandling(err)

	scanner.Scan()
	var cl [2]clock
	var cl_str [2]string
	time_str := strings.Split(scanner.Text(), " ")
	lenCheker(len(time_str), 2, 2)
	for i, t := range time_str {
		cl_str[i] = t
		lenCheker(len(t), 5, 5)
		parsedTime, err := time.Parse("15:04", t)
		errorHandling(err)
		cl[i].hour = parsedTime.Hour()
		cl[i].minute = parsedTime.Minute()

	}
	fl := 0
	clockComp(cl[0], cl[1], &fl)

	scanner.Scan()
	payment, err := strconv.Atoi(scanner.Text())
	errorHandling(err)
	return num, cl, payment, cl_str

}
func id13(timeEr clock, strErr string, writer *bufio.Writer) {
	writer.WriteString(fmt.Sprintf("%02d:%02d 13 %v\n", timeEr.hour, timeEr.minute, strErr))

}
func id12(timeEn clock, strName string, writer *bufio.Writer, numTab int) {
	writer.WriteString(fmt.Sprintf("%02d:%02d 12 %v %v\n", timeEn.hour, timeEn.minute, strName, numTab))

}
func id11(timeEn clock, strName string, writer *bufio.Writer) {
	writer.WriteString(fmt.Sprintf("%02d:%02d 11 %v\n", timeEn.hour, timeEn.minute, strName))

}
func paymTab(system *sysParametrs, tables []table, gamer *client) {
	timeGaming := system.timenow.hour*60 + system.timenow.minute - gamer.timeStTab.hour*60 - gamer.timeStTab.minute
	tables[gamer.numTab-1].sumTime.minute += timeGaming % 60
	tables[gamer.numTab-1].sumTime.hour += timeGaming / 60
	if timeGaming%60 == 0 {
		tables[gamer.numTab-1].revenue += ((timeGaming / 60) * system.hourP)
	} else {
		tables[gamer.numTab-1].revenue += (((timeGaming / 60) + 1) * system.hourP)
	}
}
func eventHandler(tables []table, numEv int, system *sysParametrs, clientInf map[string]client, queue *list.List, information []string, writer *bufio.Writer) {
	name := information[2]
	formatChecker(name)
	gamer, ok := clientInf[name]
	if numEv == 1 {
		lenCheker(len(information), 3, 3)
		flSt, flEn := 1, 1
		clockComp(system.stTime, system.timenow, &flSt)
		clockComp(system.timenow, system.enTime, &flEn)
		if !ok && flSt == 1 && flEn == 1 {
			clientInf[name] = client{}
			gamer = clientInf[name]
			gamer.numTab = 0
			gamer.timeIn = system.timenow
			gamer.timeStTab = clock{-1, -1}
			clientInf[name] = gamer

		} else {
			if ok {
				id13(system.timenow, "YouShallNotPass", writer)
			} else {
				id13(system.timenow, "NotOpenYet", writer)
			}
		}
	} else if numEv == 2 {
		lenCheker(len(information), 4, 4)
		id_table, err := strconv.Atoi(information[3])
		errorHandling(err)

		if !ok {
			id13(system.timenow, "ClientUnknown", writer)
		} else {
			if tables[id_table-1].nameClient == "" {

				if gamer.numTab == 0 {
					gamer.numTab = id_table
					gamer.timeStTab = system.timenow
					tables[id_table-1].nameClient = name
					system.lockTable++
					clientInf[name] = gamer
				} else {
					paymTab(system, tables, &gamer)
					tables[gamer.numTab-1].nameClient = ""
					tables[id_table-1].nameClient = name
					gamer.numTab = id_table
					gamer.timeStTab = system.timenow
					clientInf[name] = gamer
				}

			} else {
				id13(system.timenow, "PlaceIsBusy", writer)
			}

		}

	} else if numEv == 3 {
		lenCheker(len(information), 3, 3)
		if system.countTable-system.lockTable != 0 {
			id13(system.timenow, "ICanWaitNoLonger!", writer)
		} else if queue.Len() == system.countTable {
			id11(system.timenow, name, writer)
			delete(clientInf, name)
		} else {
			queue.PushBack(name)

		}

	} else if numEv == 4 {
		lenCheker(len(information), 3, 3)
		if !ok {
			id13(system.timenow, "ClientUnknown", writer)
		} else {
			paymTab(system, tables, &gamer)
			tables[gamer.numTab-1].nameClient = ""
			num := gamer.numTab

			delete(clientInf, name)
			if queue.Len() != 0 {

				name := fmt.Sprintf("%v", queue.Front().Value)
				queue.Remove(queue.Front())
				gamer := (clientInf)[name]
				gamer.numTab = num
				gamer.timeStTab = system.timenow
				tables[num-1].nameClient = name
				clientInf[name] = gamer
				id12(system.timenow, name, writer, num)

			} else {

				system.lockTable--
			}

		}
	} else {
		err := errors.New("Unknown event id")
		errorHandling(err)
	}

}
func endWorktime(id11Arr []string, clientInf map[string]client, writer *bufio.Writer, timeEnd clock, tables []table, system *sysParametrs) {
	for el := range clientInf {
		id11Arr = append(id11Arr, el)
	}
	sort.Strings(id11Arr)
	for _, el := range id11Arr {
		gamer := clientInf[el]
		paymTab(system, tables, &gamer)
		id11(timeEnd, el, writer)
		delete(clientInf, el)
	}

}
func revenueView(tables []table, writer *bufio.Writer) {
	for i, tab := range tables {
		writer.WriteString(fmt.Sprintf("%v %v %02d:%02d", (i + 1), tab.revenue, tab.sumTime.hour, tab.sumTime.minute))
		writer.WriteString("\n")
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("The name of the text file is not entered...")
		os.Exit(1)
	}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	errorHandling(err)
	writer := bufio.NewWriter(os.Stdout)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	n, workTime, hourPayment, clStr := fInf(scanner)
	writer.WriteString(clStr[0])
	writer.WriteString("\n")
	tables := make([]table, n)
	clientInf := make(map[string]client)
	var system sysParametrs
	system.countTable = n
	system.stTime = workTime[0]
	system.enTime = workTime[1]
	system.hourP = hourPayment
	queue := list.New()
	fmt.Printf(tables[2].nameClient)
	for scanner.Scan() {
		writer.WriteString(scanner.Text())
		writer.WriteString("\n")
		information := strings.Split(scanner.Text(), " ")
		lenCheker(len(information), 3, 4)
		t, err := time.Parse("15:04", information[0])
		errorHandling(err)
		var timeEv clock = clock{t.Hour(), t.Minute()}
		fl := 0
		clockComp(system.timenow, timeEv, &fl)
		system.timenow.hour = timeEv.hour
		system.timenow.minute = timeEv.minute
		numEv, err := strconv.Atoi(information[1])
		errorHandling(err)
		eventHandler(tables, numEv, &system, clientInf, queue, information, writer)

	}
	system.timenow = system.enTime

	id11Arr := make([]string, 0, len(clientInf))
	if len(clientInf) != 0 {
		endWorktime(id11Arr, clientInf, writer, system.enTime, tables, &system)
	}
	writer.WriteString(clStr[1])
	writer.WriteString("\n")
	revenueView(tables, writer)
	writer.Flush()

}
