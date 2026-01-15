# Mini-Scan

[Commenting inline...]

> Hello!
>
> As you've heard by now, Censys scans the internet at an incredible scale. Processing the results necessitates scaling horizontally across thousands of machines. One key aspect of our architecture is the use of distributed queues to pass data between machines.
>
> ---
>
> The `docker-compose.yml` file sets up a toy example of a scanner. It spins up a Google Pub/Sub emulator, creates a topic and subscription, and publishes scan results to the topic. It can be run via `docker compose up`.
>
> Your job is to build the data processing side. It should:
>
> 1. Pull scan results from the subscription `scan-sub`.
Done.

> 2. Maintain an up-to-date record of each unique `(ip, port, service)`. This should contain when the service was last scanned and a string containing the service's response.
>
> > **_NOTE_**
> > The scanner can publish data in two formats, shown below. In both of the following examples, the service response should be stored as: `"hello world"`.
> >
> > ```javascript
> > {
> >   // ...
> >   "data_version": 1,
> >   "data": {
> >     "response_bytes_utf8": "aGVsbG8gd29ybGQ="
> >   }
> > }
> >
> > {
> >   // ...
> >   "data_version": 2,
> >   "data": {
> >     "response_str": "hello world"
> >   }
> > }
> > ```
> 
> Your processing application should be able to be scaled horizontally, but this isn't something you need to actually do. The processing application should use `at-least-once` semantics where ever applicable.
>
Done, though the horizontal scale is out-of-scope for this exercise alone.

> You may write this in any languages you choose, but Go would be preferred.
>
Golang it is.

> You may use any data store of your choosing, with `sqlite` being one example. Like our own code, we expect the code structure to make it easy to switch data stores.
>
Sqlite it is.

> Please note that Google Pub/Sub is best effort ordering and we want to keep the latest scan. While the example scanner does not publish scans at a rate where this would be an issue, we expect the application to be able to handle extreme out of orderness. Consider what would happen if the application received a scan that is 24 hours old.
>
Done.

> cmd/scanner/main.go should not be modified
>
I'll take this as soft "requirement". This is code. There is no limiting licensing restraining modifications, I'm under
NDA and the current architecture is rather... suboptimal. Especially, there is no reason to create the topic as a
separate step, increasing the dependency surface, making the whole system more fragile.

I corrected parts which I believed should be corrected, but merely left comment on other parts which were less critical,
while still being sub-optimal (if not totally wrong).

---

> Please upload the code to a publicly accessible GitHub, GitLab or other public code repository account. This README file should be updated, briefly documenting your solution. Like our own code, we expect testing instructions: whether it’s an automated test framework, or simple manual steps.
>
Please run: go test ./...

> To help set expectations, we believe you should aim to take no more than 4 hours on this task.
>

If you expect a run-off-the-mill basic solution certainly, however doesn't this defeat the whole point of this very
exercise ? After I was done with my implementation, I consulted various forks of this project, and rather expectedly, it
was all more or less the same underlying architecture, same underlying technical solutions implemented on various
backend.  Solutions only differed on the level of attention to details. The only open question was that I do not know
how many / if any of these people ended up receiving an offer.

Now, what is the point of this exercise if not to deliver a solution close to what would someone write on a day-to-day
basis ?

Initial thinking about this toy was about 4 hours to begin with, thinking about architecture, looking up SoTA, checking
libraries documentation, brainstorming with various LLM on possible technical implementation.  Also, yes, LLM have been
used, more to provide building blocks than code the whole project. Any Google search now return LLM results, they are
the new stackoverflow. Denying their use would have been hypocrite.

This game could certainly have been vibe-coded and end up to a working solution. What a vibe-coded solution does not
provide however is 1) personal knowledge of the architecture used, and 2) re-uses of various tools and layout I've built
over the years in order to built a maintainable solution as if the service was meant to be running in production for the
time being. Areas were left as mere comments for the future.

