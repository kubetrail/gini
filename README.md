# gem
CLI for simple i/o with Google Gemini AI model

## disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.

> Don't trust, verify

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

```bash
go install github.com/kubetrail/gem@latest
```

Get an API key from [Google AI studio](https://makersuite.google.com/app/apikey)
and setup an env variable `API_KEY` for it.

## usage
```bash
gem qa
```

## example chat history

```bash
gem qa
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
