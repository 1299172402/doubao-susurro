Dim ws
Set ws = Wscript.CreateObject("Wscript.Shell")
ws.run "%s -silent", vbhide
Wscript.quit
