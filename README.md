# Complete programs 100% Golang backend code with Fyne V2 GUI
* https://www.udemy.com/course/building-gui-applications-with-fyne-and-go-golang
* https://fyne.io/

I tried to put my own "style" in the implementation. If you check out and compare the original course, you'll see that this code is a little different, but essentially the same.

The apps have a lot to improve. They don't have best coding practices, but I tried to do my best to improve the original course code a bit while still trying to finish the course on time.

## Generate executables for your system (Windows, macOS, Linux)

1) Download the code
2) Go to the program folder of your interest to generate that specific executable
    ### ./hello-world: The most simple and first contact with Fyne
    ![image](https://github.com/Agustincou/fyne-demo/assets/12106476/1adef7d5-adfa-4b35-a658-7f6f433ed63b)
    * fyne widgets
    
    ### ./containers: Slightly more complex implementation of above Hello World with Fyne
    ![image](https://github.com/Agustincou/fyne-demo/assets/12106476/bb28a4d0-d6db-4d8e-b57f-f893a1891b17)
    * fyne containers + widgets
    
    ### ./markdown: The first complete functional program of a markdown editor with Fyne
    ![image](https://github.com/Agustincou/fyne-demo/assets/12106476/35db1d87-5c00-4dea-87db-a7dfb066daca)
    * fyne containers + widgets
    * fyne windows menus with actions
    * fyne unit tests
    * fyne executable icon

    ### ./gold-watcher: The most complete program. Include and cover numerous features of Fyne V2
    ![image](https://github.com/Agustincou/fyne-demo/assets/12106476/2e0a0bdf-0cbc-4def-b855-513be9741d39)
    ![image](https://github.com/Agustincou/fyne-demo/assets/12106476/98925007-24d7-4fa0-886b-2e23e5c8d764)
    * fyne containers + widgets
    * fyne windows menus with actions
    * fyne unit tests
    * fyne preferences
    * api rest
    * local sql database

3) Run:
```sh
go install fyne.io/fyne/v2/cmd/fyne@latest
```
```sh
fyne package -appVersion 1.0.0 -name MarkDown -release -appID 1.0.0
```
(Only for Windows) The last project "Gold Watcher" has a makefile that you can use instead of the above command. Run:
```sh
make build
```
