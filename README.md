# gini
CLI for simple i/o with Google Gemini AI model

## disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.

> Don't trust, verify

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

```bash
go install github.com/kubetrail/gini@latest
```

Get an API key from [Google AI studio](https://makersuite.google.com/app/apikey)
and setup an env variable `GOOGLE_API_KEY` for it.

Setup shell completion. See more info at:
```bash
gini completion -h
```

For instance, setup `bash` completion by adding following line to your `.bashrc`
```text
source <(gini completion bash)
```

## usage
```bash
gini chat [--auto-save]
```
`--auto-save` flag will save chat history to a randomly generated filename.

## example chat history

```bash
gini chat --auto-save
```
```text
please type prompt below
press enter twice to send prompt
just enter to quit
[1]>>> hi, could you please generate a floating point number between 0.115 and 0.117

      0.115997                     

[2]>>> could you please generate a list of 10 such numbers

      • 0.115023                   
      • 0.115193
      • 0.115316
      • 0.115408
      • 0.115628
      • 0.115746
      • 0.115859
      • 0.115997
      • 0.116031
      • 0.116173

      These are all floating point numbers between 0.115 and 0.117, generated
      randomly.

[3]>>> finally, could you please tell me the average of this list of numbers

      The average of the list of numbers
      • 0.115023
      • 0.115193
      • 0.115316
      • 0.115408
      • 0.115628
      • 0.115746
      • 0.115859
      • 0.115997
      • 0.116031
      • 0.116173

      is 0.115689.

      To find the average, we add up all the numbers and divide by the total
      number of numbers:

      (0.115023 + 0.115193 + 0.115316 + 0.115408 + 0.115628 + 0.115746 +
      0.115859 + 0.115997 + 0.116031 + 0.116173) / 10 = 1.15689 / 10 = 0.115689

      Therefore, the average of the list of numbers is 0.115689.

[4]>>> 
history saved to history-0d9d6887-ce12-4e89-824d-91b87b1a636f.txt
```

The prompt detects a blank line as termination, therefore, in order to send a prompt
that has blank lines in it, start the prompt with double curly braces `{{` and end
with `}}` as shown below.
```text
gini chat 
please type prompt below and press enter twice to send it
see more info in gini chat --help
just enter to quit
[1]>>> {{
hi, can you explain to me this code below

package main

func main() {
  fmt.Println("hello world")
}

}}
      The code below is a simple "Hello, world!" program written in Go. It
      prints the string "hello world" to the standard output.

      ┃ package main
      ┃ 
      ┃ import "fmt"
      ┃ 
      ┃ func main() {
      ┃     fmt.Println("hello world")
      ┃ }

      Here's a breakdown of the code:
      • package main: This line declares the package name for the program. The
        main package is the entry point for the program.
      • import "fmt": This line imports the fmt package, which provides
        functions for formatted I/O.
      • func main() { ... }: This is the main function of the program. It's the
        entry point for the program, and it's where the program execution
        begins.
      • fmt.Println("hello world"): This line prints the string "hello world" to
        the standard output (usually the console or terminal window where the
        program is running).

      When you run the program, you should see the output:

      ┃ hello world
```

## image analysis
Images can be analyzed using a combination of raw image data and associated text prompt.
Below is an example:
```bash
gini analyze image \
  --file seagull-on-a-rock.jpg \
  --format image/jpeg \
  could you please detect objects in the provided image
```
```text
 A seagull standing on a rock in front of a blurred background of green plants.
 The seagull is gray and white with a yellow beak. The rock is brown and the background is green.
```

> Please note that image analysis is conducted using `gemini-pro-vision` model by default.
> Furthermore, please make sure `--formats` match corresponding image format, or leave it
> blank if all images are jpeg images.

## list models
Following models can be selected when performing a task. Select model by via
`--model` flag using its name. For example `gini chat --model=models/gemini-pro-vision` etc.

```bash
gini list models
```
```yaml
name: models/gemini-pro
basemodeid: ""
version: "001"
displayname: Gemini Pro
description: The best model for scaling across a wide range of tasks
inputtokenlimit: 30720
outputtokenlimit: 2048
supportedgenerationmethods:
    - generateContent
    - countTokens
temperature: 0.9
topp: 1
topk: 1
```
```yaml
name: models/gemini-pro-vision
basemodeid: ""
version: "001"
displayname: Gemini Pro Vision
description: The best image understanding model to handle a broad range of applications
inputtokenlimit: 12288
outputtokenlimit: 4096
supportedgenerationmethods:
    - generateContent
    - countTokens
temperature: 0.4
topp: 1
topk: 32
```
```yaml
name: models/embedding-001
basemodeid: ""
version: "001"
displayname: Embedding 001
description: Obtain a distributed representation of a text.
inputtokenlimit: 2048
outputtokenlimit: 1
supportedgenerationmethods:
    - embedContent
    - countTextTokens
temperature: 0
topp: 0
topk: 0
```
## safety
`--allow-harm-probability` flag is set to `negligible` to prevent output from
displaying content that could be harmful. Change it at your own risk, for example,
```bash
gini chat --allow-harm-probability=medium --auto-save
```

## advanced config
Model config params such as `--top-p`, `--top-k`, `--temperature`, `--candiate-count` and 
`--max-output-tokens` can be supplied for fine tuning
