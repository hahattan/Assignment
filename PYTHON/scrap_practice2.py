import urllib
import re
from BeautifulSoup import *

url = raw_input("Enter - ")

for i in range(7):
	print "Retrieving : ", url
	html = urllib.urlopen(url).read()
	soup = BeautifulSoup(html)

	tags = soup('a')
	iter = 1
	for tag in tags:
		# Look at the parts of a tag
	   #print 'TAG:',tag
	   #print 'URL:',tag.get('href', None)
	   #print 'Contents:',tag.contents[0]
	   #print 'Attrs:',tag.attrs

	   if iter == 18:
	   	url = tag.get('href', None)
	   	break
	   iter += 1

print "Last url : ", url
name = re.findall('by_([a-zA-Z]+)', url)
print name
print name[0]
#print type(name)
#print len(name)

