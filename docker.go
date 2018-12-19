package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func main() {
	//var rLimit syscall.Rlimit
	//rLimit.Cur = 512
	//rLimit.Max = 2048

	/*if err:=syscall.Setrlimit(syscall.RLIMIT_DATA,&rLimit);err!=nil {
		fmt.Println("Error Setting Rlimit", err)
	}*/
	switch os.Args[1] {
	case "run":
		cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
		}


		/*if err:=syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit);err!=nil{
			fmt.Println("error getting rlimit",err)
		}
		fmt.Println(rLimit)*/

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil{
			fmt.Println("error ", err)
			os.Exit(1)
		}
		if err := HandleCGroup(cmd.Process.Pid);err!=nil {
			log.Fatal(err);
	}
	case "child":
		syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND,"")
		os.MkdirAll("rootfs/oldrootfs",0700)
		syscall.PivotRoot("rootfs", "rootfs/oldrootfs")
		os.Chdir("/")
		cmd:= exec.Command(os.Args[2],os.Args[3:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run();err != nil{
			fmt.Println("error",err)
			os.Exit(1)
		}

	default:
		panic("Enter run or child")
	}
}

func HandleCGroup(pid int) error{
	os.Mkdir("/sys/fs/cgroup/memory/CustomDocker", 0777)
	if err := WriteToFile("/sys/fs/cgroup/memory/CustomDocker/tasks", os.O_WRONLY | os.O_APPEND, strconv.Itoa(pid)); err != nil {
		return err
	}
	if err := WriteToFile("/sys/fs/cgroup/memory/CustomDocker/memory.limit_in_bytes",os.O_WRONLY,"40M");err !=nil {
		return err
	}

	return nil;
}

func WriteToFile(path string, flag int, content string) error{
	file, err := os.OpenFile(path, flag ,0666)
	defer file.Close()
	if err != nil {
		return err
	}
	if _, err := file.WriteString(content); err != nil{
		return err
	}
	return nil
}