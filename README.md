# cliconnect
a small cli tool to work with [kafka connect](https://docs.confluent.io/current/connect/intro.html).

## Installation

Linux
```
curl -Lo cliconnect https://github.com/wakeful/cliconnect/releases/download/0.2.0/cliconnect-linux-amd64 && chmod +x cliconnect && sudo mv cliconnect /usr/sbin/
```

src
```
go get -u github.com/wakeful/cliconnect
```

## Usage

```
$ cliconnect -h
Usage of cliconnect:
  -url string
        kafka connect url (default "http://127.0.0.1:28083")
  -version
        show version and exit
```