> We understand that you have other responsibilities, so if you think you’ll need more than 5 business days, just let us know when you expect to send a reply.
>
> Please don’t hesitate to ask any follow-up questions for clarification.

Runtime prep:
$ docker-compose build infra-builder-image
$ docker-compose build
$ docker-compose up

Demo output:

[processor]           | 2026/01/15 19:33:47 INFO Starting processor...
[processor]           | 2026/01/15 19:33:47 INFO Connecting to pub/sub...
[pubsub]              | Executing: /google-cloud-sdk/platform/pubsub-emulator/bin/cloud-pubsub-emulator --host=0.0.0.0 --port=8085
[pubsub]              | [pubsub] This is the Google Pub/Sub fake.
[pubsub]              | [pubsub] Implementation may be incomplete or differ from the real system.
[pubsub]              | [pubsub] Jan 15, 2026 7:33:49 PM com.google.cloud.pubsub.testing.v1.Main main
[pubsub]              | [pubsub] INFO: IAM integration is disabled. IAM policy methods and ACL checks are not supported
[pubsub]              | [pubsub] Jan 15, 2026 7:33:50 PM com.google.cloud.pubsub.testing.v1.Main main
[pubsub]              | [pubsub] INFO: Server started, listening on 8085
[pubsub]              | [pubsub] Jan 15, 2026 7:33:50 PM io.gapi.emulators.netty.HttpVersionRoutingHandler channelRead
[pubsub]              | [pubsub] INFO: Detected HTTP/2 connection.
[pubsub]              | [pubsub] Jan 15, 2026 7:33:52 PM io.gapi.emulators.netty.HttpVersionRoutingHandler channelRead
[pubsub]              | [pubsub] INFO: Detected HTTP/2 connection.
[scanner]             | publishing {"ip":"1.1.1.27","port":42595,"service":"DNS","timestamp":1768505634,"data_version":2,"data":{"response_str":"service response: 20"}}
[processor]           | 2026/01/15 19:33:54 INFO processing scan="&{Ip:1.1.1.27 Port:42595 Service:DNS Timestamp:1768505634 DataVersion:2 Data:map[response_str:service response: 20]}"
[scanner]             | publishing {"ip":"1.1.1.183","port":1480,"service":"SSH","timestamp":1768505635,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogOTk="}}
[processor]           | 2026/01/15 19:33:55 INFO processing scan="&{Ip:1.1.1.183 Port:1480 Service:SSH Timestamp:1768505635 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogOTk=]}"
[scanner]             | publishing {"ip":"1.1.1.202","port":34974,"service":"SSH","timestamp":1768505636,"data_version":2,"data":{"response_str":"service response: 90"}}
[processor]           | 2026/01/15 19:33:56 INFO processing scan="&{Ip:1.1.1.202 Port:34974 Service:SSH Timestamp:1768505636 DataVersion:2 Data:map[response_str:service response: 90]}"
[scanner]             | publishing {"ip":"1.1.1.241","port":11929,"service":"SSH","timestamp":1768505637,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNjE="}}
[processor]           | 2026/01/15 19:33:57 INFO processing scan="&{Ip:1.1.1.241 Port:11929 Service:SSH Timestamp:1768505637 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNjE=]}"
[scanner]             | publishing {"ip":"1.1.1.227","port":4002,"service":"SSH","timestamp":1768505638,"data_version":2,"data":{"response_str":"service response: 47"}}
[processor]           | 2026/01/15 19:33:58 INFO processing scan="&{Ip:1.1.1.227 Port:4002 Service:SSH Timestamp:1768505638 DataVersion:2 Data:map[response_str:service response: 47]}"
[scanner]             | publishing {"ip":"1.1.1.115","port":59507,"service":"HTTP","timestamp":1768505639,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNjg="}}
[processor]           | 2026/01/15 19:33:59 INFO processing scan="&{Ip:1.1.1.115 Port:59507 Service:HTTP Timestamp:1768505639 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNjg=]}"
[scanner]             | publishing {"ip":"1.1.1.94","port":27964,"service":"SSH","timestamp":1768505640,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNzY="}}
[processor]           | 2026/01/15 19:34:00 INFO processing scan="&{Ip:1.1.1.94 Port:27964 Service:SSH Timestamp:1768505640 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNzY=]}"
[scanner]             | publishing {"ip":"1.1.1.138","port":51092,"service":"HTTP","timestamp":1768505641,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNzk="}}
[processor]           | 2026/01/15 19:34:01 INFO processing scan="&{Ip:1.1.1.138 Port:51092 Service:HTTP Timestamp:1768505641 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNzk=]}"
[scanner]             | publishing {"ip":"1.1.1.245","port":25320,"service":"HTTP","timestamp":1768505642,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMjc="}}
[processor]           | 2026/01/15 19:34:02 INFO processing scan="&{Ip:1.1.1.245 Port:25320 Service:HTTP Timestamp:1768505642 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMjc=]}"
[scanner]             | publishing {"ip":"1.1.1.123","port":21802,"service":"DNS","timestamp":1768505643,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMTk="}}
[processor]           | 2026/01/15 19:34:03 INFO processing scan="&{Ip:1.1.1.123 Port:21802 Service:DNS Timestamp:1768505643 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMTk=]}"
[scanner]             | publishing {"ip":"1.1.1.229","port":43380,"service":"SSH","timestamp":1768505644,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNjU="}}
[processor]           | 2026/01/15 19:34:04 INFO processing scan="&{Ip:1.1.1.229 Port:43380 Service:SSH Timestamp:1768505644 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNjU=]}"
[scanner]             | publishing {"ip":"1.1.1.72","port":13110,"service":"HTTP","timestamp":1768505645,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNDU="}}
[processor]           | 2026/01/15 19:34:05 INFO processing scan="&{Ip:1.1.1.72 Port:13110 Service:HTTP Timestamp:1768505645 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNDU=]}"
[scanner]             | publishing {"ip":"1.1.1.170","port":13977,"service":"DNS","timestamp":1768505646,"data_version":2,"data":{"response_str":"service response: 92"}}
[processor]           | 2026/01/15 19:34:06 INFO processing scan="&{Ip:1.1.1.170 Port:13977 Service:DNS Timestamp:1768505646 DataVersion:2 Data:map[response_str:service response: 92]}"
[scanner]             | publishing {"ip":"1.1.1.187","port":28200,"service":"SSH","timestamp":1768505647,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMTQ="}}
[processor]           | 2026/01/15 19:34:07 INFO processing scan="&{Ip:1.1.1.187 Port:28200 Service:SSH Timestamp:1768505647 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMTQ=]}"
[scanner]             | publishing {"ip":"1.1.1.196","port":26474,"service":"HTTP","timestamp":1768505648,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogOTY="}}
[processor]           | 2026/01/15 19:34:08 INFO processing scan="&{Ip:1.1.1.196 Port:26474 Service:HTTP Timestamp:1768505648 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogOTY=]}"
[scanner]             | publishing {"ip":"1.1.1.162","port":8461,"service":"HTTP","timestamp":1768505649,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNTQ="}}
[processor]           | 2026/01/15 19:34:09 INFO processing scan="&{Ip:1.1.1.162 Port:8461 Service:HTTP Timestamp:1768505649 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNTQ=]}"
[scanner]             | publishing {"ip":"1.1.1.17","port":17695,"service":"SSH","timestamp":1768505650,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMg=="}}
[processor]           | 2026/01/15 19:34:10 INFO processing scan="&{Ip:1.1.1.17 Port:17695 Service:SSH Timestamp:1768505650 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMg==]}"
[scanner]             | publishing {"ip":"1.1.1.197","port":53840,"service":"HTTP","timestamp":1768505651,"data_version":2,"data":{"response_str":"service response: 58"}}
[processor]           | 2026/01/15 19:34:11 INFO processing scan="&{Ip:1.1.1.197 Port:53840 Service:HTTP Timestamp:1768505651 DataVersion:2 Data:map[response_str:service response: 58]}"
[scanner]             | publishing {"ip":"1.1.1.208","port":6356,"service":"HTTP","timestamp":1768505652,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMjQ="}}
[processor]           | 2026/01/15 19:34:12 INFO processing scan="&{Ip:1.1.1.208 Port:6356 Service:HTTP Timestamp:1768505652 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMjQ=]}"
[scanner]             | publishing {"ip":"1.1.1.243","port":20288,"service":"HTTP","timestamp":1768505653,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMzE="}}
[processor]           | 2026/01/15 19:34:13 INFO processing scan="&{Ip:1.1.1.243 Port:20288 Service:HTTP Timestamp:1768505653 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMzE=]}"
[scanner]             | publishing {"ip":"1.1.1.206","port":53480,"service":"SSH","timestamp":1768505654,"data_version":2,"data":{"response_str":"service response: 40"}}
[processor]           | 2026/01/15 19:34:14 INFO processing scan="&{Ip:1.1.1.206 Port:53480 Service:SSH Timestamp:1768505654 DataVersion:2 Data:map[response_str:service response: 40]}"
[scanner]             | publishing {"ip":"1.1.1.80","port":59266,"service":"DNS","timestamp":1768505655,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMjY="}}
[processor]           | 2026/01/15 19:34:15 INFO processing scan="&{Ip:1.1.1.80 Port:59266 Service:DNS Timestamp:1768505655 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMjY=]}"
[scanner]             | publishing {"ip":"1.1.1.178","port":11131,"service":"HTTP","timestamp":1768505656,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNTM="}}
[processor]           | 2026/01/15 19:34:16 INFO processing scan="&{Ip:1.1.1.178 Port:11131 Service:HTTP Timestamp:1768505656 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNTM=]}"
[scanner]             | publishing {"ip":"1.1.1.19","port":52745,"service":"SSH","timestamp":1768505657,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNDI="}}
[processor]           | 2026/01/15 19:34:17 INFO processing scan="&{Ip:1.1.1.19 Port:52745 Service:SSH Timestamp:1768505657 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNDI=]}"
[scanner]             | publishing {"ip":"1.1.1.19","port":29808,"service":"SSH","timestamp":1768505658,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNzg="}}
[processor]           | 2026/01/15 19:34:18 INFO processing scan="&{Ip:1.1.1.19 Port:29808 Service:SSH Timestamp:1768505658 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNzg=]}"
[scanner]             | publishing {"ip":"1.1.1.151","port":7710,"service":"DNS","timestamp":1768505659,"data_version":2,"data":{"response_str":"service response: 7"}}
[processor]           | 2026/01/15 19:34:19 INFO processing scan="&{Ip:1.1.1.151 Port:7710 Service:DNS Timestamp:1768505659 DataVersion:2 Data:map[response_str:service response: 7]}"
[scanner]             | publishing {"ip":"1.1.1.106","port":30682,"service":"HTTP","timestamp":1768505660,"data_version":2,"data":{"response_str":"service response: 84"}}
[processor]           | 2026/01/15 19:34:20 INFO processing scan="&{Ip:1.1.1.106 Port:30682 Service:HTTP Timestamp:1768505660 DataVersion:2 Data:map[response_str:service response: 84]}"
[scanner]             | publishing {"ip":"1.1.1.214","port":37807,"service":"DNS","timestamp":1768505661,"data_version":2,"data":{"response_str":"service response: 41"}}
[processor]           | 2026/01/15 19:34:21 INFO processing scan="&{Ip:1.1.1.214 Port:37807 Service:DNS Timestamp:1768505661 DataVersion:2 Data:map[response_str:service response: 41]}"
[scanner]             | publishing {"ip":"1.1.1.105","port":36359,"service":"SSH","timestamp":1768505662,"data_version":2,"data":{"response_str":"service response: 6"}}
[processor]           | 2026/01/15 19:34:22 INFO processing scan="&{Ip:1.1.1.105 Port:36359 Service:SSH Timestamp:1768505662 DataVersion:2 Data:map[response_str:service response: 6]}"
[scanner]             | publishing {"ip":"1.1.1.236","port":60301,"service":"SSH","timestamp":1768505663,"data_version":2,"data":{"response_str":"service response: 20"}}
[processor]           | 2026/01/15 19:34:23 INFO processing scan="&{Ip:1.1.1.236 Port:60301 Service:SSH Timestamp:1768505663 DataVersion:2 Data:map[response_str:service response: 20]}"
[scanner]             | publishing {"ip":"1.1.1.161","port":1961,"service":"DNS","timestamp":1768505664,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMA=="}}
[processor]           | 2026/01/15 19:34:24 INFO processing scan="&{Ip:1.1.1.161 Port:1961 Service:DNS Timestamp:1768505664 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMA==]}"
[scanner]             | publishing {"ip":"1.1.1.6","port":15777,"service":"DNS","timestamp":1768505665,"data_version":2,"data":{"response_str":"service response: 29"}}
[processor]           | 2026/01/15 19:34:25 INFO processing scan="&{Ip:1.1.1.6 Port:15777 Service:DNS Timestamp:1768505665 DataVersion:2 Data:map[response_str:service response: 29]}"
[scanner]             | publishing {"ip":"1.1.1.185","port":33171,"service":"HTTP","timestamp":1768505666,"data_version":2,"data":{"response_str":"service response: 1"}}
[processor]           | 2026/01/15 19:34:26 INFO processing scan="&{Ip:1.1.1.185 Port:33171 Service:HTTP Timestamp:1768505666 DataVersion:2 Data:map[response_str:service response: 1]}"
[scanner]             | publishing {"ip":"1.1.1.11","port":18535,"service":"DNS","timestamp":1768505667,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogOTI="}}
[processor]           | 2026/01/15 19:34:27 INFO processing scan="&{Ip:1.1.1.11 Port:18535 Service:DNS Timestamp:1768505667 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogOTI=]}"
[scanner]             | publishing {"ip":"1.1.1.159","port":3516,"service":"DNS","timestamp":1768505668,"data_version":2,"data":{"response_str":"service response: 29"}}
[processor]           | 2026/01/15 19:34:28 INFO processing scan="&{Ip:1.1.1.159 Port:3516 Service:DNS Timestamp:1768505668 DataVersion:2 Data:map[response_str:service response: 29]}"
[scanner]             | publishing {"ip":"1.1.1.131","port":33069,"service":"HTTP","timestamp":1768505669,"data_version":2,"data":{"response_str":"service response: 30"}}
[processor]           | 2026/01/15 19:34:29 INFO processing scan="&{Ip:1.1.1.131 Port:33069 Service:HTTP Timestamp:1768505669 DataVersion:2 Data:map[response_str:service response: 30]}"
[scanner]             | publishing {"ip":"1.1.1.77","port":52958,"service":"SSH","timestamp":1768505670,"data_version":2,"data":{"response_str":"service response: 26"}}
[processor]           | 2026/01/15 19:34:30 INFO processing scan="&{Ip:1.1.1.77 Port:52958 Service:SSH Timestamp:1768505670 DataVersion:2 Data:map[response_str:service response: 26]}"
[scanner]             | publishing {"ip":"1.1.1.89","port":3171,"service":"DNS","timestamp":1768505671,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogOTA="}}
[processor]           | 2026/01/15 19:34:31 INFO processing scan="&{Ip:1.1.1.89 Port:3171 Service:DNS Timestamp:1768505671 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogOTA=]}"
[scanner]             | publishing {"ip":"1.1.1.119","port":22366,"service":"DNS","timestamp":1768505672,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNDU="}}

Associated database available under `data/processor.db'.
