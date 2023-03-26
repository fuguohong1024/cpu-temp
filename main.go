package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/StackExchange/wmi"
)

type Sensor struct {
	Name   string
	Parent string
	Value  float32
}

func main() {
	// 启动OpenHardwareMonitor
	ohmPath, err := filepath.Abs("OpenHardwareMonitor/OpenHardwareMonitor.exe")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ohm := exec.Command(ohmPath)
	ohm.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err = ohm.Start()
	if err != nil {
		fmt.Printf("Error starting OpenHardwareMonitor: %v\n", err)
		return
	}

	// 延迟关闭OpenHardwareMonitor
	defer func() {
		err = ohm.Process.Kill()
		if err != nil {
			fmt.Printf("Error stopping OpenHardwareMonitor: %v\n", err)
		}
	}()

	// 等待OpenHardwareMonitor启动
	time.Sleep(3 * time.Second)

	// 查询CPU温度
	var sensors []Sensor
	query := "SELECT Name, Parent, Value FROM Sensor WHERE SensorType = 'Temperature'"
	err = wmi.QueryNamespace(query, &sensors, "root\\OpenHardwareMonitor")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, sensor := range sensors {
		if sensor.Parent != "" && !strings.Contains(sensor.Name, "Temperature") {
			fmt.Printf(" Parent: %s, Sensor: %s, Temperature: %.1f°C\n", sensor.Parent, sensor.Name, sensor.Value)
		}
	}
}

/*
package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

func main() {
	// 执行nvidia-smi命令
	nvidiaSmi := exec.Command("nvidia-smi", "--query-gpu=index,name,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used", "--format=csv,noheader,nounits")
	nvidiaSmi.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// 获取命令输出
	var out bytes.Buffer
	nvidiaSmi.Stdout = &out
	err := nvidiaSmi.Run()
	if err != nil {
		fmt.Printf("Error running nvidia-smi: %v\n", err)
		return
	}

	// 输出结果
	fmt.Printf("GPU information:\n%s", out.String())
}

*/
