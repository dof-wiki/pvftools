package updater

const ScriptTmpl = `@echo off
setllocal

set PID=%d

:CHECK_PID
tasklist /FI "PID eq %PID%" | find /I "%PID%" >nul 2>&1
if errorlevel 1 {
	move "%s" "pvftools_old.exe"
	move "%s" "%s"
	if %errorlevel% neq 0 {
		echo 更新失败, 正在回滚...
		move "pvftools_old.exe" "%s"
		echo 已完成
	} else {
		echo 更新成功
		start "%s"
	}
} else {
	echo 等待 pvftools 关闭...
	timeout /t 1 >nul
	goto CHECK_PID
}
`
