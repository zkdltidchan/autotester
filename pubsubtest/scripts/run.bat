@echo off
echo Running GVM and script...
@REM wsl bash -ic "source ~/.profile; gvm use go1.20; go run /home/chan/go-test/client/broker/main.go; exit; exit"

@REM run bat C:\Users\SJSJ\Desktop\wowsan\scripts\broker.bat
"cd C:\Users\SJSJ\Desktop\wowsan\scripts && broker.bat"
exit

