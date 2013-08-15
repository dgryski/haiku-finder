# haiku-finder

This program searches through stdin for sentences that match a 5-7-5 syllable count and prints them out, formatted as a haiku.

You will need the cmudict syllable dictionary.  Download it with:

    curl -o cmudict.0.7a 'http://sourceforge.net/p/cmusphinx/code/11879/tree/trunk/cmudict/cmudict.0.7a?format=raw'

Example:

    bash$ go run main.go </usr/share/common-licenses/GPL-3 |head -3
    When we speak of free
    software, we are referring
    to freedom, not price.

