# WasmStuff

This is the repository for WasmStuff.
A modification of relab/hotstuff that compiles to WebAssembly.

## Run it yourself

The easiest way of testing the system is found by downloading the Release Tag WasmStuff-1.0.
That release contains all the necessary files to run and test the system locally.

To demo the project on a browser perform the following steps:

1. Download the release named WasmStuff-1.0
2. Open a terminal and navigate to the folder 'HotstuffWASM/newNetwork'
3. Start the web server by inputting this command: 'websocket/websocket.exe'
4. Open 4 windows of http://127.0.0.1:8080/websocket/server.html
5. Choose Start Srv 1, 2, 3 and 4 in each window
6. For each window choose Benchmark or Chess(go to step 9.)
7. Input number of commands to benchmark (default 1000)
8. Click GO in all windows to start the benchmark - After completion execution times will be printed in the console
9. When the system is ready, buttons are enabled
10. Choose which other server to challenge to a game of chess

A commandline version is also included in this release.\\
Located in the folder 'HotstuffWASM/newNetwork/main'.\\
To demo the system from command line locate the 'main.exe' file.\\
The file take to command line arguments server ID and number of commands to run.\\
Example command: 'main.exe 1 500'\\
The commandline system can be executed separate or together with the browser setup.
