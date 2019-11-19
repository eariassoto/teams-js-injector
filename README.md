# Microsoft Team Javascript Injector

This program sends a javascript payload to be evaluated by Microsoft Team's Javascript VM. The program opens a Teams process with the Remote Debugging tool enabled. This debugger allow us to get the current pages/connections. The program then gets the address for the Chat service websocket and sends a message to it. The message request the "Runtime.evaluate" method that will make the application to execute the payload code.

Usage of ./teams-js-injector:
```
  -debug-port int
        Port number for Chromium remote debugging (default 9222)
  -payload-file string
        Javascript file to inject (default "payload.js")
  -teams-path string
        Location of Teams executable (default "C:\\Users\\%USERNAME%\\AppData\\Local\\Microsoft\\Teams\\current\\Teams.exe")
```

This program was inspired by [this article](https://medium.com/@dany74q/injecting-js-into-electron-apps-and-adding-rtl-support-for-microsoft-teams-d315dfb212a6)