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
type sys_parametrs struct {
	timenow     clock
	st_time     clock
	en_time     clock
	hour_p      int
	lock_table  int
	count_table int
}
type client struct {
	time_in     clock
	time_st_tab clock
	num_tab     int
}
type table struct {
	name_client string
	revenue     int
	sum_time    clock
}

// Обработчик ошибок
func error_handling(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
func cl_comp(time1, time2 clock, fl *int) {
	if time2.hour < time1.hour {
		if *fl == 0 {
			error_handling(errors.New("Time_input_error....."))
		}
		*fl = -1

	} else if time2.hour == time1.hour && time2.minute < time1.minute {
		if *fl == 0 {
			error_handling(errors.New("Time_input_error....."))
		}
		*fl = -1
	}

}
func len_cheker(size int, n1 int, n2 int) {
	if size != n1 && size != n2 {
		err := errors.New("Invalid string format")
		error_handling(err)
	}
}
func format_checker(s string) {
	for _, char := range s {
		if !(char > '0' && char < '9') && !(unicode.IsLower(char)) && char != '_' && char != '-' {
			err := errors.New("Input name format error")
			error_handling(err)
		}

	}
}

func f_inf(scanner *bufio.Scanner) (int, [2]clock, int, [2]string) {
	scanner.Scan()
	num, err := strconv.Atoi(scanner.Text())
	error_handling(err)

	scanner.Scan()
	var cl [2]clock
	var cl_str [2]string
	time_str := strings.Split(scanner.Text(), " ")
	len_cheker(len(time_str), 2, 2)
	for i, t := range time_str {
		cl_str[i] = t
		len_cheker(len(t), 5, 5)
		parsedTime, err := time.Parse("15:04", t)
		error_handling(err)
		cl[i].hour = parsedTime.Hour()
		cl[i].minute = parsedTime.Minute()

	}
	fl := 0
	cl_comp(cl[0], cl[1], &fl)

	scanner.Scan()
	payment, err := strconv.Atoi(scanner.Text())
	error_handling(err)
	return num, cl, payment, cl_str

}
func id_13(time_er clock, str_err string, writer *bufio.Writer) {
	writer.WriteString(fmt.Sprintf("%02d:%02d 13 %v\n", time_er.hour, time_er.minute, str_err))

}
func id_12(time_en clock, str_name string, writer *bufio.Writer, num_tab int) {
	writer.WriteString(fmt.Sprintf("%02d:%02d 12 %v %v\n", time_en.hour, time_en.minute, str_name, num_tab))

}
func id_11(time_en clock, str_name string, writer *bufio.Writer) {
	writer.WriteString(fmt.Sprintf("%02d:%02d 11 %v\n", time_en.hour, time_en.minute, str_name))

}
func paym_tab(system *sys_parametrs, tables []table, gamer *client) {
	time_gaming := system.timenow.hour*60 + system.timenow.minute - gamer.time_st_tab.hour*60 - gamer.time_st_tab.minute
	tables[gamer.num_tab-1].sum_time.minute += time_gaming % 60
	tables[gamer.num_tab-1].sum_time.hour += time_gaming / 60
	if time_gaming%60 == 0 {
		tables[gamer.num_tab-1].revenue += ((time_gaming / 60) * system.hour_p)
	} else {
		tables[gamer.num_tab-1].revenue += (((time_gaming / 60) + 1) * system.hour_p)
	}
}
func event_handler(tables []table, num_ev int, system *sys_parametrs, client_inf map[string]client, queue *list.List, information []string, writer *bufio.Writer) {
	name := information[2]
	format_checker(name)
	gamer, ok := client_inf[name]
	if num_ev == 1 {
		len_cheker(len(information), 3, 3)
		fl_st, fl_en := 1, 1
		cl_comp(system.st_time, system.timenow, &fl_st)
		cl_comp(system.timenow, system.en_time, &fl_en)
		if !ok && fl_st == 1 && fl_en == 1 {
			client_inf[name] = client{}
			gamer = client_inf[name]
			gamer.num_tab = 0
			gamer.time_in = system.timenow
			gamer.time_st_tab = clock{-1, -1}
			client_inf[name] = gamer

		} else {
			if ok {
				id_13(system.timenow, "YouShallNotPass", writer)
			} else {
				id_13(system.timenow, "NotOpenYet", writer)
			}
		}
	} else if num_ev == 2 {
		len_cheker(len(information), 4, 4)
		id_table, err := strconv.Atoi(information[3])
		error_handling(err)

		if !ok {
			id_13(system.timenow, "ClientUnknown", writer)
		} else {
			if tables[id_table-1].name_client == "" {

				if gamer.num_tab == 0 {
					gamer.num_tab = id_table
					gamer.time_st_tab = system.timenow
					tables[id_table-1].name_client = name
					system.lock_table++
					client_inf[name] = gamer
				} else {
					paym_tab(system, tables, &gamer)
					tables[gamer.num_tab-1].name_client = ""
					tables[id_table-1].name_client = name
					gamer.num_tab = id_table
					gamer.time_st_tab = system.timenow
					client_inf[name] = gamer
				}

			} else {
				id_13(system.timenow, "PlaceIsBusy", writer)
			}

		}

	} else if num_ev == 3 {
		len_cheker(len(information), 3, 3)
		if system.count_table-system.lock_table != 0 {
			id_13(system.timenow, "ICanWaitNoLonger!", writer)
		} else if queue.Len() == system.count_table {
			id_11(system.timenow, name, writer)
			delete(client_inf, name)
		} else {
			queue.PushBack(name)

		}

	} else if num_ev == 4 {
		len_cheker(len(information), 3, 3)
		if !ok {
			id_13(system.timenow, "ClientUnknown", writer)
		} else {
			paym_tab(system, tables, &gamer)
			tables[gamer.num_tab-1].name_client = ""
			num := gamer.num_tab

			delete(client_inf, name)
			if queue.Len() != 0 {

				name := fmt.Sprintf("%v", queue.Front().Value)
				queue.Remove(queue.Front())
				gamer := (client_inf)[name]
				gamer.num_tab = num
				gamer.time_st_tab = system.timenow
				tables[num-1].name_client = name
				client_inf[name] = gamer
				id_12(system.timenow, name, writer, num)

			} else {

				system.lock_table--
			}

		}
	} else {
		err := errors.New("Unknown event id")
		error_handling(err)
	}

}
func end_worktime(id_11_arr []string, client_inf map[string]client, writer *bufio.Writer, time_end clock, tables []table, system *sys_parametrs) {
	for el := range client_inf {
		id_11_arr = append(id_11_arr, el)
	}
	sort.Strings(id_11_arr)
	for _, el := range id_11_arr {
		gamer := client_inf[el]
		paym_tab(system, tables, &gamer)
		id_11(time_end, el, writer)
		delete(client_inf, el)
	}

}
func revenue_view(tables []table, writer *bufio.Writer) {
	for i, tab := range tables {
		writer.WriteString(fmt.Sprintf("%v %v %02d:%02d", (i + 1), tab.revenue, tab.sum_time.hour, tab.sum_time.minute))
		writer.WriteString("\n")
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("The name of the text file is not entered...")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	error_handling(err)
	writer := bufio.NewWriter(os.Stdout)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	n, work_time, hour_payment, cl_str := f_inf(scanner)
	writer.WriteString(cl_str[0])
	writer.WriteString("\n")
	tables := make([]table, n)
	client_inf := make(map[string]client)
	var system sys_parametrs
	system.count_table = n
	system.st_time = work_time[0]
	system.en_time = work_time[1]
	system.hour_p = hour_payment
	queue := list.New()
	fmt.Printf(tables[2].name_client)
	for scanner.Scan() {
		writer.WriteString(scanner.Text())
		writer.WriteString("\n")
		information := strings.Split(scanner.Text(), " ")
		len_cheker(len(information), 3, 4)
		t, err := time.Parse("15:04", information[0])
		error_handling(err)
		var time_ev clock = clock{t.Hour(), t.Minute()}
		fl := 0
		cl_comp(system.timenow, time_ev, &fl)
		system.timenow.hour = time_ev.hour
		system.timenow.minute = time_ev.minute
		num_ev, err := strconv.Atoi(information[1])
		error_handling(err)
		event_handler(tables, num_ev, &system, client_inf, queue, information, writer)

	}
	system.timenow = system.en_time

	id_11_arr := make([]string, 0, len(client_inf))
	if len(client_inf) != 0 {
		end_worktime(id_11_arr, client_inf, writer, system.en_time, tables, &system)
	}
	writer.WriteString(cl_str[1])
	writer.WriteString("\n")
	revenue_view(tables, writer)
    writer.Flush()

}
