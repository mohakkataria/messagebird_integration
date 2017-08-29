# Synopsis
MessageBird Golang SDK Integration.
Go Report Card : https://goreportcard.com/report/github.com/mohakkataria/messagebird_integration

# Working :

* An API that accepts SMS messages submitted via a POST request containing a JSON object as request body.
Example payloads
```
{"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}
{"recipient":"31612345678,31612345678","originator":"MessageBird","message":"This is a test message. 😀"}
```
   
* Validates the payload before queueing it to golang channel (being used as a FIFO queue). Other alternatives could have been using Redis as the con of this method is that it can only enable per process rate limiting.

* A go routine takes a message off the channel and proceeds onto sending the message via MessageBird SDK. We limit the number of requests made to the MessageBird per second to 1 using golang time ticker. Every second a message is processed for sending. Another approach to this would have been to save a last sent timestamp, and always compare before sending

* When an incoming message content/body is longer than 160 chars (70 in case of unicode), we split it into multiple parts (known as concatenated SMS) by setting the appropriate UDH (User Data Header). The rate limit specified above also applies to each segment of the SMS and only 1 segment is sent per second.

* When a plain or unicode message exceeds the maximum character limit, we truncate the message to 1377 and 603 characters for plain and unicode messages respectively.

* Unit tests have been added in the corresponding packages

* Using viper for configuration management of Access Key, Rate Limit factor etc.


# TODO :

* Currently, the implementation responds with a pending status to API, since we are queuing the message. We can maybe use a system to return the actual status while keeping the API call waiting for response, if required. Or respond with a transaction id(random) and then hit back the sender at a callback URL with response.
* Log sent messages somewhere to have some accountability.
* Log API errors in a better way.
* Increase Test Coverage.