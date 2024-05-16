@echo off
echo Running GVM and script...
@REM wsl bash -ic "source ~/.profile; gvm use go1.20; go run /home/chan/go-test/client/simulator/main.go; exit; exit"
"cd C:\Users\SJSJ\Desktop\wowsan && go .\cmd\simulator\main.go"
exit