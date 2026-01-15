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
