Mouseion collects log data from applications

Why? Other logging applications are:
* Expensive. LogTail is one of the most affordable, and jumps to $30/month after their free tier.
* Overcomplicated. Small software businesses don't need complex log analysis. They need access to the logs, usually within a period of time.
* Difficult to get started. Self-hosting is a poorly-documented nightmare, and the load tools are even worse.

What Mouseion is:
* Simple. Easy to set up, easy to load with data, easy to use
* Affordable. Self-hostable, affordable hosted option
* Efficient. Takes steps to reduce the storage for customer data (and pass that on to the customer)

What Mouseion is not:
* Massively scaleable. I don't know how much Mouseion can handle right now, but it's probably not going to handle Facebook-level logging.
* Machine-learning AI Smart Big Data Analytics. It's simple. It's not going to tell you how many times the word "cat" has appeared in your logs
* Bloated. We're not going to implement all the features every business needs. Again, it's simple.

How it works:
1. A _source_ sends log data to the Mouseion _server_. The source sends the entry with a timestamp, the entry text, and a set of tags (if desired)
2. The server stores those log entries, deduplicating as necessary.
3. When something happens, you log in with the _client_ and look at your entries by time and optionally by tag.
4. You fix things and save the day.
