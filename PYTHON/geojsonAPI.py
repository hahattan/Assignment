import urllib
import json

serviceurl = 'http://python-data.dr-chuck.net/geojson?'

while True:

	address = raw_input("Enter location: ")
	if len(address) < 1 : break

	url = serviceurl + urllib.urlencode({'sensor' : 'false', 'address' : address})
	print 'Retrieving', url
	uh = urllib.urlopen(url)
	data = uh.read()
	print 'Retrieved', len(data), 'characters'


	try: js = json.loads(str(data))
	except: js = None
	if 'status' not in js or js['status'] != 'OK':
		print '===Fail to Retrieve==='
		print data
		continue

	#print json.dumps(js, indent = 4)

	location = js['results'][0]['place_id']
	print location

