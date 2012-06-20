# Twitter Topword Tracker

This tool follows a twitter stream and prints out the most frequently
encountered words over the given reporting frequency.

This is related to a concept we gave to a customer who wanted to solve
this problem a while back.  This is a significantly scaled down
version, but it can exercise a server pretty hard by doing the
following:

# What it Does

## Count Words

Every tweet that comes in is split into the words that make it up
and those words are counted within the tweet.  Each word is then
`incr`emented against the current time window.

## Maintain Top Lists

Also after each tweet, we maintain the top 1,000 words.  In order to
do this, we do a `cas` loop where we `get` the current list, do a
giant `multi-get` of the union of all of the words currently in the
list and all of the words we just found to find their current counts,
then update the list with this value.  The `cas` frequently fails, and
will likely fail more with more `-workers`.

## Reporting

Based on the `-interval` parameter, a report will be sent to stdout
showing the top of the top list.  That's currently just a simple `get`
and in-memory ops.

# Doing More

In general, twitter only provides a 1% stream of tweets.  Using the
`-multiplier` parameter, you can send every tweet through multiple
times such that `-multipler=100` should approximately equal the
traffic of all of twitter (though the words being processed would
vary, etc...).
