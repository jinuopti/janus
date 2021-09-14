package process

import (
	"errors"
	"os"
	"os/exec"

	"github.com/jinuopti/janus/database/gorm"
	"github.com/jinuopti/janus/database/gorm/processdb"
	. "github.com/jinuopti/janus/log"
	"github.com/shirou/gopsutil/v3/process"

	// "github.com/shirou/gopsutil/v3/host"
	// "github.com/shirou/gopsutil/v3/memory"
	// "github.com/shirou/gopsutil/v3/disk"
	// "github.com/shirou/gopsutil/v3/load"
	"strings"
	"syscall"
	"time"
)

type ProcInfo struct {
    // manual insert
    Name        string      // unique name, ex: PshmIOC
    Cmd         string
    Args        string
    Type        string      // process type, ex: IOC
    UseDb       bool
    AutoRestart bool

    // auto insert
    Pid         int
    Process     *process.Process
    Status      string      // "zombie" is dead

    // callback functions
    PreRunFunc  func(*ProcInfo)     // RunPorocess 실행 전 수행할 작업
    PostRunFunc func(*ProcInfo)     // RunPorocess 실행 후 수행할 작업
    PreDelFunc  func(*ProcInfo)     // DelPorocess 실행 전 수행할 작업
    PostDelFunc func(*ProcInfo)     // DelPorocess 실행 후 수행할 작업
    KilledCallbackFunc  func(*ProcInfo)     // process 가 중단될 때 호출
}

var (
    processMap map[string]*ProcInfo // 직접 실행한 process 관리
    processMapPid map[int]*ProcInfo
)

// GetRunningProcesses
// pName process 를 모두 검색
func GetRunningProcesses(pName string) []*process.Process {
    var pidList []*process.Process
    var count int
    psList, _ := process.Processes()
    for _, p := range psList {
        name, _ := p.Name()
        if name == pName {
            count++
            pidList = append(pidList, p)
        }
    }
    return pidList
}

// GetRunningProcessesContains
// 실행명령어 내 contain 문자열이 포함된 모든 process 검색
func GetRunningProcessesContains(contain string) []*process.Process {
    var pidList []*process.Process
    var count int
    psList, _ := process.Processes()
    myPid := os.Getpid()
    for _, p := range psList {
        cmdLine, _ := p.Cmdline()
        if strings.Contains(cmdLine, contain) && p.Pid != int32(myPid) {
            count++
            pidList = append(pidList, p)
        }
    }
    return pidList
}

// GetRunningProcess
// 현재 실행중인 process list 중 pid, cmd 일치하는 process 가 있는지 확인, 없으면 nil 반환
func GetRunningProcess(pid int, cmd string) *process.Process {
    psList, _ := process.Processes()
    for _, p := range psList {
        c, _ := p.Cmdline()
        if strings.Contains(c, cmd) && p.Pid == int32(pid) {
            return p
        }
    }
    return nil
}

func GetMyProcessAll() map[string]*ProcInfo {
    return processMap
}

func GetProcess(name string) *ProcInfo {
    return processMap[name]
}

func GetProcessPid(pid int) *ProcInfo {
    return processMapPid[pid]
}

func InsertNewProcess(proc *ProcInfo) {
    if processMap == nil {
        processMap = make(map[string]*ProcInfo)
    }
    if processMapPid == nil {
        processMapPid = make(map[int]*ProcInfo)
    }

    processMap[proc.Name] = proc
    processMapPid[proc.Pid] = proc
}

func (p *ProcInfo) DeleteProcess() {
    if p.PreDelFunc != nil {
        p.PreDelFunc(p)
    }

    _ = processdb.DeleteProcess(p.Name)
    delete(processMap, p.Name)
    delete(processMapPid, p.Pid)

    if p.PostDelFunc != nil {
        p.PostRunFunc(p)
    }
}

// RunProcess
// process 를 background 로 실행
func (p *ProcInfo) RunProcess() error {
    if p == nil || len(p.Name) < 1 || len(p.Cmd) < 1 {
        return errors.New("ProcInfo or Name or Cmd is nil")
    }

    if p.PreRunFunc != nil {
        p.PreRunFunc(p)
    }

    args := strings.Split(p.Args," ")
    cmd := exec.Command(p.Cmd, args...)
    //cmd := exec.Command("sh", "-c", proc.Cmd + " " + proc.Args)
    cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    err := cmd.Start()
    if err != nil {
        return err
    }
    
    pid := cmd.Process.Pid
    pp, err := process.NewProcess(int32(pid))
    if err != nil {
        return err
    }

    p.Pid = pid
    p.Process = pp

    Logd("Start Process [%s], pid:%d, cmd:[%s], args:[%s]", p.Name, p.Pid, p.Cmd, p.Args)

    if p.UseDb {
        record := &processdb.Processes{
            Name: p.Name,
            Pid: p.Pid,
            Cmd: p.Cmd,
            Args: p.Args,
            Type: p.Type,
            Running: true,
            AutoRestart: p.AutoRestart,
        }
        _ = processdb.InsertProcess(record)
    }

    InsertNewProcess(p)

    if p.PostRunFunc != nil {
        p.PostRunFunc(p)
    }

    go func() {
       err = cmd.Wait()
       Logd("Command finished with error: %v", err)
        if p.KilledCallbackFunc != nil {
            p.KilledCallbackFunc(p)
        }
        p.DeleteProcess()
    }()

    return nil
}

