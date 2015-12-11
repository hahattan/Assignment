import urllib
import json

url = 'http://python-data.dr-chuck.net/comments_194442.json '
sum = 0
count = 0

uh = urllib.urlopen(url)
data = uh.read()	#type : string
print "Retrieved", len(data), "characters"
	
info = json.loads(data)

for iter in info['comments']:
	sum += iter['count']
	count += 1

print count
print sum