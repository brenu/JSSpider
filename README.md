# :spider_web:JSSpider

Tool that automatically navigates through websites in order to spot all of its .js files. Ideally created to be combined with other tools, like SecretFinder or LinkFinder.

## :gear: Installation

If you've got Go installed and configured, you can install JSSpider using the command below:

```console
foo@bar:~$ go get -u github.com/brenu/jsspider
```

## :arrow_forward: Usage

The usage of JSSpider consists of passing a list of domains by using the pipe, just like the examples below:

```console
foo@bar:~$ cat subdomains.txt | jsspider
foo@bar:~$ cat subdomains.txt | jsspider -o output.txt
foo@bar:~$ cat subdomains.txt | jsspider > output.txt
```

## :balance_scale: Disclaim :spider:

Use it with caution. You are responsible for your actions. I assume no liability and I'm not responsible for any misuse or damage.