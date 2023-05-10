# fyne-demo
https://www.udemy.com/course/building-gui-applications-with-fyne-and-go-golang
https://fyne.io/

I try put my own style in the implementation. If you get and compare the original course, you see that this code it's a little different, but in essence the same.

## Generate executables for your system (Windows, macOS, Linux)

1) Download the code
2) Go to the program folder of your interest to generate that executable
    1) ./hello-world: The most simple and first contact with Fyne
        * fyne widgets
    2) ./containers: A little more complex implementation of Hello World with Fyne
        * fyne containers + widgets
    3) ./markdown: The first complete functional program of a markdown editor with Fyne
        * fyne containers + widgets
        * fyne windows menus with actions
        * fyne unit tests
4) Run:
```sh
go install fyne.io/fyne/v2/cmd/fyne@latest
```
```sh
fyne package -appVersion 1.0.0 -name MarkDown -release -appID 1.0.0
```
