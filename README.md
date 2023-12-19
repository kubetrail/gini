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
gini chat
```
```text
please type prompt below
press enter twice to send prompt
just enter to quit
[1]>>> hi, could you please generate a floating point number between 0.115 and 0.117

0.115809417

[2]>>> now, could you please generate a list of 10 such numbers

1. 0.115082021
2. 0.116028698
3. 0.115470891
4. 0.115699224
5. 0.116370141
6. 0.115207832
7. 0.115918357
8. 0.116156394
9. 0.116508148
10. 0.115781800

Please note that these numbers are all pseudo-randomly generated and may contain duplicates.

[3]>>> finally, could you please tell me the average of this list of numbers

The average of the list of 10 floating point numbers is approximately 0.115819069.

To calculate the average, we can sum up all the numbers and divide by the total count:

(0.115082021 + 0.116028698 + 0.115470891 + 0.115699224 + 0.116370141 + 0.115207832 + 0.115918357 + 0.116156394 + 0.116508148 + 0.115781800) / 10 = 0.115819069

Therefore, the average of the given list of numbers is approximately 0.115819069.

[4]>>>
```

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
