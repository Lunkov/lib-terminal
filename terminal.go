package terminal

import (
  "os"
  "fmt"
	"bufio"
	"os/signal"
	"strings"
	"syscall"
)

// techEcho() - turns terminal echo on or off.
func termEcho(on bool) {
  // Common settings and variables for both stty calls.
  attrs := syscall.ProcAttr{
    Dir:   "",
    Env:   []string{},
    Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
    Sys:   nil}
  var ws syscall.WaitStatus
  cmd := "echo"
  if on == false {
    cmd = "-echo"
  }

  // Enable/disable echoing.
  pid, err := syscall.ForkExec(
    "/bin/stty",
    []string{"stty", cmd},
    &attrs)
  if err != nil {
    panic(err)
  }

  // Wait for the stty process to complete.
  _, err = syscall.Wait4(pid, &ws, 0, nil)
  if err != nil {
    panic(err)
  }
}

func GetPassword(prompt string) string {
  fmt.Print(prompt)

  // Catch a ^C interrupt.
  // Make sure that we reset term echo before exiting.
  signalChannel := make(chan os.Signal, 1)
  signal.Notify(signalChannel, os.Interrupt)
  go func() {
    for _ = range signalChannel {
      fmt.Println("\n^C interrupt.")
      termEcho(true)
      os.Exit(1)
    }
  }()

  // Echo is disabled, now grab the data.
  termEcho(false) // disable terminal echo
  reader := bufio.NewReader(os.Stdin)
  text, err := reader.ReadString('\n')
  termEcho(true) // always re-enable terminal echo
  fmt.Println("")
  if err != nil {
    // The terminal has been reset, go ahead and exit.
    fmt.Println("ERROR:", err.Error())
    os.Exit(1)
  }
  return strings.TrimSpace(text)
}

func GetText(prompt string) string {
  fmt.Print(prompt)
  reader := bufio.NewReader(os.Stdin)
  // ReadString will block until the delimiter is entered
  input, err := reader.ReadString('\n')
  if err != nil {
    fmt.Println("An error occured while reading input. Please try again", err)
    return ""
  }

  // remove the delimeter from the string
  return strings.TrimSuffix(input, "\n")
}