// KillProcess
// Process 를 kill 한다.
// name : process name
// killChild : true -> 자식의 w
func KillProcess(proc interface{}, killChild bool) error {
    var err error
    var p *ProcInfo

    switch proc.(type) {
    case string:
        p = GetProcess(proc.(string))
    case int:
        p = GetProcessPid(proc.(int))
    }

    if p == nil {
        return errors.New("process not found")
    } else {
        err = syscall.Kill(p.Pid, syscall.SIGKILL)
        if err != nil {
            return err
        }
    }

    if killChild {
        err = syscall.Kill(-p.Pid, syscall.SIGKILL) // 음수 pid 는 pgid (자식프로세스까지) 사용 모두 kill
    } else {
        err = p.Process.Kill()
    }
    if err != nil {
        return err
    }

    Logd("Process [%s] killed, pid:%d, cmd:[%s], args:[%s]", p.Name, p.Pid, p.Cmd, p.Args)

    p.DeleteProcess()

    return nil
}

// CheckProcessStatus
// R: Running S: Sleep T: Stop I: Idle Z: Zombie W: Wait L: Lock
func CheckProcessStatus() error {
    for _, proc := range processMap {
        fault := false
        status, err := proc.Process.Status()
        if err != nil {
            Logd("error, %s", err)
            fault = true
            goto CheckFault
        }

        proc.Status = status[0]
        //Logd("Process[%d][%s] Status: %s", proc.Pid, proc.Cmd, proc.Status)

        switch proc.Status {
        case process.Zombie:
            Loge("zombie process detected, kill process [%d][%s][%s]", proc.Pid, proc.Cmd, proc.Args)
            fault = true
        case process.Running:
        case process.Sleep:
        case process.Idle:
        case process.Lock:
        case process.Wait:
        }

CheckFault:
        if fault {
            isDelete := false
            _ = proc.Process.Kill()
            if proc.AutoRestart {       // 자동 재시작
                err = proc.RunProcess()
                if err != nil {
                    isDelete = true
                    Loge("error!! Auto RunProcess [%s]", proc.Name)
                } else {
                    Logd("Auto Restart process [%s]", proc.Name)
                }
            } else {
                isDelete = true
                Logd("Don't auto restart process [%s]", proc.Name)
            }
            if isDelete {
                if proc.KilledCallbackFunc != nil {
                    proc.KilledCallbackFunc(proc)
                }
                proc.DeleteProcess()
            }
        }
    }

    return nil
}

func MonitorProcessStatus() {
    ticker := time.NewTicker(time.Second * 1)

    for {
        select {
        case _ = <- ticker.C:
            _ = CheckProcessStatus()
            //Logd("Check Process Status.. [%s]", tick.String())
        }
    }
}

func DefaultKilledCallback(p *ProcInfo) {
    if p != nil {
        Logd("Process [%s] killed, pid:%d, cmd:[%s], args:[%s]", p.Name, p.Pid, p.Cmd, p.Args)
    }
}

// Test 테스트 function
// name: unique process name
// callback: process 상태 변경 시 호출되는 callback function
// cmd: 실행 process command
// args: 실행 process arguments
func Test(name string, callback func(*ProcInfo), cmd string, arg string) {
    // process 실행
    procInfo := &ProcInfo{
        Name: name,
        Cmd: cmd,
        Args: arg,
        UseDb: true,
        AutoRestart: true,
    }

    if procInfo.UseDb {
        err := gormdb.InitSingletonDB()
        if err == nil {
            processdb.InitProcessTable()
        }
    }

    err := procInfo.RunProcess()
    if err != nil {
        Loge("error, %s", err)
        return
    }
    Logd("run process [%s], pid:%d, cmd:[%s], args:[%s]", procInfo.Name, procInfo.Pid, procInfo.Cmd, procInfo.Args)

    // process monitor goroutine with callback function
    go MonitorProcessStatus()

    // 관리되고 있는 모든 process 출력
    pids := GetMyProcessAll()
    for i, pid := range pids {
        Logd("%s: pid[%d], cmd[%s], args[%s]", i, pid.Pid, pid.Cmd, pid.Args)
    }

    // 10초 sleep
    time.Sleep(time.Second * 10)

    // cmd or pid 검색
    p := GetProcess(name)
    // p := GetProcessPid(pid)
    if p != nil {
       Logd("Kill process [%s] pid:%d, cmd:[%s], args[%s]", p.Name, p.Pid, p.Cmd, p.Args)
       err = KillProcess(p.Name, true)
       if err != nil {
           Loge("error, %s", err)
       }
    } else {
        Logd("process not found")
    }
}
