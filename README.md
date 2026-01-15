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

Demo output:

[pubsub]              | Executing: /google-cloud-sdk/platform/pubsub-emulator/bin/cloud-pubsub-emulator --host=0.0.0.0 --port=8085
[pubsub]              | [pubsub] This is the Google Pub/Sub fake.
[pubsub]              | [pubsub] Implementation may be incomplete or differ from the real system.
[pubsub]              | [pubsub] Jan 15, 2026 6:49:04 PM com.google.cloud.pubsub.testing.v1.Main main
[pubsub]              | [pubsub] INFO: IAM integration is disabled. IAM policy methods and ACL checks are not supported
[pubsub]              | [pubsub] Jan 15, 2026 6:49:05 PM com.google.cloud.pubsub.testing.v1.Main main
[pubsub]              | [pubsub] INFO: Server started, listening on 8085
[pubsub]              | [pubsub] Jan 15, 2026 6:49:06 PM io.gapi.emulators.netty.HttpVersionRoutingHandler channelRead
[pubsub]              | [pubsub] INFO: Detected HTTP/2 connection.
[pubsub]              | [pubsub] Jan 15, 2026 6:49:06 PM io.gapi.emulators.netty.HttpVersionRoutingHandler channelRead
[pubsub]              | [pubsub] INFO: Detected HTTP/2 connection.
[scanner]             | publishing {"ip":"1.1.1.118","port":29799,"service":"SSH","timestamp":1768502947,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogODA="}}
[processor]           | projects/test-project/subscriptions/scan-subscription
[scanner]             | publishing {"ip":"1.1.1.222","port":8283,"service":"HTTP","timestamp":1768502948,"data_version":2,"data":{"response_str":"service response: 89"}}
[processor]           | 2026/01/15 18:49:08 INFO processing !BADKEY="&{Ip:1.1.1.222 Port:8283 Service:HTTP Timestamp:1768502948 DataVersion:2 Data:map[response_str:service response: 89]}"
[scanner]             | publishing {"ip":"1.1.1.2","port":31682,"service":"HTTP","timestamp":1768502949,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMjE="}}
[processor]           | 2026/01/15 18:49:09 INFO processing !BADKEY="&{Ip:1.1.1.2 Port:31682 Service:HTTP Timestamp:1768502949 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMjE=]}"
[scanner]             | publishing {"ip":"1.1.1.222","port":64469,"service":"DNS","timestamp":1768502950,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogODk="}}
[processor]           | 2026/01/15 18:49:10 INFO processing !BADKEY="&{Ip:1.1.1.222 Port:64469 Service:DNS Timestamp:1768502950 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogODk=]}"
[scanner]             | publishing {"ip":"1.1.1.216","port":51545,"service":"DNS","timestamp":1768502951,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMzY="}}
[processor]           | 2026/01/15 18:49:11 INFO processing !BADKEY="&{Ip:1.1.1.216 Port:51545 Service:DNS Timestamp:1768502951 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMzY=]}"
[scanner]             | publishing {"ip":"1.1.1.129","port":5828,"service":"SSH","timestamp":1768502952,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMzE="}}
[processor]           | 2026/01/15 18:49:12 INFO processing !BADKEY="&{Ip:1.1.1.129 Port:5828 Service:SSH Timestamp:1768502952 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMzE=]}"
[scanner]             | publishing {"ip":"1.1.1.71","port":3230,"service":"DNS","timestamp":1768502953,"data_version":2,"data":{"response_str":"service response: 64"}}
[processor]           | 2026/01/15 18:49:13 INFO processing !BADKEY="&{Ip:1.1.1.71 Port:3230 Service:DNS Timestamp:1768502953 DataVersion:2 Data:map[response_str:service response: 64]}"
[scanner]             | publishing {"ip":"1.1.1.109","port":62717,"service":"SSH","timestamp":1768502954,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNg=="}}
[processor]           | 2026/01/15 18:49:14 INFO processing !BADKEY="&{Ip:1.1.1.109 Port:62717 Service:SSH Timestamp:1768502954 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNg==]}"
[scanner]             | publishing {"ip":"1.1.1.150","port":44795,"service":"SSH","timestamp":1768502955,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogODE="}}
[processor]           | 2026/01/15 18:49:15 INFO processing !BADKEY="&{Ip:1.1.1.150 Port:44795 Service:SSH Timestamp:1768502955 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogODE=]}"
[scanner]             | publishing {"ip":"1.1.1.67","port":55668,"service":"DNS","timestamp":1768502956,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNTU="}}
[processor]           | 2026/01/15 18:49:16 INFO processing !BADKEY="&{Ip:1.1.1.67 Port:55668 Service:DNS Timestamp:1768502956 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNTU=]}"
[scanner]             | publishing {"ip":"1.1.1.126","port":41912,"service":"SSH","timestamp":1768502957,"data_version":2,"data":{"response_str":"service response: 80"}}
[processor]           | 2026/01/15 18:49:17 INFO processing !BADKEY="&{Ip:1.1.1.126 Port:41912 Service:SSH Timestamp:1768502957 DataVersion:2 Data:map[response_str:service response: 80]}"
[scanner]             | publishing {"ip":"1.1.1.28","port":61262,"service":"HTTP","timestamp":1768502958,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNTI="}}
[processor]           | 2026/01/15 18:49:18 INFO processing !BADKEY="&{Ip:1.1.1.28 Port:61262 Service:HTTP Timestamp:1768502958 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNTI=]}"
[scanner]             | publishing {"ip":"1.1.1.197","port":22648,"service":"DNS","timestamp":1768502959,"data_version":2,"data":{"response_str":"service response: 96"}}
[processor]           | 2026/01/15 18:49:19 INFO processing !BADKEY="&{Ip:1.1.1.197 Port:22648 Service:DNS Timestamp:1768502959 DataVersion:2 Data:map[response_str:service response: 96]}"
[scanner]             | publishing {"ip":"1.1.1.1","port":10348,"service":"SSH","timestamp":1768502960,"data_version":2,"data":{"response_str":"service response: 38"}}
[processor]           | 2026/01/15 18:49:20 INFO processing !BADKEY="&{Ip:1.1.1.1 Port:10348 Service:SSH Timestamp:1768502960 DataVersion:2 Data:map[response_str:service response: 38]}"
[scanner]             | publishing {"ip":"1.1.1.129","port":63103,"service":"SSH","timestamp":1768502961,"data_version":2,"data":{"response_str":"service response: 4"}}
[processor]           | 2026/01/15 18:49:21 INFO processing !BADKEY="&{Ip:1.1.1.129 Port:63103 Service:SSH Timestamp:1768502961 DataVersion:2 Data:map[response_str:service response: 4]}"
[scanner]             | publishing {"ip":"1.1.1.6","port":32560,"service":"DNS","timestamp":1768502962,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogOA=="}}
[processor]           | 2026/01/15 18:49:22 INFO processing !BADKEY="&{Ip:1.1.1.6 Port:32560 Service:DNS Timestamp:1768502962 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogOA==]}"
[scanner]             | publishing {"ip":"1.1.1.99","port":56803,"service":"SSH","timestamp":1768502963,"data_version":2,"data":{"response_str":"service response: 0"}}
[processor]           | 2026/01/15 18:49:23 INFO processing !BADKEY="&{Ip:1.1.1.99 Port:56803 Service:SSH Timestamp:1768502963 DataVersion:2 Data:map[response_str:service response: 0]}"
[scanner]             | publishing {"ip":"1.1.1.192","port":36617,"service":"SSH","timestamp":1768502964,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMzg="}}
[processor]           | 2026/01/15 18:49:24 INFO processing !BADKEY="&{Ip:1.1.1.192 Port:36617 Service:SSH Timestamp:1768502964 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMzg=]}"
[scanner]             | publishing {"ip":"1.1.1.96","port":1111,"service":"DNS","timestamp":1768502965,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMTY="}}
[processor]           | 2026/01/15 18:49:25 INFO processing !BADKEY="&{Ip:1.1.1.96 Port:1111 Service:DNS Timestamp:1768502965 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMTY=]}"
[scanner]             | publishing {"ip":"1.1.1.226","port":59269,"service":"DNS","timestamp":1768502966,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMzY="}}
[processor]           | 2026/01/15 18:49:26 INFO processing !BADKEY="&{Ip:1.1.1.226 Port:59269 Service:DNS Timestamp:1768502966 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMzY=]}"
[scanner]             | publishing {"ip":"1.1.1.33","port":14297,"service":"DNS","timestamp":1768502967,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMQ=="}}
[processor]           | 2026/01/15 18:49:27 INFO processing !BADKEY="&{Ip:1.1.1.33 Port:14297 Service:DNS Timestamp:1768502967 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMQ==]}"
[scanner]             | publishing {"ip":"1.1.1.221","port":26075,"service":"DNS","timestamp":1768502968,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogOTg="}}
[processor]           | 2026/01/15 18:49:28 INFO processing !BADKEY="&{Ip:1.1.1.221 Port:26075 Service:DNS Timestamp:1768502968 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogOTg=]}"
[scanner]             | publishing {"ip":"1.1.1.187","port":3086,"service":"HTTP","timestamp":1768502969,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNjE="}}
[processor]           | 2026/01/15 18:49:29 INFO processing !BADKEY="&{Ip:1.1.1.187 Port:3086 Service:HTTP Timestamp:1768502969 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNjE=]}"
[scanner]             | publishing {"ip":"1.1.1.142","port":38047,"service":"SSH","timestamp":1768502970,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNTA="}}
[processor]           | 2026/01/15 18:49:30 INFO processing !BADKEY="&{Ip:1.1.1.142 Port:38047 Service:SSH Timestamp:1768502970 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNTA=]}"
[scanner]             | publishing {"ip":"1.1.1.196","port":36576,"service":"SSH","timestamp":1768502971,"data_version":2,"data":{"response_str":"service response: 69"}}
[processor]           | 2026/01/15 18:49:31 INFO processing !BADKEY="&{Ip:1.1.1.196 Port:36576 Service:SSH Timestamp:1768502971 DataVersion:2 Data:map[response_str:service response: 69]}"
[scanner]             | publishing {"ip":"1.1.1.136","port":48997,"service":"SSH","timestamp":1768502972,"data_version":2,"data":{"response_str":"service response: 53"}}
[processor]           | 2026/01/15 18:49:32 INFO processing !BADKEY="&{Ip:1.1.1.136 Port:48997 Service:SSH Timestamp:1768502972 DataVersion:2 Data:map[response_str:service response: 53]}"
[scanner]             | publishing {"ip":"1.1.1.79","port":61381,"service":"DNS","timestamp":1768502973,"data_version":2,"data":{"response_str":"service response: 52"}}
[processor]           | 2026/01/15 18:49:33 INFO processing !BADKEY="&{Ip:1.1.1.79 Port:61381 Service:DNS Timestamp:1768502973 DataVersion:2 Data:map[response_str:service response: 52]}"
[scanner]             | publishing {"ip":"1.1.1.239","port":22023,"service":"SSH","timestamp":1768502974,"data_version":2,"data":{"response_str":"service response: 50"}}
[processor]           | 2026/01/15 18:49:34 INFO processing !BADKEY="&{Ip:1.1.1.239 Port:22023 Service:SSH Timestamp:1768502974 DataVersion:2 Data:map[response_str:service response: 50]}"
[scanner]             | publishing {"ip":"1.1.1.91","port":16921,"service":"DNS","timestamp":1768502975,"data_version":2,"data":{"response_str":"service response: 91"}}
[processor]           | 2026/01/15 18:49:35 INFO processing !BADKEY="&{Ip:1.1.1.91 Port:16921 Service:DNS Timestamp:1768502975 DataVersion:2 Data:map[response_str:service response: 91]}"
[scanner]             | publishing {"ip":"1.1.1.209","port":22344,"service":"HTTP","timestamp":1768502976,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogMjU="}}
[processor]           | 2026/01/15 18:49:36 INFO processing !BADKEY="&{Ip:1.1.1.209 Port:22344 Service:HTTP Timestamp:1768502976 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogMjU=]}"
[scanner]             | publishing {"ip":"1.1.1.162","port":47473,"service":"HTTP","timestamp":1768502977,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogODE="}}
[processor]           | 2026/01/15 18:49:37 INFO processing !BADKEY="&{Ip:1.1.1.162 Port:47473 Service:HTTP Timestamp:1768502977 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogODE=]}"
[scanner]             | publishing {"ip":"1.1.1.181","port":59479,"service":"HTTP","timestamp":1768502978,"data_version":1,"data":{"response_bytes_utf8":"c2VydmljZSByZXNwb25zZTogNjE="}}
[processor]           | 2026/01/15 18:49:38 INFO processing !BADKEY="&{Ip:1.1.1.181 Port:59479 Service:HTTP Timestamp:1768502978 DataVersion:1 Data:map[response_bytes_utf8:c2VydmljZSByZXNwb25zZTogNjE=]}

Associated database available under `data/processor.db'.
