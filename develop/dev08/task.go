package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
	Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:


	- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
	- pwd - показать путь до текущего каталога
	- echo <args> - вывод аргумента в STDOUT
	- kill <args> - "убить" процесс, переданный в качестве аргумента (пример: такой-то пример)
	- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*
	Так же требуется поддерживать функционал fork/exec-команд
	Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).
	*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
	в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
	и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор,
	пока не будет введена команда выхода (например \quit).
*/

const (
	CmdEcho     = "echo"
	CmdPwd      = "pwd"
	CmdCd       = "cd"
	CmdForkExec = "fork/exec"
	CmdPs       = "ps"
	CmdKill     = "kill"
)

// Бесконечный цикл обработки shell
func Shell() {
	sc := bufio.NewScanner(os.Stdin)
	for fmt.Print(">"); sc.Scan(); fmt.Print(">") {
		cmd := sc.Text()
		if cmd == "quit" {
			break
		}
		execCmd(cmd)
	}
}

// Обрабатывает введенную команду
func execCmd(cmd string) {
	args := strings.Fields(cmd)

	switch args[0] {
	case CmdEcho:
		if len(args) < 2 {
			fmt.Println("missing argument for echo")
			return
		}
		echo(strings.Join(args[1:], " "))
	case CmdPwd:
		pwd()
	case CmdCd:
		if len(args) < 2 {
			fmt.Println("missing argument for cd")
			return
		}
		cd(args[1])
	case CmdForkExec:
		if len(args) < 2 {
			fmt.Println("missing argument for fork/exec")
			return
		}
		forkExec(args[1])
	case CmdPs:
		ps()
	case CmdKill:
		if len(args) < 2 {
			fmt.Println("missing argument for kill")
			return
		}
		kill(args[1])
	default:
		fmt.Println("unknown command")
	}
}

// Выводит аргумент в shell
func echo(arg string) {
	fmt.Println(arg)
}

// Выводит рабочую директорию в shell
func pwd() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("pwd: %v", err)
		return
	}

	fmt.Println(path)
}

// Меняет рабочую директорию
func cd(dir string) {
	if err := os.Chdir(dir); err != nil {
		fmt.Println("no such directory")
		return
	}
}

// Запускает новый процесс, с соответствующими аргументами
func forkExec(arg string) {
	args := strings.Split(arg, " ")

	if len(args) < 1 {
		fmt.Println("invalid arguments")
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	go func() {
		if err := cmd.Run(); err != nil {
			fmt.Printf("fork/exec: %v", err)
		}
	}()
}

// Выводит список процессов
func ps() {
	output, err := exec.Command("ps", "-e").Output()
	if err != nil {
		fmt.Printf("ps: %v", err)
	}

	fmt.Println(string(output))
}

// Убивает процесс по PID
func kill(pidStr string) {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Printf("kill: %v", err)
		return
	}

	prc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("kill: %v", err)
		return
	}

	if err := prc.Kill(); err != nil {
		fmt.Printf("cannot kill process: %v", err)
	}
}

func main() {
	Shell()
}
