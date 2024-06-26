package updater

import (
	"bytes"
	"os"
	"pvftools/backend/common/log"
	"strconv"
	"text/template"
)

const ScriptTmpl = `@echo off
setllocal

set PID={{.Pid}}

:CHECK_PID
tasklist /FI "PID eq %PID%" | find /I "%PID%" >nul 2>&1
if errorlevel 1 {
	move "{{.OldExecutable}}" "pvftools_old.exe"
	move "{{.NewExecutable}}" "{{.OldExecutable}}"
	if %errorlevel% neq 0 {
		echo 更新失败, 正在回滚...
		move "pvftools_old.exe" "{{.OldExecutable}}"
		echo 已完成
	} else {
		echo 更新成功
		start "{{.OldExecutable}}"
	}
} else {
	echo 等待 pvftools 关闭...
	timeout /t 1 >nul
	goto CHECK_PID
}
`

func GenScript(oldExecutable, newExecutable string) {
	pid := os.Getpid()
	tmpl, err := template.New("script").Parse(ScriptTmpl)
	if err != nil {
		log.LogError("create tmpl err %v", err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{
		"Pid":           strconv.Itoa(pid),
		"OldExecutable": oldExecutable,
		"NewExecutable": newExecutable,
	})
	if err != nil {
		log.LogError("execute tmpl err %v", err)
		return
	}
	f, err := os.Create("pvftools_updater.bat")
	if err != nil {
		log.LogError("create update.bat err %v", err)
		return
	}
	defer f.Close()
	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.LogError("write update.bat err %v", err)
		return
	}
}
